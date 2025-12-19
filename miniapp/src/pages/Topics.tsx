import { useState, useEffect } from 'react'
import { api, Topic } from '../lib/api'

const PRIORITY_LABELS: Record<number, string> = {
  0: 'Buffer',
  1: 'Emergency',
  2: 'ASAP',
  3: 'Must discuss',
  4: 'Soon',
  5: 'When possible',
}

const PRIORITY_COLORS: Record<number, string> = {
  0: 'bg-blue-500',
  1: 'bg-red-500',
  2: 'bg-orange-500',
  3: 'bg-yellow-500',
  4: 'bg-green-500',
  5: 'bg-gray-500',
}

function TopicsPage() {
  const [topics, setTopics] = useState<Topic[]>([])
  const [total, setTotal] = useState(0)
  const [filter, setFilter] = useState<'pending' | 'discussed' | 'all'>('pending')
  const [loading, setLoading] = useState(true)
  const [showAdd, setShowAdd] = useState(false)
  const [newTitle, setNewTitle] = useState('')
  const [newPriority, setNewPriority] = useState(3)

  useEffect(() => {
    loadTopics()
  }, [filter])

  const loadTopics = async () => {
    setLoading(true)
    try {
      const status = filter === 'all' ? '' : filter
      const resp = await api.topics.list({ status, limit: 50 })
      setTopics(resp.topics || [])
      setTotal(resp.total)
    } catch (err) {
      console.error('Failed to load topics:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!newTitle.trim()) return

    try {
      await api.topics.create({ title: newTitle, priority: newPriority })
      setNewTitle('')
      setNewPriority(3)
      setShowAdd(false)
      loadTopics()
    } catch (err) {
      console.error('Failed to create topic:', err)
    }
  }

  const handleMarkDiscussed = async (id: number) => {
    try {
      await api.topics.markDiscussed(id)
      loadTopics()
    } catch (err) {
      console.error('Failed to mark topic:', err)
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Delete this topic?')) return
    try {
      await api.topics.delete(id)
      loadTopics()
    } catch (err) {
      console.error('Failed to delete topic:', err)
    }
  }

  return (
    <div className="space-y-4">
      <header className="flex justify-between items-center">
        <h1 className="text-xl font-bold">Topics ({total})</h1>
        <button
          onClick={() => setShowAdd(!showAdd)}
          className="bg-dna-primary px-4 py-2 rounded-lg text-sm"
        >
          {showAdd ? 'Cancel' : 'Add'}
        </button>
      </header>

      {showAdd && (
        <form onSubmit={handleCreate} className="bg-dna-surface rounded-lg p-4 space-y-3">
          <input
            type="text"
            value={newTitle}
            onChange={(e) => setNewTitle(e.target.value)}
            placeholder="Topic title..."
            className="w-full bg-dna-bg border border-dna-border rounded-lg px-3 py-2 text-sm"
            autoFocus
          />
          <div className="flex gap-2 items-center">
            <label className="text-sm text-gray-400">Priority:</label>
            <select
              value={newPriority}
              onChange={(e) => setNewPriority(Number(e.target.value))}
              className="bg-dna-bg border border-dna-border rounded px-2 py-1 text-sm"
            >
              {[0, 1, 2, 3, 4, 5].map((p) => (
                <option key={p} value={p}>
                  {p} - {PRIORITY_LABELS[p]}
                </option>
              ))}
            </select>
          </div>
          <button
            type="submit"
            className="w-full bg-dna-primary py-2 rounded-lg text-sm font-medium"
          >
            Create Topic
          </button>
        </form>
      )}

      <div className="flex gap-2">
        <FilterButton active={filter === 'pending'} onClick={() => setFilter('pending')}>
          Pending
        </FilterButton>
        <FilterButton active={filter === 'discussed'} onClick={() => setFilter('discussed')}>
          Discussed
        </FilterButton>
        <FilterButton active={filter === 'all'} onClick={() => setFilter('all')}>
          All
        </FilterButton>
      </div>

      {loading ? (
        <div className="text-center py-8 text-gray-400">Loading...</div>
      ) : topics.length === 0 ? (
        <div className="text-center py-8 text-gray-400">No topics found</div>
      ) : (
        <div className="space-y-2">
          {topics.map((topic) => (
            <div
              key={topic.id}
              className="bg-dna-surface rounded-lg p-3 flex items-start gap-3"
            >
              <span className={`w-2 h-2 rounded-full mt-2 ${PRIORITY_COLORS[topic.priority]}`} />
              <div className="flex-1 min-w-0">
                <div className="font-medium">{topic.title}</div>
                <div className="text-xs text-gray-400 mt-1">
                  P{topic.priority} - {PRIORITY_LABELS[topic.priority]}
                </div>
              </div>
              <div className="flex gap-1">
                {topic.status === 'pending' && (
                  <button
                    onClick={() => handleMarkDiscussed(topic.id)}
                    className="text-xs bg-green-600 px-2 py-1 rounded"
                  >
                    Done
                  </button>
                )}
                <button
                  onClick={() => handleDelete(topic.id)}
                  className="text-xs bg-red-600/50 px-2 py-1 rounded"
                >
                  Del
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      <div className="bg-dna-surface rounded-lg p-4">
        <h3 className="text-sm font-semibold mb-2 text-gray-400">Priority Legend</h3>
        <div className="grid grid-cols-2 gap-2 text-xs">
          {[0, 1, 2, 3, 4, 5].map((p) => (
            <div key={p} className="flex items-center gap-2">
              <span className={`w-3 h-3 rounded-full ${PRIORITY_COLORS[p]}`} />
              <span>{PRIORITY_LABELS[p]}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

function FilterButton({
  active,
  onClick,
  children,
}: {
  active: boolean
  onClick: () => void
  children: React.ReactNode
}) {
  return (
    <button
      onClick={onClick}
      className={`px-3 py-1 rounded-full text-sm ${
        active ? 'bg-dna-primary text-white' : 'bg-dna-surface text-gray-400'
      }`}
    >
      {children}
    </button>
  )
}

export default TopicsPage
