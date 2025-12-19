import { useEffect, useState } from 'react'
import { api, Topic } from '../api/client'

const priorityLabels: Record<number, string> = {
  0: 'Buffer',
  1: 'Emergency',
  2: 'ASAP',
  3: 'Must discuss',
  4: 'Soon',
  5: 'When possible',
}

export default function Topics() {
  const [topics, setTopics] = useState<Topic[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [editingTopic, setEditingTopic] = useState<Topic | null>(null)
  const [formData, setFormData] = useState({ title: '', description: '', priority: 3 })
  const [filter, setFilter] = useState<'pending' | 'discussed' | 'all'>('pending')

  useEffect(() => {
    loadTopics()
  }, [filter])

  const loadTopics = async () => {
    try {
      const status = filter === 'all' ? undefined : filter
      const data = await api.topics.list(status, 50)
      setTopics(data.topics)
    } catch (err) {
      console.error('Failed to load topics:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (editingTopic) {
        await api.topics.update(editingTopic.id, formData)
      } else {
        await api.topics.create(formData)
      }
      setShowForm(false)
      setEditingTopic(null)
      setFormData({ title: '', description: '', priority: 3 })
      loadTopics()
    } catch (err) {
      console.error('Failed to save topic:', err)
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

  const handleMarkDiscussed = async (topic: Topic) => {
    try {
      await api.topics.markDiscussed(topic.id)
      loadTopics()
    } catch (err) {
      console.error('Failed to mark as discussed:', err)
    }
  }

  const startEdit = (topic: Topic) => {
    setEditingTopic(topic)
    setFormData({ title: topic.title, description: topic.description, priority: topic.priority })
    setShowForm(true)
  }

  if (loading) {
    return <div className="text-gray-400">Loading...</div>
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-white">Topics</h1>
        <button
          onClick={() => { setShowForm(true); setEditingTopic(null); setFormData({ title: '', description: '', priority: 3 }) }}
          className="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded"
        >
          New Topic
        </button>
      </div>

      <div className="flex gap-2 mb-6">
        {(['pending', 'discussed', 'all'] as const).map(f => (
          <button
            key={f}
            onClick={() => setFilter(f)}
            className={`px-4 py-2 rounded ${filter === f ? 'bg-indigo-600 text-white' : 'bg-gray-700 text-gray-300 hover:bg-gray-600'}`}
          >
            {f.charAt(0).toUpperCase() + f.slice(1)}
          </button>
        ))}
      </div>

      {showForm && (
        <div className="bg-gray-800 rounded-lg p-6 mb-6">
          <h2 className="text-xl font-semibold text-white mb-4">
            {editingTopic ? 'Edit Topic' : 'New Topic'}
          </h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-gray-300 mb-2">Title</label>
              <input
                type="text"
                value={formData.title}
                onChange={e => setFormData({ ...formData, title: e.target.value })}
                className="w-full bg-gray-700 text-white px-4 py-2 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500"
                required
              />
            </div>
            <div>
              <label className="block text-gray-300 mb-2">Description</label>
              <textarea
                value={formData.description}
                onChange={e => setFormData({ ...formData, description: e.target.value })}
                className="w-full bg-gray-700 text-white px-4 py-2 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500 h-24"
              />
            </div>
            <div>
              <label className="block text-gray-300 mb-2">Priority</label>
              <select
                value={formData.priority}
                onChange={e => setFormData({ ...formData, priority: parseInt(e.target.value) })}
                className="w-full bg-gray-700 text-white px-4 py-2 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500"
              >
                {Object.entries(priorityLabels).map(([value, label]) => (
                  <option key={value} value={value}>{value} - {label}</option>
                ))}
              </select>
            </div>
            <div className="flex gap-2">
              <button type="submit" className="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded">
                {editingTopic ? 'Update' : 'Create'}
              </button>
              <button
                type="button"
                onClick={() => { setShowForm(false); setEditingTopic(null) }}
                className="bg-gray-600 hover:bg-gray-500 text-white px-4 py-2 rounded"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {topics.length === 0 ? (
        <div className="text-gray-500 text-center py-8">No topics found</div>
      ) : (
        <div className="space-y-4">
          {topics.map(topic => (
            <div key={topic.id} className="bg-gray-800 rounded-lg p-4">
              <div className="flex justify-between items-start">
                <div className="flex-1">
                  <div className="flex items-center gap-2">
                    <h3 className={`text-lg font-medium ${topic.status === 'discussed' ? 'text-gray-500 line-through' : 'text-white'}`}>
                      {topic.title}
                    </h3>
                    <span className="bg-indigo-500/20 text-indigo-400 text-xs px-2 py-1 rounded">
                      P{topic.priority} - {priorityLabels[topic.priority]}
                    </span>
                    <span className={`text-xs px-2 py-1 rounded ${topic.status === 'pending' ? 'bg-yellow-500/20 text-yellow-400' : 'bg-green-500/20 text-green-400'}`}>
                      {topic.status}
                    </span>
                  </div>
                  {topic.description && (
                    <p className="text-gray-400 mt-2">{topic.description}</p>
                  )}
                  <p className="text-gray-500 text-sm mt-2">
                    Created {new Date(topic.created_at).toLocaleDateString()}
                  </p>
                </div>
                <div className="flex gap-2 ml-4">
                  {topic.status === 'pending' && (
                    <button
                      onClick={() => handleMarkDiscussed(topic)}
                      className="text-green-400 hover:text-green-300 text-sm"
                    >
                      Mark Discussed
                    </button>
                  )}
                  <button
                    onClick={() => startEdit(topic)}
                    className="text-indigo-400 hover:text-indigo-300 text-sm"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => handleDelete(topic.id)}
                    className="text-red-400 hover:text-red-300 text-sm"
                  >
                    Delete
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}