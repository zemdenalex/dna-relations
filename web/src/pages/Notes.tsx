import { useEffect, useState } from 'react'
import { api, Note } from '../api/client'

export default function Notes() {
  const [notes, setNotes] = useState<Note[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [editingNote, setEditingNote] = useState<Note | null>(null)
  const [formData, setFormData] = useState({ title: '', content: '', note_type: 'general' })
  const [searchQuery, setSearchQuery] = useState('')
  const [filterType, setFilterType] = useState<string>('')

  useEffect(() => {
    loadNotes()
  }, [filterType])

  const loadNotes = async () => {
    try {
      const data = await api.notes.list(filterType || undefined, 100)
      setNotes(data.notes)
    } catch (err) {
      console.error('Failed to load notes:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (editingNote) {
        await api.notes.update(editingNote.id, formData)
      } else {
        await api.notes.create({ ...formData, shared: true })
      }
      setShowForm(false)
      setEditingNote(null)
      setFormData({ title: '', content: '', note_type: 'general' })
      loadNotes()
    } catch (err) {
      console.error('Failed to save note:', err)
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

  const handleTogglePin = async (note: Note) => {
    try {
      await api.notes.togglePin(note.id)
      loadNotes()
    } catch (err) {
      console.error('Failed to toggle pin:', err)
    }
  }

  const startEdit = (note: Note) => {
    setEditingNote(note)
    setFormData({ title: note.title, content: note.content, note_type: note.note_type })
    setShowForm(true)
  }

  const filteredNotes = notes.filter(note => {
    if (!searchQuery) return true
    const query = searchQuery.toLowerCase()
    return note.title.toLowerCase().includes(query) || note.content.toLowerCase().includes(query)
  })

  const pinnedNotes = filteredNotes.filter(n => n.is_pinned)
  const regularNotes = filteredNotes.filter(n => !n.is_pinned)

  if (loading) {
    return <div className="text-gray-400">Loading...</div>
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-white">Notes</h1>
        <button
          onClick={() => { setShowForm(true); setEditingNote(null); setFormData({ title: '', content: '', note_type: 'general' }) }}
          className="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded"
        >
          New Note
        </button>
      </div>

      <div className="flex gap-4 mb-6">
        <input
          type="text"
          placeholder="Search notes..."
          value={searchQuery}
          onChange={e => setSearchQuery(e.target.value)}
          className="flex-1 bg-gray-800 text-white px-4 py-2 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500"
        />
        <select
          value={filterType}
          onChange={e => setFilterType(e.target.value)}
          className="bg-gray-800 text-white px-4 py-2 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500"
        >
          <option value="">All types</option>
          <option value="general">General</option>
          <option value="rule">Rules</option>
          <option value="idea">Ideas</option>
        </select>
      </div>

      {showForm && (
        <div className="bg-gray-800 rounded-lg p-6 mb-6">
          <h2 className="text-xl font-semibold text-white mb-4">
            {editingNote ? 'Edit Note' : 'New Note'}
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
              <label className="block text-gray-300 mb-2">Type</label>
              <select
                value={formData.note_type}
                onChange={e => setFormData({ ...formData, note_type: e.target.value })}
                className="w-full bg-gray-700 text-white px-4 py-2 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500"
              >
                <option value="general">General</option>
                <option value="rule">Rule</option>
                <option value="idea">Idea</option>
              </select>
            </div>
            <div>
              <label className="block text-gray-300 mb-2">Content</label>
              <textarea
                value={formData.content}
                onChange={e => setFormData({ ...formData, content: e.target.value })}
                className="w-full bg-gray-700 text-white px-4 py-2 rounded focus:outline-none focus:ring-2 focus:ring-indigo-500 h-40 font-mono"
              />
            </div>
            <div className="flex gap-2">
              <button type="submit" className="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded">
                {editingNote ? 'Update' : 'Create'}
              </button>
              <button
                type="button"
                onClick={() => { setShowForm(false); setEditingNote(null) }}
                className="bg-gray-600 hover:bg-gray-500 text-white px-4 py-2 rounded"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {filteredNotes.length === 0 ? (
        <div className="text-gray-500 text-center py-8">
          {searchQuery ? 'No notes found matching your search' : 'No notes yet. Create one!'}
        </div>
      ) : (
        <div className="space-y-6">
          {pinnedNotes.length > 0 && (
            <div>
              <h2 className="text-lg font-medium text-gray-400 mb-3">Pinned</h2>
              <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                {pinnedNotes.map(note => (
                  <NoteCard
                    key={note.id}
                    note={note}
                    onEdit={startEdit}
                    onDelete={handleDelete}
                    onTogglePin={handleTogglePin}
                  />
                ))}
              </div>
            </div>
          )}

          {regularNotes.length > 0 && (
            <div>
              {pinnedNotes.length > 0 && <h2 className="text-lg font-medium text-gray-400 mb-3">All Notes</h2>}
              <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                {regularNotes.map(note => (
                  <NoteCard
                    key={note.id}
                    note={note}
                    onEdit={startEdit}
                    onDelete={handleDelete}
                    onTogglePin={handleTogglePin}
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

function NoteCard({
  note,
  onEdit,
  onDelete,
  onTogglePin,
}: {
  note: Note
  onEdit: (note: Note) => void
  onDelete: (id: number) => void
  onTogglePin: (note: Note) => void
}) {
  const typeColors: Record<string, string> = {
    general: 'bg-gray-500/20 text-gray-400',
    rule: 'bg-red-500/20 text-red-400',
    idea: 'bg-yellow-500/20 text-yellow-400',
  }

  return (
    <div className="bg-gray-800 rounded-lg p-4 group">
      <div className="flex justify-between items-start mb-2">
        <div className="flex items-center gap-2">
          {note.is_pinned && <span>📌</span>}
          <h3 className="text-white font-medium">{note.title}</h3>
        </div>
        <span className={`text-xs px-2 py-1 rounded ${typeColors[note.note_type] || typeColors.general}`}>
          {note.note_type}
        </span>
      </div>
      {note.content && (
        <p className="text-gray-400 text-sm whitespace-pre-wrap line-clamp-4 mb-3">{note.content}</p>
      )}
      <div className="flex justify-between items-center">
        <p className="text-gray-500 text-xs">
          {new Date(note.created_at).toLocaleDateString()}
        </p>
        <div className="flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
          <button onClick={() => onTogglePin(note)} className="text-gray-400 hover:text-yellow-400 text-sm">
            {note.is_pinned ? 'Unpin' : 'Pin'}
          </button>
          <button onClick={() => onEdit(note)} className="text-indigo-400 hover:text-indigo-300 text-sm">
            Edit
          </button>
          <button onClick={() => onDelete(note.id)} className="text-red-400 hover:text-red-300 text-sm">
            Delete
          </button>
        </div>
      </div>
    </div>
  )
}