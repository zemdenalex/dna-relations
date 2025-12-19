const API_BASE = '/api/v1'

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
    credentials: 'include',
  })

  if (!res.ok) {
    const error = await res.text()
    throw new Error(error || res.statusText)
  }

  if (res.status === 204) {
    return {} as T
  }

  return res.json()
}

export const api = {
  auth: {
    login: (username: string, password: string) =>
      request<LoginResponse>('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ username, password }),
      }),
    logout: () =>
      request<void>('/auth/logout', { method: 'POST' }),
    me: () =>
      request<{ user: User }>('/auth/me'),
  },

  topics: {
    list: (status?: string, limit?: number) => {
      const params = new URLSearchParams()
      if (status) params.set('status', status)
      if (limit) params.set('limit', limit.toString())
      const query = params.toString()
      return request<TopicsResponse>(`/topics${query ? `?${query}` : ''}`)
    },
    get: (id: number) => request<{ topic: Topic }>(`/topics/${id}`),
    create: (data: CreateTopicRequest) =>
      request<{ topic: Topic }>('/topics', {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    update: (id: number, data: Partial<Topic>) =>
      request<{ topic: Topic }>(`/topics/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data),
      }),
    delete: (id: number) =>
      request<void>(`/topics/${id}`, { method: 'DELETE' }),
    markDiscussed: (id: number) =>
      request<{ topic: Topic }>(`/topics/${id}/discuss`, { method: 'POST' }),
  },

  events: {
    list: (from?: string, to?: string) => {
      const params = new URLSearchParams()
      if (from) params.set('from', from)
      if (to) params.set('to', to)
      const query = params.toString()
      return request<EventsResponse>(`/events${query ? `?${query}` : ''}`)
    },
    get: (id: number) => request<{ event: CalendarEvent }>(`/events/${id}`),
    create: (data: CreateEventRequest) =>
      request<{ event: CalendarEvent }>('/events', {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    update: (id: number, data: Partial<CalendarEvent>) =>
      request<{ event: CalendarEvent }>(`/events/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data),
      }),
    delete: (id: number) =>
      request<void>(`/events/${id}`, { method: 'DELETE' }),
  },

  notes: {
    list: (noteType?: string, limit?: number) => {
      const params = new URLSearchParams()
      if (noteType) params.set('note_type', noteType)
      if (limit) params.set('limit', limit.toString())
      const query = params.toString()
      return request<NotesResponse>(`/notes${query ? `?${query}` : ''}`)
    },
    get: (id: number) => request<{ note: Note }>(`/notes/${id}`),
    create: (data: CreateNoteRequest) =>
      request<{ note: Note }>('/notes', {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    update: (id: number, data: Partial<Note>) =>
      request<{ note: Note }>(`/notes/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data),
      }),
    delete: (id: number) =>
      request<void>(`/notes/${id}`, { method: 'DELETE' }),
    togglePin: (id: number) =>
      request<{ note: Note }>(`/notes/${id}/pin`, { method: 'POST' }),
  },
}

export interface User {
  id: number
  username: string
  display_name: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface Topic {
  id: number
  title: string
  description: string
  priority: number
  status: string
  created_by: number
  created_at: string
  discussed_at: string | null
}

export interface TopicsResponse {
  topics: Topic[]
  total: number
}

export interface CreateTopicRequest {
  title: string
  description?: string
  priority?: number
}

export interface CalendarEvent {
  id: number
  title: string
  description: string
  start_at: string
  end_at: string | null
  all_day: boolean
  color: string
  shared: boolean
  created_by: number
  created_at: string
}

export interface EventsResponse {
  events: CalendarEvent[]
}

export interface CreateEventRequest {
  title: string
  description?: string
  start_at: string
  end_at?: string
  all_day?: boolean
  color?: string
  shared?: boolean
}

export interface Note {
  id: number
  title: string
  content: string
  note_type: string
  is_pinned: boolean
  shared: boolean
  created_by: number
  created_at: string
  updated_at: string
}

export interface NotesResponse {
  notes: Note[]
  total: number
}

export interface CreateNoteRequest {
  title: string
  content?: string
  note_type?: string
  shared?: boolean
}