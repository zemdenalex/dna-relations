import { useEffect, useState } from 'react'
import { api, Ritual } from '../api/client'

export default function Rituals() {
  const [rituals, setRituals] = useState<Ritual[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    frequency: 'daily' as 'daily' | 'weekly' | 'monthly',
    time_of_day: '09:00',
  })

  useEffect(() => {
    loadRituals()
  }, [])

  const loadRituals = async () => {
    try {
      const data = await api.rituals.list()
      setRituals(data)
    } catch (err) {
      console.error('Failed to load rituals:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await api.rituals.create({ ...formData, is_active: true })
      setShowForm(false)
      setFormData({ title: '', description: '', frequency: 'daily', time_of_day: '09:00' })
      loadRituals()
    } catch (err) {
      console.error('Failed to create ritual:', err)
    }
  }

  const handleComplete = async (id: number) => {
    try {
      await api.rituals.complete(id)
      loadRituals()
    } catch (err) {
      console.error('Failed to complete ritual:', err)
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Delete this ritual?')) return
    try {
      await api.rituals.delete(id)
      loadRituals()
    } catch (err) {
      console.error('Failed to delete ritual:', err)
    }
  }

  const handleToggleActive = async (ritual: Ritual) => {
    try {
      await api.rituals.update(ritual.id, { is_active: !ritual.is_active })
      loadRituals()
    } catch (err) {
      console.error('Failed to toggle ritual:', err)
    }
  }

  const activeRituals = rituals.filter(r => r.is_active)
  const inactiveRituals = rituals.filter(r => !r.is_active)

  if (loading) {
    return <div className="text-gray-400">Loading...</div>
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-white">Rituals</h1>
        <button
          onClick={() => setShowForm(true)}
          className="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded"
        >
          New Ritual
        </button>
      </div>

      {showForm && (
        <div className="bg-gray-800 rounded-lg p-6 mb-6">
          <h2 className="text-xl font-semibold text-white mb-4">New Ritual</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-gray-300 mb-2">Title</label>
              <input
                type="text"
                value={formData.title}
                onChange={e => setFormData({ ...formData, title: e.target.value })}
                className="w-full bg-gray-700 text-white px-4 py-2 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500"
                placeholder="e.g., Morning check-in"
                required
              />
            </div>
            <div>
              <label className="block text-gray-300 mb-2">Description</label>
              <textarea
                value={formData.description}
                onChange={e => setFormData({ ...formData, description: e.target.value })}
                className="w-full bg-gray-700 text-white px-4 py-2 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500 h-20"
                placeholder="What does this ritual involve?"
              />
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-gray-300 mb-2">Frequency</label>
                <select
                  value={formData.frequency}
                  onChange={e => setFormData({ ...formData, frequency: e.target.value as 'daily' | 'weekly' | 'monthly' })}
                  className="w-full bg-gray-700 text-white px-4 py-2 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500"
                >
                  <option value="daily">Daily</option>
                  <option value="weekly">Weekly</option>
                  <option value="monthly">Monthly</option>
                </select>
              </div>
              <div>
                <label className="block text-gray-300 mb-2">Time</label>
                <input
                  type="time"
                  value={formData.time_of_day}
                  onChange={e => setFormData({ ...formData, time_of_day: e.target.value })}
                  className="w-full bg-gray-700 text-white px-4 py-2 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500"
                />
              </div>
            </div>
            <div className="flex gap-2">
              <button type="submit" className="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded">
                Create
              </button>
              <button
                type="button"
                onClick={() => setShowForm(false)}
                className="bg-gray-600 hover:bg-gray-500 text-white px-4 py-2 rounded"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {rituals.length === 0 ? (
        <div className="text-gray-500 text-center py-8">No rituals yet. Create one to build healthy habits together!</div>
      ) : (
        <div className="space-y-6">
          {activeRituals.length > 0 && (
            <div>
              <h2 className="text-lg font-medium text-gray-400 mb-3">Active Rituals</h2>
              <div className="space-y-3">
                {activeRituals.map(ritual => (
                  <RitualCard
                    key={ritual.id}
                    ritual={ritual}
                    onComplete={handleComplete}
                    onDelete={handleDelete}
                    onToggleActive={handleToggleActive}
                  />
                ))}
              </div>
            </div>
          )}

          {inactiveRituals.length > 0 && (
            <div>
              <h2 className="text-lg font-medium text-gray-400 mb-3">Inactive Rituals</h2>
              <div className="space-y-3">
                {inactiveRituals.map(ritual => (
                  <RitualCard
                    key={ritual.id}
                    ritual={ritual}
                    onComplete={handleComplete}
                    onDelete={handleDelete}
                    onToggleActive={handleToggleActive}
                  />
                ))}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  )
}

function RitualCard({
  ritual,
  onComplete,
  onDelete,
  onToggleActive,
}: {
  ritual: Ritual
  onComplete: (id: number) => void
  onDelete: (id: number) => void
  onToggleActive: (ritual: Ritual) => void
}) {
  const frequencyColors = {
    daily: 'bg-green-500/20 text-green-400',
    weekly: 'bg-blue-500/20 text-blue-400',
    monthly: 'bg-purple-500/20 text-purple-400',
  }

  return (
    <div className={`bg-gray-800 rounded-lg p-4 ${!ritual.is_active ? 'opacity-60' : ''}`}>
      <div className="flex justify-between items-start">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-1">
            <h3 className="text-white font-medium">{ritual.title}</h3>
            <span className={`text-xs px-2 py-1 rounded ${frequencyColors[ritual.frequency]}`}>
              {ritual.frequency}
            </span>
            <span className="text-gray-500 text-sm">{ritual.time_of_day}</span>
          </div>
          {ritual.description && <p className="text-gray-400 text-sm">{ritual.description}</p>}
        </div>
        <div className="flex gap-2 ml-4">
          {ritual.is_active && (
            <button
              onClick={() => onComplete(ritual.id)}
              className="bg-green-600 hover:bg-green-700 text-white px-3 py-1 rounded text-sm"
            >
              Complete
            </button>
          )}
          <button
            onClick={() => onToggleActive(ritual)}
            className="text-gray-400 hover:text-gray-300 text-sm"
          >
            {ritual.is_active ? 'Pause' : 'Resume'}
          </button>
          <button onClick={() => onDelete(ritual.id)} className="text-red-400 hover:text-red-300 text-sm">
            Delete
          </button>
        </div>
      </div>
    </div>
  )
}