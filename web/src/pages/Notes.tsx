import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { api, Note, CreateNoteRequest } from '../api/client'

const noteTypeLabels: Record<string, string> = {
  general: 'Общее',
  rule: 'Правило',
  thought: 'Мысль',
  resource: 'Ресурс',
}

const noteTypeColors: Record<string, string> = {
  general: '#6366f1',
  rule: '#ec4899',
  thought: '#8b5cf6',
  resource: '#22c55e',
}

export default function Notes() {
  const [notes, setNotes] = useState<Note[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [filter, setFilter] = useState<string>('')
  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState<CreateNoteRequest>({
    title: '',
    content: '',
    note_type: 'general',
    shared: true,
  })
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    loadNotes()
  }, [filter])

  const loadNotes = async () => {
    setLoading(true)
    try {
      const params: { note_type?: string; limit: number } = { limit: 100 }
      if (filter) {
        params.note_type = filter
      }
      const resp = await api.notes.list(params)
      setNotes(resp.notes || [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка загрузки')
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!formData.content.trim()) return

    setSaving(true)
    try {
      await api.notes.create(formData)
      setFormData({ title: '', content: '', note_type: 'general', shared: true })
      setShowForm(false)
      loadNotes()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка создания')
    } finally {
      setSaving(false)
    }
  }

  const togglePin = async (id: number) => {
    try {
      await api.notes.togglePin(id)
      loadNotes()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка')
    }
  }

  const deleteNote = async (id: number) => {
    if (!confirm('Удалить заметку?')) return
    try {
      await api.notes.delete(id)
      loadNotes()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка удаления')
    }
  }

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString('ru-RU', {
      day: 'numeric',
      month: 'short',
    })
  }

  const pinnedNotes = notes.filter((n) => n.is_pinned)
  const regularNotes = notes.filter((n) => !n.is_pinned)

  return (
    <div className="notes-page">
      <header className="page-header">
        <Link to="/" className="back-btn">←</Link>
        <h1>Заметки</h1>
        <button onClick={() => setShowForm(!showForm)} className="add-btn">
          {showForm ? '✕' : '+'}
        </button>
      </header>

      {error && <div className="error-banner">{error}</div>}

      {showForm && (
        <form onSubmit={handleSubmit} className="note-form card">
          <div className="form-group">
            <label>Заголовок</label>
            <input
              type="text"
              value={formData.title}
              onChange={(e) => setFormData({ ...formData, title: e.target.value })}
              placeholder="Название (необязательно)"
            />
          </div>

          <div className="form-group">
            <label>Содержание</label>
            <textarea
              value={formData.content}
              onChange={(e) => setFormData({ ...formData, content: e.target.value })}
              placeholder="Текст заметки..."
              rows={5}
              required
            />
          </div>

          <div className="form-group">
            <label>Тип</label>
            <div className="type-selector">
              {Object.entries(noteTypeLabels).map(([type, label]) => (
                <button
                  key={type}
                  type="button"
                  className={`type-btn ${formData.note_type === type ? 'active' : ''}`}
                  style={{
                    backgroundColor: formData.note_type === type ? noteTypeColors[type] : 'transparent',
                    borderColor: noteTypeColors[type],
                    color: formData.note_type === type ? '#fff' : noteTypeColors[type],
                  }}
                  onClick={() => setFormData({ ...formData, note_type: type })}
                >
                  {label}
                </button>
              ))}
            </div>
          </div>

          <button type="submit" className="btn-primary" disabled={saving}>
            {saving ? 'Сохранение...' : 'Добавить заметку'}
          </button>
        </form>
      )}

      <div className="filter-tabs">
        <button
          className={`filter-tab ${filter === '' ? 'active' : ''}`}
          onClick={() => setFilter('')}
        >
          Все
        </button>
        {Object.entries(noteTypeLabels).map(([type, label]) => (
          <button
            key={type}
            className={`filter-tab ${filter === type ? 'active' : ''}`}
            onClick={() => setFilter(type)}
          >
            {label}
          </button>
        ))}
      </div>

      {loading ? (
        <div className="loading">Загрузка...</div>
      ) : notes.length === 0 ? (
        <div className="empty-state">
          <p>Нет заметок</p>
          <button onClick={() => setShowForm(true)} className="btn-primary">
            Создать первую заметку
          </button>
        </div>
      ) : (
        <div className="notes-container">
          {pinnedNotes.length > 0 && (
            <section className="pinned-section">
              <h2>Закреплённые</h2>
              <ul className="notes-list">
                {pinnedNotes.map((note) => (
                  <li key={note.id} className="note-card card pinned">
                    <div className="note-header">
                      <span
                        className="type-badge"
                        style={{ backgroundColor: noteTypeColors[note.note_type] }}
                      >
                        {noteTypeLabels[note.note_type]}
                      </span>
                      <span className="note-date">{formatDate(note.created_at)}</span>
                    </div>
                    {note.title && <h3 className="note-title">{note.title}</h3>}
                    <p className="note-content">{note.content}</p>
                    <div className="note-actions">
                      <button onClick={() => togglePin(note.id)} className="btn-icon">
                        📌
                      </button>
                      <button onClick={() => deleteNote(note.id)} className="btn-icon danger">
                        🗑
                      </button>
                    </div>
                  </li>
                ))}
              </ul>
            </section>
          )}

          {regularNotes.length > 0 && (
            <section className="regular-section">
              {pinnedNotes.length > 0 && <h2>Все заметки</h2>}
              <ul className="notes-list">
                {regularNotes.map((note) => (
                  <li key={note.id} className="note-card card">
                    <div className="note-header">
                      <span
                        className="type-badge"
                        style={{ backgroundColor: noteTypeColors[note.note_type] }}
                      >
                        {noteTypeLabels[note.note_type]}
                      </span>
                      <span className="note-date">{formatDate(note.created_at)}</span>
                    </div>
                    {note.title && <h3 className="note-title">{note.title}</h3>}
                    <p className="note-content">{note.content}</p>
                    <div className="note-actions">
                      <button onClick={() => togglePin(note.id)} className="btn-icon">
                        📌
                      </button>
                      <button onClick={() => deleteNote(note.id)} className="btn-icon danger">
                        🗑
                      </button>
                    </div>
                  </li>
                ))}
              </ul>
            </section>
          )}
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
        <Link to="/calendar" className="nav-item">
          <span className="nav-icon">📅</span>
          <span>Календарь</span>
        </Link>
        <Link to="/notes" className="nav-item active">
          <span className="nav-icon">📝</span>
          <span>Заметки</span>
        </Link>
      </nav>
    </div>
  )
}
