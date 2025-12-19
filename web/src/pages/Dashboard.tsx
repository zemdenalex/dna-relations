import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { api, Topic, CalendarEvent, Note } from '../api/client'

export default function Dashboard() {
  const [stats, setStats] = useState({ topicCount: 0, eventCount: 0, noteCount: 0 })
  const [pendingTopics, setPendingTopics] = useState<Topic[]>([])
  const [upcomingEvents, setUpcomingEvents] = useState<CalendarEvent[]>([])
  const [pinnedNotes, setPinnedNotes] = useState<Note[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    try {
      const now = new Date()
      const weekEnd = new Date(now)
      weekEnd.setDate(weekEnd.getDate() + 7)

      const [topicsRes, eventsRes, notesRes] = await Promise.all([
        api.topics.list('pending', 10),
        api.events.list(now.toISOString().split('T')[0], weekEnd.toISOString().split('T')[0]),
        api.notes.list(undefined, 50),
      ])

      const pinned = notesRes.notes.filter(n => n.is_pinned)

      setStats({
        topicCount: topicsRes.topics.length,
        eventCount: eventsRes.events.length,
        noteCount: notesRes.total,
      })

      setPendingTopics(topicsRes.topics.slice(0, 5))
      setUpcomingEvents(eventsRes.events.slice(0, 5))
      setPinnedNotes(pinned.slice(0, 5))
    } catch (err) {
      console.error('Failed to load dashboard data:', err)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div className="text-gray-400">Loading...</div>
  }

  return (
    <div>
      <h1 className="text-3xl font-bold text-white mb-2">Dashboard</h1>
      <p className="text-gray-400 mb-8">Welcome to DNA Relations</p>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
        <StatCard title="Topics" value={stats.topicCount} subtitle="pending" to="/topics" />
        <StatCard title="Events" value={stats.eventCount} subtitle="this week" to="/calendar" />
        <StatCard title="Notes" value={stats.noteCount} subtitle="total" to="/notes" />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Section title="Priority Topics" to="/topics">
          {pendingTopics.length === 0 ? (
            <p className="text-gray-500">No pending topics</p>
          ) : (
            <ul className="space-y-2">
              {pendingTopics.map(topic => (
                <li key={topic.id} className="text-gray-300">
                  <span className="text-indigo-400 mr-2">[P{topic.priority}]</span>
                  {topic.title}
                </li>
              ))}
            </ul>
          )}
        </Section>

        <Section title="Upcoming Events" to="/calendar">
          {upcomingEvents.length === 0 ? (
            <p className="text-gray-500">No upcoming events</p>
          ) : (
            <ul className="space-y-2">
              {upcomingEvents.map(event => {
                const date = new Date(event.start_at)
                return (
                  <li key={event.id} className="text-gray-300">
                    <span className="text-indigo-400 mr-2">
                      {date.toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric' })}
                    </span>
                    {event.title}
                  </li>
                )
              })}
            </ul>
          )}
        </Section>

        <Section title="Pinned Notes" to="/notes" className="md:col-span-2">
          {pinnedNotes.length === 0 ? (
            <p className="text-gray-500">No pinned notes</p>
          ) : (
            <ul className="space-y-2">
              {pinnedNotes.map(note => (
                <li key={note.id} className="text-gray-300">
                  <span className="mr-2">📌</span>
                  <span className="font-medium">{note.title}</span>
                  {note.content && (
                    <span className="text-gray-500 ml-2">
                      - {note.content.substring(0, 50)}{note.content.length > 50 ? '...' : ''}
                    </span>
                  )}
                </li>
              ))}
            </ul>
          )}
        </Section>
      </div>
    </div>
  )
}

function StatCard({ title, value, subtitle, to }: { title: string; value: number; subtitle: string; to: string }) {
  return (
    <Link to={to} className="bg-gray-800 rounded-lg p-6 hover:bg-gray-750 transition-colors">
      <div className="text-gray-400 text-sm">{title}</div>
      <div className="text-3xl font-bold text-white mt-1">{value}</div>
      <div className="text-gray-500 text-sm">{subtitle}</div>
    </Link>
  )
}

function Section({ title, to, children, className = '' }: { title: string; to: string; children: React.ReactNode; className?: string }) {
  return (
    <div className={`bg-gray-800 rounded-lg p-6 ${className}`}>
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold text-white">{title}</h2>
        <Link to={to} className="text-indigo-400 hover:text-indigo-300 text-sm">View all</Link>
      </div>
      {children}
    </div>
  )
}