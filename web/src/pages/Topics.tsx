import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { api, Topic, CreateTopicRequest } from '../api/client'

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

type TabType = 'pending' | 'discussed' | 'all'

export default function Topics() {
  const [topics, setTopics] = useState<Topic[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [activeTab, setActiveTab] = useState<TabType>('pending')
  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState<CreateTopicRequest>({
    title: '',
    description: '',
    priority: 3,
  })
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    loadTopics()
  }, [activeTab])

  const loadTopics = async () => {
    setLoading(true)
    try {
      const params: { status?: string; limit: number } = { limit: 50 }
      if (activeTab !== 'all') {
        params.status = activeTab
      }
      const resp = await api.topics.list(params)
      setTopics(resp.topics || [])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка загрузки')
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!formData.title.trim()) return

    setSaving(true)
    try {
      await api.topics.create(formData)
      setFormData({ title: '', description: '', priority: 3 })
      setShowForm(false)
      loadTopics()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка создания')
    } finally {
      setSaving(false)
    }
  }

  const markDiscussed = async (id: number) => {
    try {
      await api.topics.markDiscussed(id)
      loadTopics()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка')
    }
  }

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString('ru-RU', {
      day: 'numeric',
      month: 'short',
    })
  }

  return (
    <div className="topics-page">
      <header className="page-header">
        <Link to="/" className="back-btn">←</Link>
        <h1>Темы</h1>
        <button onClick={() => setShowForm(!showForm)} className="add-btn">
          {showForm ? '✕' : '+'}
        </button>
      </header>

      {error && <div className="error-banner">{error}</div>}

      {showForm && (
        <form onSubmit={handleSubmit} className="topic-form card">
          <div className="form-group">
            <label>Тема</label>
            <input
              type="text"
              value={formData.title}
              onChange={(e) => setFormData({ ...formData, title: e.target.value })}
              placeholder="О чём хотите поговорить?"
              required
            />
          </div>

          <div className="form-group">
            <label>Описание</label>
            <textarea
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              placeholder="Детали (необязательно)"
              rows={3}
            />
          </div>

          <div className="form-group">
            <label>Приоритет</label>
            <div className="priority-selector">
              {[0, 1, 2, 3, 4, 5].map((p) => (
                <button
                  key={p}
                  type="button"
                  className={`priority-btn ${formData.priority === p ? 'active' : ''}`}
                  style={{
                    backgroundColor: formData.priority === p ? priorityColors[p] : 'transparent',
                    borderColor: priorityColors[p],
                    color: formData.priority === p ? '#fff' : priorityColors[p],
                  }}
                  onClick={() => setFormData({ ...formData, priority: p })}
                >
                  {priorityLabels[p]}
                </button>
              ))}
            </div>
          </div>

          <button type="submit" className="btn-primary" disabled={saving}>
            {saving ? 'Сохранение...' : 'Добавить тему'}
          </button>
        </form>
      )}

      <div className="tabs">
        <button
          className={`tab ${activeTab === 'pending' ? 'active' : ''}`}
          onClick={() => setActiveTab('pending')}
        >
          Ожидают
        </button>
        <button
          className={`tab ${activeTab === 'discussed' ? 'active' : ''}`}
          onClick={() => setActiveTab('discussed')}
        >
          Обсуждённые
        </button>
        <button
          className={`tab ${activeTab === 'all' ? 'active' : ''}`}
          onClick={() => setActiveTab('all')}
        >
          Все
        </button>
      </div>

      {loading ? (
        <div className="loading">Загрузка...</div>
      ) : topics.length === 0 ? (
        <div className="empty-state">
          <p>Нет тем</p>
          <button onClick={() => setShowForm(true)} className="btn-primary">
            Добавить первую тему
          </button>
        </div>
      ) : (
        <ul className="topics-list-full">
          {topics.map((topic) => (
            <li key={topic.id} className="topic-card card">
              <div className="topic-header">
                <span
                  className="priority-badge"
                  style={{ backgroundColor: priorityColors[topic.priority] }}
                >
                  {priorityLabels[topic.priority]}
                </span>
                <span className="topic-date">{formatDate(topic.created_at)}</span>
              </div>
              <h3 className="topic-title">{topic.title}</h3>
              {topic.description && (
                <p className="topic-description">{topic.description}</p>
              )}
              {topic.status === 'pending' && (
                <button
                  onClick={() => markDiscussed(topic.id)}
                  className="btn-success"
                >
                  Обсудили ✓
                </button>
              )}
              {topic.status === 'discussed' && topic.discussed_at && (
                <p className="discussed-date">
                  Обсуждено: {formatDate(topic.discussed_at)}
                </p>
              )}
            </li>
          ))}
        </ul>
      )}

      <nav className="bottom-nav">
        <Link to="/" className="nav-item">
          <span className="nav-icon">🏠</span>
          <span>Главная</span>
        </Link>
        <Link to="/topics" className="nav-item active">
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
