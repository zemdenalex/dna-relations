import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { api, CalendarEvent, CreateEventRequest } from '../api/client'

const weekDays = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']
const monthNames = [
  'Январь', 'Февраль', 'Март', 'Апрель', 'Май', 'Июнь',
  'Июль', 'Август', 'Сентябрь', 'Октябрь', 'Ноябрь', 'Декабрь'
]

const eventColors = [
  '#6366f1', '#ec4899', '#22c55e', '#f97316', '#8b5cf6', '#ef4444'
]

export default function Calendar() {
  const [events, setEvents] = useState<CalendarEvent[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [currentDate, setCurrentDate] = useState(new Date())
  const [view, setView] = useState<'week' | 'month'>('week')
  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState<CreateEventRequest>({
    title: '',
    description: '',
    start_at: '',
    end_at: '',
    all_day: false,
    color: '#6366f1',
    shared: true,
  })
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    loadEvents()
  }, [currentDate, view])

  const loadEvents = async () => {
    setLoading(true)
    try {
      const { from, to } = getDateRange()
      const resp = await api.events.list({ from, to })
      setEvents(resp.events || [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка загрузки')
    } finally {
      setLoading(false)
    }
  }

  const getDateRange = () => {
    if (view === 'week') {
      const start = getWeekStart(currentDate)
      const end = new Date(start)
      end.setDate(end.getDate() + 6)
      return {
        from: start.toISOString().split('T')[0],
        to: end.toISOString().split('T')[0],
      }
    } else {
      const start = new Date(currentDate.getFullYear(), currentDate.getMonth(), 1)
      const end = new Date(currentDate.getFullYear(), currentDate.getMonth() + 1, 0)
      return {
        from: start.toISOString().split('T')[0],
        to: end.toISOString().split('T')[0],
      }
    }
  }

  const getWeekStart = (date: Date) => {
    const d = new Date(date)
    const day = d.getDay()
    const diff = d.getDate() - day + (day === 0 ? -6 : 1)
    return new Date(d.setDate(diff))
  }

  const getWeekDays = () => {
    const start = getWeekStart(currentDate)
    return Array.from({ length: 7 }, (_, i) => {
      const d = new Date(start)
      d.setDate(d.getDate() + i)
      return d
    })
  }

  const getMonthDays = () => {
    const year = currentDate.getFullYear()
    const month = currentDate.getMonth()
    const firstDay = new Date(year, month, 1)
    const lastDay = new Date(year, month + 1, 0)
    const days: (Date | null)[] = []

    let startPadding = firstDay.getDay() - 1
    if (startPadding < 0) startPadding = 6

    for (let i = 0; i < startPadding; i++) {
      days.push(null)
    }

    for (let i = 1; i <= lastDay.getDate(); i++) {
      days.push(new Date(year, month, i))
    }

    return days
  }

  const getEventsForDate = (date: Date) => {
    const dateStr = date.toISOString().split('T')[0]
    return events.filter((e) => e.start_at.split('T')[0] === dateStr)
  }

  const isToday = (date: Date) => {
    const today = new Date()
    return date.toDateString() === today.toDateString()
  }

  const navigatePrev = () => {
    const newDate = new Date(currentDate)
    if (view === 'week') {
      newDate.setDate(newDate.getDate() - 7)
    } else {
      newDate.setMonth(newDate.getMonth() - 1)
    }
    setCurrentDate(newDate)
  }

  const navigateNext = () => {
    const newDate = new Date(currentDate)
    if (view === 'week') {
      newDate.setDate(newDate.getDate() + 7)
    } else {
      newDate.setMonth(newDate.getMonth() + 1)
    }
    setCurrentDate(newDate)
  }

  const goToToday = () => {
    setCurrentDate(new Date())
  }

  const openAddForm = (date?: Date) => {
    const d = date || new Date()
    const dateStr = d.toISOString().split('T')[0]
    setFormData({
      ...formData,
      start_at: `${dateStr}T12:00`,
      end_at: `${dateStr}T13:00`,
    })
    setShowForm(true)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!formData.title.trim() || !formData.start_at) return

    setSaving(true)
    try {
      await api.events.create({
        ...formData,
        start_at: new Date(formData.start_at).toISOString(),
        end_at: formData.end_at ? new Date(formData.end_at).toISOString() : undefined,
      })
      setShowForm(false)
      setFormData({
        title: '',
        description: '',
        start_at: '',
        end_at: '',
        all_day: false,
        color: '#6366f1',
        shared: true,
      })
      loadEvents()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка создания')
    } finally {
      setSaving(false)
    }
  }

  const deleteEvent = async (id: number) => {
    if (!confirm('Удалить событие?')) return
    try {
      await api.events.delete(id)
      loadEvents()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка удаления')
    }
  }

  const formatTime = (dateStr: string) => {
    return new Date(dateStr).toLocaleTimeString('ru-RU', {
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  const getHeaderTitle = () => {
    if (view === 'week') {
      const days = getWeekDays()
      const start = days[0]
      const end = days[6]
      if (start.getMonth() === end.getMonth()) {
        return `${start.getDate()} - ${end.getDate()} ${monthNames[start.getMonth()]} ${start.getFullYear()}`
      }
      return `${start.getDate()} ${monthNames[start.getMonth()].slice(0, 3)} - ${end.getDate()} ${monthNames[end.getMonth()].slice(0, 3)} ${end.getFullYear()}`
    }
    return `${monthNames[currentDate.getMonth()]} ${currentDate.getFullYear()}`
  }

  return (
    <div className="calendar-page">
      <header className="page-header">
        <Link to="/" className="back-btn">←</Link>
        <h1>Календарь</h1>
        <button onClick={() => openAddForm()} className="add-btn">+</button>
      </header>

      {error && <div className="error-banner">{error}</div>}

      <div className="calendar-nav">
        <button onClick={navigatePrev} className="nav-btn">‹</button>
        <div className="calendar-title">
          <span>{getHeaderTitle()}</span>
          <button onClick={goToToday} className="today-btn">Сегодня</button>
        </div>
        <button onClick={navigateNext} className="nav-btn">›</button>
      </div>

      <div className="view-toggle">
        <button
          className={`toggle-btn ${view === 'week' ? 'active' : ''}`}
          onClick={() => setView('week')}
        >
          Неделя
        </button>
        <button
          className={`toggle-btn ${view === 'month' ? 'active' : ''}`}
          onClick={() => setView('month')}
        >
          Месяц
        </button>
      </div>

      {showForm && (
        <div className="modal-overlay" onClick={() => setShowForm(false)}>
          <form
            onSubmit={handleSubmit}
            className="event-form card"
            onClick={(e) => e.stopPropagation()}
          >
            <h2>Новое событие</h2>

            <div className="form-group">
              <label>Название</label>
              <input
                type="text"
                value={formData.title}
                onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                placeholder="Что запланировано?"
                required
              />
            </div>

            <div className="form-row">
              <div className="form-group">
                <label>Начало</label>
                <input
                  type="datetime-local"
                  value={formData.start_at}
                  onChange={(e) => setFormData({ ...formData, start_at: e.target.value })}
                  required
                />
              </div>

              <div className="form-group">
                <label>Конец</label>
                <input
                  type="datetime-local"
                  value={formData.end_at}
                  onChange={(e) => setFormData({ ...formData, end_at: e.target.value })}
                />
              </div>
            </div>

            <div className="form-group">
              <label>Описание</label>
              <textarea
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                placeholder="Детали (необязательно)"
                rows={2}
              />
            </div>

            <div className="form-group">
              <label>Цвет</label>
              <div className="color-picker">
                {eventColors.map((c) => (
                  <button
                    key={c}
                    type="button"
                    className={`color-btn ${formData.color === c ? 'active' : ''}`}
                    style={{ backgroundColor: c }}
                    onClick={() => setFormData({ ...formData, color: c })}
                  />
                ))}
              </div>
            </div>

            <div className="form-actions">
              <button type="button" onClick={() => setShowForm(false)} className="btn-secondary">
                Отмена
              </button>
              <button type="submit" className="btn-primary" disabled={saving}>
                {saving ? 'Сохранение...' : 'Создать'}
              </button>
            </div>
          </form>
        </div>
      )}

      {loading ? (
        <div className="loading">Загрузка...</div>
      ) : view === 'week' ? (
        <div className="week-view">
          <div className="week-header">
            {weekDays.map((day, i) => (
              <div key={i} className="week-day-header">{day}</div>
            ))}
          </div>
          <div className="week-grid">
            {getWeekDays().map((date, i) => {
              const dayEvents = getEventsForDate(date)
              return (
                <div
                  key={i}
                  className={`week-day ${isToday(date) ? 'today' : ''}`}
                  onClick={() => openAddForm(date)}
                >
                  <div className="day-number">{date.getDate()}</div>
                  <div className="day-events">
                    {dayEvents.map((event) => (
                      <div
                        key={event.id}
                        className="event-chip"
                        style={{ backgroundColor: event.color || '#6366f1' }}
                        onClick={(e) => {
                          e.stopPropagation()
                        }}
                      >
                        <span className="event-time">{formatTime(event.start_at)}</span>
                        <span className="event-title">{event.title}</span>
                        <button
                          className="delete-btn"
                          onClick={(e) => {
                            e.stopPropagation()
                            deleteEvent(event.id)
                          }}
                        >
                          ×
                        </button>
                      </div>
                    ))}
                  </div>
                </div>
              )
            })}
          </div>
        </div>
      ) : (
        <div className="month-view">
          <div className="month-header">
            {weekDays.map((day, i) => (
              <div key={i} className="month-day-header">{day}</div>
            ))}
          </div>
          <div className="month-grid">
            {getMonthDays().map((date, i) => {
              if (!date) {
                return <div key={i} className="month-day empty" />
              }
              const dayEvents = getEventsForDate(date)
              return (
                <div
                  key={i}
                  className={`month-day ${isToday(date) ? 'today' : ''}`}
                  onClick={() => openAddForm(date)}
                >
                  <div className="day-number">{date.getDate()}</div>
                  {dayEvents.length > 0 && (
                    <div className="event-dots">
                      {dayEvents.slice(0, 3).map((e) => (
                        <span
                          key={e.id}
                          className="event-dot"
                          style={{ backgroundColor: e.color || '#6366f1' }}
                        />
                      ))}
                    </div>
                  )}
                </div>
              )
            })}
          </div>
        </div>
      )}

      <nav className="bottom-nav">
        <Link to="/" className="nav-item">
          <span className="nav-icon">🏠</span>
          <span>Главная</span>
        </Link>
        <Link to="/topics" className="nav-item">
          <span className="nav-icon">💬</span>
          <span>Темы</span>
        </Link>
        <Link to="/calendar" className="nav-item active">
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
