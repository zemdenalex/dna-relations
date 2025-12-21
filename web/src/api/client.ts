const API_BASE = '/api/v1'

let authToken: string | null = localStorage.getItem('dna_token')

export function setToken(token: string) {
  authToken = token
  localStorage.setItem('dna_token', token)
}

export function getToken(): string | null {
  return authToken
}

export function clearToken() {
  authToken = null
  localStorage.removeItem('dna_token')
}

export function isAuthenticated(): boolean {
  return !!authToken
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }

  if (authToken) {
    headers['Authorization'] = `Bearer ${authToken}`
  }

  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers: {
      ...headers,
      ...options?.headers,
    },
  })

  if (res.status === 401) {
    clearToken()
    window.location.href = '/login'
    throw new Error('Unauthorized')
  }

  if (!res.ok) {
    const error = await res.text()
    throw new Error(error || res.statusText)
  }

  if (res.status === 204) {
    return {} as T
  }

  return res.json()
}

export interface User {
  id: number
  username: string
  display_name?: string
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
  owner_id: number
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
  title?: string
  content: string
  note_type?: string
  shared?: boolean
}

export const api = {
  auth: {
    login: async (username: string, password: string) => {
      const resp = await request<LoginResponse>('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ username, password }),
      })
      setToken(resp.token)
      return resp
    },
    logout: () => {
      clearToken()
    },
    me: () => request<{ user: User }>('/auth/me'),
  },

  topics: {
    list: (params?: { status?: string; priority?: number; limit?: number; offset?: number }) => {
      const query = new URLSearchParams()
      if (params?.status) query.set('status', params.status)
      if (params?.priority !== undefined) query.set('priority', params.priority.toString())
      if (params?.limit) query.set('limit', params.limit.toString())
      if (params?.offset) query.set('offset', params.offset.toString())
      return request<TopicsResponse>(`/topics?${query}`)
    },
    create: (data: CreateTopicRequest) =>
      request<{ topic: Topic }>('/topics', {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    get: (id: number) => request<{ topic: Topic }>(`/topics/${id}`),
    update: (id: number, data: Partial<Topic>) =>
      request<{ topic: Topic }>(`/topics/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data),
      }),
    delete: (id: number) => request<void>(`/topics/${id}`, { method: 'DELETE' }),
    markDiscussed: (id: number) =>
      request<{ topic: Topic }>(`/topics/${id}/discuss`, { method: 'POST' }),
  },

  events: {
    list: (params: { from: string; to: string; owner_id?: number; shared_only?: boolean }) => {
      const query = new URLSearchParams()
      query.set('from', params.from)
      query.set('to', params.to)
      if (params.owner_id) query.set('owner_id', params.owner_id.toString())
      if (params.shared_only) query.set('shared_only', 'true')
      return request<EventsResponse>(`/events?${query}`)
    },
    create: (data: CreateEventRequest) =>
      request<{ event: CalendarEvent }>('/events', {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    get: (id: number) => request<{ event: CalendarEvent }>(`/events/${id}`),
    update: (id: number, data: Partial<CalendarEvent>) =>
      request<{ event: CalendarEvent }>(`/events/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data),
      }),
    delete: (id: number) => request<void>(`/events/${id}`, { method: 'DELETE' }),
  },

  notes: {
    list: (params?: { note_type?: string; pinned_only?: boolean; limit?: number; offset?: number }) => {
      const query = new URLSearchParams()
      if (params?.note_type) query.set('note_type', params.note_type)
      if (params?.pinned_only) query.set('pinned_only', 'true')
      if (params?.limit) query.set('limit', params.limit.toString())
      if (params?.offset) query.set('offset', params.offset.toString())
      return request<NotesResponse>(`/notes?${query}`)
    },
    create: (data: CreateNoteRequest) =>
      request<{ note: Note }>('/notes', {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    get: (id: number) => request<{ note: Note }>(`/notes/${id}`),
    update: (id: number, data: Partial<Note>) =>
      request<{ note: Note }>(`/notes/${id}`, {
        method: 'PUT',
        body: JSON.stringify(data),
      }),
    delete: (id: number) => request<void>(`/notes/${id}`, { method: 'DELETE' }),
    togglePin: (id: number) => request<{ note: Note }>(`/notes/${id}/pin`, { method: 'POST' }),
  },
}
