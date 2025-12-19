export interface User {
  id: number
  username: string
  telegram_id?: number
  created_at: string
  updated_at: string
}

export interface Topic {
  id: number
  title: string
  description?: string
  priority: 0 | 1 | 2 | 3 | 4 | 5
  category_id?: number
  category?: Category
  status: 'pending' | 'discussed' | 'archived'
  tags?: Tag[]
  created_by: number
  created_at: string
  updated_at: string
  discussed_at?: string
}

export interface Category {
  id: number
  name: string
  color: string
  created_at: string
}

export interface Tag {
  id: number
  name: string
}

export interface Event {
  id: number
  title: string
  description?: string
  start_at: string
  end_at?: string
  all_day: boolean
  owner_id: number
  owner?: User
  shared: boolean
  color: string
  recurrence?: string
  created_at: string
  updated_at: string
}

export interface Note {
  id: number
  title?: string
  content: string
  note_type: 'general' | 'rule' | 'thought' | 'resource' | 'credential'
  is_pinned: boolean
  created_by: number
  creator?: User
  shared: boolean
  created_at: string
  updated_at: string
}

export interface Favorite {
  id: number
  user_id: number
  item_type: string
  item_name: string
  is_favorite: boolean
  notes?: string
  created_at: string
}

export interface MediaSession {
  id: number
  media_type: 'youtube' | 'spotify' | 'local'
  media_url: string
  media_title?: string
  host_user_id: number
  host?: User
  current_position_ms: number
  is_playing: boolean
  participants?: User[]
  created_at: string
  updated_at: string
}
