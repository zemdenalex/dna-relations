import { useEffect, useState } from 'react'
import { api, CalendarEvent } from '../api/client'

export default function Calendar() {
  const [events, setEvents] = useState<CalendarEvent[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [currentWeekStart, setCurrentWeekStart] = useState(() => {
    const now = new Date()
    const day = now.getDay()
    const diff = now.getDate() - day + (day === 0 ? -6 : 1)
    return new Date(now.setDate(diff))
  })
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    date: '',
    time: '12:00',
    all_day: false,
  })

  useEffect(() => {
    loadEvents()
  }, [currentWeekStart])

  const loadEvents = async () => {
    setLoading(true)
    try {
      const weekEnd = new Date(currentWeekStart)
      weekEnd.setDate(weekEnd.getDate() + 7)

      const data = await api.events.list(
        currentWeekStart.toISOString().split('T')[0],
        weekEnd.toISOString().split('T')[0]
      )
      setEvents(data.events)
    } catch (err) {
      console.error('Failed to load events:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      const startAt = formData.all_day
        ? `${formData.date}T00:00:00Z`
        : `${formData.date}T${formData.time}:00Z`

      const endAt = formData.all_day
        ? `${formData.date}T23:59:59Z`
        : new Date(new Date(startAt).getTime() + 60 * 60 * 1000).toISOString()

      await api.events.create({
        title: formData.title,
        description: formData.description,
        start_at: startAt,
        end_at: endAt,
        all_day: formData.all_day,
        shared: true,
      })

      setShowForm(false)
      setFormData({ title: '', description: '', date: '', time: '12:00', all_day: false })
      loadEvents()
    } catch (err) {
      console.error('Failed to create event:', err)
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Delete this event?')) return
    try {
      await api.events.delete(id)
      loadEvents()
    } catch (err) {
      console.error('Failed to delete event:', err)
    }
  }

  const previousWeek = () => {
    const newStart = new Date(currentWeekStart)
    newStart.setDate(newStart.getDate() - 7)
    setCurrentWeekStart(newStart)
  }

  const nextWeek = () => {
    const newStart = new Date(currentWeekStart)
    newStart.setDate(newStart.getDate() + 7)
    setCurrentWeekStart(newStart)
  }

  const getDaysOfWeek = () => {
    const days = []
    for (let i = 0; i < 7; i++) {
      const day = new Date(currentWeekStart)
      day.setDate(day.getDate() + i)
      days.push(day)
    }
    return days
  }

  const getEventsForDay = (day: Date) => {
    return events.filter(event => {
      const eventDate = new Date(event.start_at)
      return eventDate.toDateString() === day.toDateString()
    })
  }

  const formatWeekRange = () => {
    const weekEnd = new Date(currentWeekStart)
    weekEnd.setDate(weekEnd.getDate() + 6)
    return `${currentWeekStart.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })} - ${weekEnd.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}`
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-white">Calendar</h1>
        <button
          onClick={() => setShowForm(true)}
          className="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded"
        >
          New Event
        </button>
      </div>

      <div className="flex justify-between items-center mb-6">
        <button onClick={previousWeek} className="text-gray-400 hover:text-white px-4 py-2">
          Previous
        </button>
        <span className="text-white font-medium">{formatWeekRange()}</span>
        <button onClick={nextWeek} className="text-gray-400 hover:text-white px-4 py-2">
          Next
        </button>
      </div>

      {showForm && (
        <div className="bg-gray-800 rounded-lg p-6 mb-6">
          <h2 className="text-xl font-semibold text-white mb-4">New Event</h2>
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
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-gray-300 mb-2">Date</label>
                <input
                  type="date"
                  value={formData.date}
                  onChange={e => setFormData({ ...formData, date: e.target.value })}
                  className="w-full bg-gray-700 text-white px-4 py-2 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  required
                />
              </div>
              <div>
                <label className="block text-gray-300 mb-2">Time</label>
                <input
                  type="time"
                  value={formData.time}
                  onChange={e => setFormData({ ...formData, time: e.target.value })}
                  className="w-full bg-gray-700 text-white px-4 py-2 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  disabled={formData.all_day}
                />
              </div>
            </div>
            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                id="all_day"
                checked={formData.all_day}
                onChange={e => setFormData({ ...formData, all_day: e.target.checked })}
                className="rounded"
              />
              <label htmlFor="all_day" className="text-gray-300">All day event</label>
            </div>
            <div>
              <label className="block text-gray-300 mb-2">Description</label>
              <textarea
                value={formData.description}
                onChange={e => setFormData({ ...formData, description: e.target.value })}
                className="w-full bg-gray-700 text-white px-4 py-2 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500 h-20"
              />
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

      {loading ? (
        <div className="text-gray-400">Loading...</div>
      ) : (
        <div className="grid grid-cols-7 gap-2">
          {getDaysOfWeek().map(day => {
            const dayEvents = getEventsForDay(day)
            const isToday = day.toDateString() === new Date().toDateString()

            return (
              <div
                key={day.toISOString()}
                className={`bg-gray-800 rounded-lg p-3 min-h-32 ${isToday ? 'ring-2 ring-indigo-500' : ''}`}
              >
                <div className={`text-sm mb-2 ${isToday ? 'text-indigo-400 font-bold' : 'text-gray-400'}`}>
                  {day.toLocaleDateString('en-US', { weekday: 'short', day: 'numeric' })}
                </div>
                <div className="space-y-1">
                  {dayEvents.map(event => (
                    <div
                      key={event.id}
                      className="bg-indigo-600/30 text-indigo-200 text-xs p-2 rounded group relative"
                    >
                      <div className="font-medium">{event.title}</div>
                      {!event.all_day && (
                        <div className="text-indigo-300">
                          {new Date(event.start_at).toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' })}
                        </div>
                      )}
                      <button
                        onClick={() => handleDelete(event.id)}
                        className="absolute top-1 right-1 text-red-400 hover:text-red-300 opacity-0 group-hover:opacity-100"
                      >
                        x
                      </button>
                    </div>
                  ))}
                </div>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}