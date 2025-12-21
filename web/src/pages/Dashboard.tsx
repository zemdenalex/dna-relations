import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { api, Topic, CalendarEvent, Note } from '../api/client'
import { useAuth } from '../contexts/AuthContext'

const priorityLabels: Record<number, string> = {
  0: 'Буфер',
  1: 'Срочно',
  2: 'ASAP',
  3: 'Надо обсудить',
  4: 'Скоро',
  5: 'Когда-нибудь',
}

const priorityColors: Record<number, string> = {
  0: '#94a3b8',
  1: '#ef4444',
  2: '#f97316',
  3: '#eab308',
  4: '#22c55e',
  5: '#6366f1',
}

export default function Dashboard() {
  const { user, logout } = useAuth()
  const [topics, setTopics] = useState<Topic[]>([])
  const [events, setEvents] = useState<CalendarEvent[]>([])
  const [notes, setNotes] = useState<Note[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      const today = new Date()
      const nextWeek = new Date(today)
      nextWeek.setDate(nextWeek.getDate() + 7)

      const [topicsResp, eventsResp, notesResp] = await Promise.all([
        api.topics.list({ status: 'pending', limit: 5 }),
        api.events.list({
          from: today.toISOString().split('T')[0],
          to: nextWeek.toISOString().split('T')[0],
        }),
        api.notes.list({ limit: 5 }),
      ])

      setTopics(topicsResp.topics || [])
      setEvents(eventsResp.events || [])
      setNotes(notesResp.notes || [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка загрузки')
    } finally {
      setLoading(false)
    }
  }

  const formatDate = (dateStr: string) => {
    const date = new Date(dateStr)
    const today = new Date()
    const tomorrow = new Date(today)
    tomorrow.setDate(tomorrow.getDate() + 1)

    if (date.toDateString() === today.toDateString()) return 'Сегодня'
    if (date.toDateString() === tomorrow.toDateString()) return 'Завтра'

    return date.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' })
  }

  const formatTime = (dateStr: string) => {
    return new Date(dateStr).toLocaleTimeString('ru-RU', {
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  if (loading) {
    return <div className="loading">Загрузка...</div>
  }

  return (
    <div className="dashboard">
      <header className="dashboard-header">
        <div className="header-content">
          <h1>Привет, {user?.username}</h1>
          <button onClick={logout} className="btn-text">Выйти</button>
        </div>
      </header>

      {error && <div className="error-banner">{error}</div>}

      <div className="dashboard-grid">
        <section className="card topics-card">
          <div className="card-header">
            <h2>Темы для обсуждения</h2>
            <Link to="/topics" className="btn-link">Все</Link>
          </div>
          {topics.length === 0 ? (
            <p className="empty-state">Нет ожидающих тем</p>
          ) : (
            <ul className="topics-list">
              {topics.map((topic) => (
                <li key={topic.id} className="topic-item">
                  <span
                    className="priority-badge"
                    style={{ backgroundColor: priorityColors[topic.priority] }}
                  >
                    {priorityLabels[topic.priority]}
                  </span>
                  <span className="topic-title">{topic.title}</span>
                </li>
              ))}
            </ul>
          )}
          <Link to="/topics" className="btn-secondary">Добавить тему</Link>
        </section>

        <section className="card calendar-card">
          <div className="card-header">
            <h2>Ближайшие события</h2>
            <Link to="/calendar" className="btn-link">Календарь</Link>
          </div>
          {events.length === 0 ? (
            <p className="empty-state">Нет событий на неделю</p>
          ) : (
            <ul className="events-list">
              {events.slice(0, 5).map((event) => (
                <li key={event.id} className="event-item">
                  <div className="event-date">{formatDate(event.start_at)}</div>
                  <div className="event-info">
                    <span className="event-time">{formatTime(event.start_at)}</span>
                    <span className="event-title">{event.title}</span>
                  </div>
                </li>
              ))}
            </ul>
          )}
          <Link to="/calendar" className="btn-secondary">Открыть календарь</Link>
        </section>

        <section className="card notes-card">
          <div className="card-header">
            <h2>Заметки</h2>
            <Link to="/notes" className="btn-link">Все</Link>
          </div>
          {notes.length === 0 ? (
            <p className="empty-state">Нет заметок</p>
          ) : (
            <ul className="notes-list">
              {notes.map((note) => (
                <li key={note.id} className="note-item">
                  {note.is_pinned && <span className="pin-icon">📌</span>}
                  <span className="note-title">{note.title || note.content.slice(0, 50)}</span>
                </li>
              ))}
            </ul>
          )}
          <Link to="/notes" className="btn-secondary">Добавить заметку</Link>
        </section>
      </div>

      <nav className="bottom-nav">
        <Link to="/" className="nav-item active">
          <span className="nav-icon">🏠</span>
          <span>Главная</span>
        </Link>
        <Link to="/topics" className="nav-item">
          <span className="nav-icon">💬</span>
          <span>Темы</span>
        </Link>
        <Link to="/calendar" className="nav-item">
          <span className="nav-icon">📅</span>
          <span>Календарь</span>
        </Link>
        <Link to="/notes" className="nav-item">
          <span className="nav-icon">📝</span>
          <span>Заметки</span>
        </Link>
      </nav>
    </div>
  )
}
