import { useState, useEffect } from 'react'
import { api, Note } from '../lib/api'

type NoteType = 'general' | 'rule' | 'thought' | 'resource' | 'credential'

const NOTE_TYPE_LABELS: Record<NoteType, string> = {
  general: 'General',
  rule: 'Rules',
  thought: 'Thoughts',
  resource: 'Resources',
  credential: 'Credentials',
}

function NotesPage() {
  const [notes, setNotes] = useState<Note[]>([])
  const [total, setTotal] = useState(0)
  const [filter, setFilter] = useState<NoteType | 'all'>('all')
  const [showPinnedOnly, setShowPinnedOnly] = useState(false)
  const [loading, setLoading] = useState(true)
  const [showAdd, setShowAdd] = useState(false)
  const [newContent, setNewContent] = useState('')
  const [newType, setNewType] = useState<NoteType>('general')

  useEffect(() => {
    loadNotes()
  }, [filter, showPinnedOnly])

  const loadNotes = async () => {
    setLoading(true)
    try {
      const params: { note_type?: string; pinned_only?: boolean; limit?: number } = { limit: 50 }
      if (filter !== 'all') params.note_type = filter
      if (showPinnedOnly) params.pinned_only = true
      const resp = await api.notes.list(params)
      setNotes(resp.notes || [])
      setTotal(resp.total)
    } catch (err) {
      console.error('Failed to load notes:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!newContent.trim()) return

    try {
      await api.notes.create({ content: newContent, note_type: newType, shared: true })
      setNewContent('')
      setShowAdd(false)
      loadNotes()
    } catch (err) {
      console.error('Failed to create note:', err)
    }
  }

  const handleTogglePin = async (id: number) => {
    try {
      await api.notes.togglePin(id)
      loadNotes()
    } catch (err) {
      console.error('Failed to toggle pin:', err)
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Delete this note?')) return
    try {
      await api.notes.delete(id)
      loadNotes()
    } catch (err) {
      console.error('Failed to delete note:', err)
    }
  }

  return (
    <div className="space-y-4">
      <header className="flex justify-between items-center">
        <h1 className="text-xl font-bold">Notes ({total})</h1>
        <button
          onClick={() => setShowAdd(!showAdd)}
          className="bg-dna-primary px-4 py-2 rounded-lg text-sm"
        >
          {showAdd ? 'Cancel' : 'Add'}
        </button>
      </header>

      {showAdd && (
        <form onSubmit={handleCreate} className="bg-dna-surface rounded-lg p-4 space-y-3">
          <textarea
            value={newContent}
            onChange={(e) => setNewContent(e.target.value)}
            placeholder="Note content..."
            rows={3}
            className="w-full bg-dna-bg border border-dna-border rounded-lg px-3 py-2 text-sm resize-none"
            autoFocus
          />
          <div className="flex gap-2 items-center">
            <label className="text-sm text-gray-400">Type:</label>
            <select
              value={newType}
              onChange={(e) => setNewType(e.target.value as NoteType)}
              className="bg-dna-bg border border-dna-border rounded px-2 py-1 text-sm"
            >
              {(Object.keys(NOTE_TYPE_LABELS) as NoteType[]).map((t) => (
                <option key={t} value={t}>
                  {NOTE_TYPE_LABELS[t]}
                </option>
              ))}
            </select>
          </div>
          <button
            type="submit"
            className="w-full bg-dna-primary py-2 rounded-lg text-sm font-medium"
          >
            Create Note
          </button>
        </form>
      )}

      <div className="flex gap-2 flex-wrap">
        <FilterButton active={filter === 'all'} onClick={() => setFilter('all')}>
          All
        </FilterButton>
        {(Object.keys(NOTE_TYPE_LABELS) as NoteType[]).map((type) => (
          <FilterButton
            key={type}
            active={filter === type}
            onClick={() => setFilter(type)}
          >
            {NOTE_TYPE_LABELS[type]}
          </FilterButton>
        ))}
      </div>

      <div className="flex items-center gap-2">
        <input
          type="checkbox"
          id="pinnedOnly"
          checked={showPinnedOnly}
          onChange={(e) => setShowPinnedOnly(e.target.checked)}
          className="rounded"
        />
        <label htmlFor="pinnedOnly" className="text-sm text-gray-400">
          Pinned only
        </label>
      </div>

      {loading ? (
        <div className="text-center py-8 text-gray-400">Loading...</div>
      ) : notes.length === 0 ? (
        <div className="text-center py-8 text-gray-400">No notes found</div>
      ) : (
        <div className="space-y-2">
          {notes.map((note) => (
            <div key={note.id} className="bg-dna-surface rounded-lg p-3">
              <div className="flex items-start justify-between gap-2">
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 text-xs text-gray-400 mb-1">
                    <span className="bg-dna-border px-2 py-0.5 rounded">{note.note_type}</span>
                    {note.is_pinned && <span className="text-yellow-500">pinned</span>}
                  </div>
                  <div className="text-sm whitespace-pre-wrap">{note.content}</div>
                </div>
                <div className="flex gap-1 flex-shrink-0">
                  <button
                    onClick={() => handleTogglePin(note.id)}
                    className={`text-xs px-2 py-1 rounded ${
                      note.is_pinned ? 'bg-yellow-600' : 'bg-dna-border'
                    }`}
                  >
                    Pin
                  </button>
                  <button
                    onClick={() => handleDelete(note.id)}
                    className="text-xs bg-red-600/50 px-2 py-1 rounded"
                  >
                    Del
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      <section className="bg-dna-surface rounded-lg p-4">
        <h3 className="text-sm font-semibold mb-2 text-gray-400">Note Types</h3>
        <ul className="text-xs text-gray-400 space-y-1">
          <li><strong>Rules</strong> - Relationship agreements</li>
          <li><strong>Thoughts</strong> - Key realizations</li>
          <li><strong>Resources</strong> - Links, books, media</li>
          <li><strong>Credentials</strong> - Shared logins</li>
        </ul>
      </section>
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

export default NotesPage
