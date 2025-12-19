const API_BASE = import.meta.env.VITE_API_BASE || '/api/v1'

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

async function request<T>(
  method: string,
  path: string,
  body?: unknown
): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }

  if (authToken) {
    headers['Authorization'] = `Bearer ${authToken}`
  }

  const res = await fetch(`${API_BASE}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  })

  if (res.status === 401) {
    clearToken()
    window.location.href = '/miniapp/login'
    throw new Error('Unauthorized')
  }

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: 'Unknown error' }))
    throw new Error(err.error || err.message || 'Request failed')
  }

  if (res.status === 204) {
    return {} as T
  }

  return res.json()
}

export interface User {
  id: number
  username: string
  telegram_id?: number
}

export interface Topic {
  id: number
  title: string
  description?: string
  priority: number
  status: 'pending' | 'discussed' | 'archived'
  created_by: number
  created_at: string
  discussed_at?: string
}

export interface Event {
  id: number
  title: string
  description?: string
  start_at: string
  end_at?: string
  all_day: boolean
  owner_id: number
  shared: boolean
  color: string
}

export interface Note {
  id: number
  title?: string
  content: string
  note_type: string
  is_pinned: boolean
  created_by: number
  shared: boolean
  created_at: string
}

export const api = {
  auth: {
    login: async (username: string, password: string) => {
      const resp = await request<{ token: string; user: User }>('POST', '/auth/login', {
        username,
        password,
      })
      setToken(resp.token)
      return resp
    },
    me: () => request<{ user: User }>('GET', '/auth/me'),
    logout: () => {
      clearToken()
    },
  },

  topics: {
    list: (params?: { status?: string; priority?: number; limit?: number; offset?: number }) => {
      const query = new URLSearchParams()
      if (params?.status) query.set('status', params.status)
      if (params?.priority !== undefined) query.set('priority', params.priority.toString())
      if (params?.limit) query.set('limit', params.limit.toString())
      if (params?.offset) query.set('offset', params.offset.toString())
      return request<{ topics: Topic[]; total: number }>('GET', `/topics?${query}`)
    },
    create: (data: { title: string; description?: string; priority: number }) =>
      request<{ topic: Topic }>('POST', '/topics', data),
    get: (id: number) => request<{ topic: Topic }>('GET', `/topics/${id}`),
    update: (id: number, data: Partial<{ title: string; description: string; priority: number; status: string }>) =>
      request<{ topic: Topic }>('PUT', `/topics/${id}`, data),
    delete: (id: number) => request<void>('DELETE', `/topics/${id}`),
    markDiscussed: (id: number) => request<{ topic: Topic }>('POST', `/topics/${id}/discuss`),
  },

  events: {
    list: (params: { from: string; to: string; owner_id?: number; shared_only?: boolean }) => {
      const query = new URLSearchParams()
      query.set('from', params.from)
      query.set('to', params.to)
      if (params.owner_id) query.set('owner_id', params.owner_id.toString())
      if (params.shared_only) query.set('shared_only', 'true')
      return request<{ events: Event[] }>('GET', `/events?${query}`)
    },
    create: (data: {
      title: string
      description?: string
      start_at: string
      end_at?: string
      all_day?: boolean
      shared?: boolean
      color?: string
    }) => request<{ event: Event }>('POST', '/events', data),
    get: (id: number) => request<{ event: Event }>('GET', `/events/${id}`),
    update: (id: number, data: Partial<Event>) => request<{ event: Event }>('PUT', `/events/${id}`, data),
    delete: (id: number) => request<void>('DELETE', `/events/${id}`),
  },

  notes: {
    list: (params?: { note_type?: string; pinned_only?: boolean; limit?: number; offset?: number }) => {
      const query = new URLSearchParams()
      if (params?.note_type) query.set('note_type', params.note_type)
      if (params?.pinned_only) query.set('pinned_only', 'true')
      if (params?.limit) query.set('limit', params.limit.toString())
      if (params?.offset) query.set('offset', params.offset.toString())
      return request<{ notes: Note[]; total: number }>('GET', `/notes?${query}`)
    },
    create: (data: { title?: string; content: string; note_type: string; shared?: boolean }) =>
      request<{ note: Note }>('POST', '/notes', data),
    get: (id: number) => request<{ note: Note }>('GET', `/notes/${id}`),
    update: (id: number, data: Partial<Note>) => request<{ note: Note }>('PUT', `/notes/${id}`, data),
    delete: (id: number) => request<void>('DELETE', `/notes/${id}`),
    togglePin: (id: number) => request<{ note: Note }>('POST', `/notes/${id}/pin`),
  },
}
