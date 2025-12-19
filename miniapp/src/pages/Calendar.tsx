import { useState, useEffect } from 'react'
import { api, Event } from '../lib/api'

function CalendarPage() {
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(true)
  const [currentDate, setCurrentDate] = useState(new Date())
  const [view, setView] = useState<'week' | 'month'>('week')
  const [showAdd, setShowAdd] = useState(false)
  const [newTitle, setNewTitle] = useState('')
  const [newDate, setNewDate] = useState('')
  const [newTime, setNewTime] = useState('')
  const [newAllDay, setNewAllDay] = useState(false)

  useEffect(() => {
    loadEvents()
  }, [currentDate, view])

  const loadEvents = async () => {
    setLoading(true)
    try {
      const from = getStartDate()
      const to = getEndDate()
      const resp = await api.events.list({
        from: from.toISOString(),
        to: to.toISOString(),
      })
      setEvents(resp.events || [])
    } catch (err) {
      console.error('Failed to load events:', err)
    } finally {
      setLoading(false)
    }
  }

  const getStartDate = () => {
    const d = new Date(currentDate)
    if (view === 'week') {
      d.setDate(d.getDate() - d.getDay() + 1)
    } else {
      d.setDate(1)
    }
    d.setHours(0, 0, 0, 0)
    return d
  }

  const getEndDate = () => {
    const d = new Date(currentDate)
    if (view === 'week') {
      d.setDate(d.getDate() + (7 - d.getDay()))
    } else {
      d.setMonth(d.getMonth() + 1, 0)
    }
    d.setHours(23, 59, 59, 999)
    return d
  }

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!newTitle.trim() || !newDate) return

    try {
      const startAt = newAllDay || !newTime
        ? `${newDate}T00:00:00`
        : `${newDate}T${newTime}:00`

      await api.events.create({
        title: newTitle,
        start_at: startAt,
        all_day: newAllDay,
        shared: true,
      })
      setNewTitle('')
      setNewDate('')
      setNewTime('')
      setNewAllDay(false)
      setShowAdd(false)
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

  const navigatePrev = () => {
    const d = new Date(currentDate)
    if (view === 'week') {
      d.setDate(d.getDate() - 7)
    } else {
      d.setMonth(d.getMonth() - 1)
    }
    setCurrentDate(d)
  }

  const navigateNext = () => {
    const d = new Date(currentDate)
    if (view === 'week') {
      d.setDate(d.getDate() + 7)
    } else {
      d.setMonth(d.getMonth() + 1)
    }
    setCurrentDate(d)
  }

  const goToToday = () => {
    setCurrentDate(new Date())
  }

  const formatDateRange = () => {
    const start = getStartDate()
    const end = getEndDate()
    const opts: Intl.DateTimeFormatOptions = { month: 'short', day: 'numeric' }
    if (view === 'month') {
      return currentDate.toLocaleDateString('en-US', { month: 'long', year: 'numeric' })
    }
    return `${start.toLocaleDateString('en-US', opts)} - ${end.toLocaleDateString('en-US', opts)}`
  }

  const groupEventsByDate = () => {
    const grouped: Record<string, Event[]> = {}
    events.forEach((event) => {
      const date = new Date(event.start_at).toDateString()
      if (!grouped[date]) grouped[date] = []
      grouped[date].push(event)
    })
    return grouped
  }

  const grouped = groupEventsByDate()

  return (
    <div className="space-y-4">
      <header className="flex justify-between items-center">
        <h1 className="text-xl font-bold">Calendar</h1>
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
            placeholder="Event title..."
            className="w-full bg-dna-bg border border-dna-border rounded-lg px-3 py-2 text-sm"
            autoFocus
          />
          <div className="flex gap-2">
            <input
              type="date"
              value={newDate}
              onChange={(e) => setNewDate(e.target.value)}
              className="flex-1 bg-dna-bg border border-dna-border rounded-lg px-3 py-2 text-sm"
              required
            />
            {!newAllDay && (
              <input
                type="time"
                value={newTime}
                onChange={(e) => setNewTime(e.target.value)}
                className="bg-dna-bg border border-dna-border rounded-lg px-3 py-2 text-sm"
              />
            )}
          </div>
          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              id="allDay"
              checked={newAllDay}
              onChange={(e) => setNewAllDay(e.target.checked)}
              className="rounded"
            />
            <label htmlFor="allDay" className="text-sm text-gray-400">All day</label>
          </div>
          <button
            type="submit"
            className="w-full bg-dna-primary py-2 rounded-lg text-sm font-medium"
          >
            Create Event
          </button>
        </form>
      )}

      <div className="flex items-center justify-between">
        <div className="flex gap-2">
          <button
            onClick={() => setView('week')}
            className={`px-3 py-1 rounded-full text-sm ${
              view === 'week' ? 'bg-dna-primary' : 'bg-dna-surface text-gray-400'
            }`}
          >
            Week
          </button>
          <button
            onClick={() => setView('month')}
            className={`px-3 py-1 rounded-full text-sm ${
              view === 'month' ? 'bg-dna-primary' : 'bg-dna-surface text-gray-400'
            }`}
          >
            Month
          </button>
        </div>
        <button onClick={goToToday} className="text-sm text-dna-primary">
          Today
        </button>
      </div>

      <div className="flex items-center justify-between bg-dna-surface rounded-lg p-3">
        <button onClick={navigatePrev} className="p-2 hover:bg-dna-border rounded">&lt;</button>
        <span className="font-medium">{formatDateRange()}</span>
        <button onClick={navigateNext} className="p-2 hover:bg-dna-border rounded">&gt;</button>
      </div>

      {loading ? (
        <div className="text-center py-8 text-gray-400">Loading...</div>
      ) : events.length === 0 ? (
        <div className="text-center py-8 text-gray-400">No events in this period</div>
      ) : (
        <div className="space-y-4">
          {Object.entries(grouped)
            .sort(([a], [b]) => new Date(a).getTime() - new Date(b).getTime())
            .map(([date, dayEvents]) => (
              <div key={date}>
                <div className="text-sm font-semibold text-gray-400 mb-2">
                  {new Date(date).toLocaleDateString('en-US', {
                    weekday: 'short',
                    month: 'short',
                    day: 'numeric',
                  })}
                </div>
                <div className="space-y-2">
                  {dayEvents.map((event) => (
                    <div
                      key={event.id}
                      className="bg-dna-surface rounded-lg p-3 flex items-center gap-3"
                      style={{ borderLeft: `3px solid ${event.color || '#6366f1'}` }}
                    >
                      <div className="flex-1">
                        <div className="font-medium">{event.title}</div>
                        <div className="text-xs text-gray-400">
                          {event.all_day
                            ? 'All day'
                            : new Date(event.start_at).toLocaleTimeString('en-US', {
                                hour: '2-digit',
                                minute: '2-digit',
                              })}
                          {event.shared && ' • Shared'}
                        </div>
                      </div>
                      <button
                        onClick={() => handleDelete(event.id)}
                        className="text-xs bg-red-600/50 px-2 py-1 rounded"
                      >
                        Del
                      </button>
                    </div>
                  ))}
                </div>
              </div>
            ))}
        </div>
      )}
    </div>
  )
}

export default CalendarPage
