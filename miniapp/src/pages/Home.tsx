import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { api, Topic, Event, Note } from '../lib/api'

function HomePage() {
  const [loading, setLoading] = useState(true)
  const [urgentTopics, setUrgentTopics] = useState<Topic[]>([])
  const [todayEvents, setTodayEvents] = useState<Event[]>([])
  const [pinnedNotes, setPinnedNotes] = useState<Note[]>([])
  const [stats, setStats] = useState({ topics: 0, events: 0, notes: 0 })

  useEffect(() => {
    loadDashboard()
  }, [])

  const loadDashboard = async () => {
    setLoading(true)
    try {
      const [topicsResp, eventsResp, notesResp] = await Promise.all([
        api.topics.list({ status: 'pending', limit: 5 }),
        api.events.list({
          from: new Date().toISOString().split('T')[0],
          to: new Date(Date.now() + 86400000).toISOString().split('T')[0],
        }),
        api.notes.list({ pinned_only: true, limit: 5 }),
      ])

      setUrgentTopics((topicsResp.topics || []).filter(t => t.priority <= 2))
      setTodayEvents(eventsResp.events || [])
      setPinnedNotes(notesResp.notes || [])
      setStats({
        topics: topicsResp.total,
        events: eventsResp.events?.length || 0,
        notes: notesResp.total,
      })
    } catch (err) {
      console.error('Failed to load dashboard:', err)
    } finally {
      setLoading(false)
    }
  }

  const user = window.Telegram?.WebApp?.initDataUnsafe?.user

  const handleLogout = () => {
    api.auth.logout()
    window.location.href = '/miniapp/login'
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-dna-primary">Loading...</div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <header className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-dna-primary">
            DNA Relations
          </h1>
          {user && (
            <p className="text-gray-400 text-sm mt-1">
              Hey {user.first_name}!
            </p>
          )}
        </div>
        <button
          onClick={handleLogout}
          className="text-sm text-gray-400 hover:text-white"
        >
          Logout
        </button>
      </header>

      <div className="grid grid-cols-3 gap-3">
        <StatCard label="Topics" value={stats.topics} to="/topics" />
        <StatCard label="Today" value={stats.events} to="/calendar" />
        <StatCard label="Notes" value={stats.notes} to="/notes" />
      </div>

      {urgentTopics.length > 0 && (
        <section className="bg-dna-surface rounded-lg p-4">
          <div className="flex justify-between items-center mb-3">
            <h2 className="font-semibold text-red-400">Urgent Topics</h2>
            <Link to="/topics" className="text-xs text-dna-primary">View all</Link>
          </div>
          <div className="space-y-2">
            {urgentTopics.map((topic) => (
              <div key={topic.id} className="flex items-center gap-2 text-sm">
                <span className={`w-2 h-2 rounded-full ${
                  topic.priority === 1 ? 'bg-red-500' : 'bg-orange-500'
                }`} />
                <span className="truncate">{topic.title}</span>
              </div>
            ))}
          </div>
        </section>
      )}

      {todayEvents.length > 0 && (
        <section className="bg-dna-surface rounded-lg p-4">
          <div className="flex justify-between items-center mb-3">
            <h2 className="font-semibold">Today</h2>
            <Link to="/calendar" className="text-xs text-dna-primary">Calendar</Link>
          </div>
          <div className="space-y-2">
            {todayEvents.map((event) => (
              <div key={event.id} className="flex items-center gap-2 text-sm">
                <span
                  className="w-2 h-2 rounded-full"
                  style={{ backgroundColor: event.color || '#6366f1' }}
                />
                <span className="text-gray-400 text-xs w-12">
                  {event.all_day
                    ? 'All day'
                    : new Date(event.start_at).toLocaleTimeString('en-US', {
                        hour: '2-digit',
                        minute: '2-digit',
                      })}
                </span>
                <span className="truncate">{event.title}</span>
              </div>
            ))}
          </div>
        </section>
      )}

      {pinnedNotes.length > 0 && (
        <section className="bg-dna-surface rounded-lg p-4">
          <div className="flex justify-between items-center mb-3">
            <h2 className="font-semibold text-yellow-400">Pinned</h2>
            <Link to="/notes" className="text-xs text-dna-primary">All notes</Link>
          </div>
          <div className="space-y-2">
            {pinnedNotes.map((note) => (
              <div key={note.id} className="text-sm">
                <span className="text-xs text-gray-500 mr-2">[{note.note_type}]</span>
                <span className="truncate">
                  {note.title || note.content.slice(0, 50)}
                </span>
              </div>
            ))}
          </div>
        </section>
      )}

      <section className="bg-dna-surface rounded-lg p-4">
        <h2 className="font-semibold mb-3">Quick Actions</h2>
        <div className="grid grid-cols-2 gap-2">
          <Link
            to="/topics"
            className="bg-dna-bg p-3 rounded-lg text-center text-sm hover:bg-dna-border"
          >
            Add Topic
          </Link>
          <Link
            to="/calendar"
            className="bg-dna-bg p-3 rounded-lg text-center text-sm hover:bg-dna-border"
          >
            Add Event
          </Link>
          <Link
            to="/notes"
            className="bg-dna-bg p-3 rounded-lg text-center text-sm hover:bg-dna-border"
          >
            Add Note
          </Link>
          <Link
            to="/notes?type=rule"
            className="bg-dna-bg p-3 rounded-lg text-center text-sm hover:bg-dna-border"
          >
            View Rules
          </Link>
        </div>
      </section>
    </div>
  )
}

function StatCard({ label, value, to }: { label: string; value: number; to: string }) {
  return (
    <Link
      to={to}
      className="bg-dna-surface rounded-lg p-3 text-center hover:bg-dna-border transition-colors"
    >
      <div className="text-2xl font-bold text-dna-primary">{value}</div>
      <div className="text-xs text-gray-400">{label}</div>
    </Link>
  )
}

export default HomePage
