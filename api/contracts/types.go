package contracts

import "time"

type User struct {
	ID         int        `json:"id"`
	Username   string     `json:"username"`
	TelegramID *int64     `json:"telegram_id,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type Topic struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Priority    int        `json:"priority"`
	CategoryID  *int       `json:"category_id,omitempty"`
	Category    *Category  `json:"category,omitempty"`
	Status      string     `json:"status"`
	Tags        []Tag      `json:"tags,omitempty"`
	CreatedBy   int        `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DiscussedAt *time.Time `json:"discussed_at,omitempty"`
}

type Category struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Event struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	StartAt     time.Time  `json:"start_at"`
	EndAt       *time.Time `json:"end_at,omitempty"`
	AllDay      bool       `json:"all_day"`
	OwnerID     int        `json:"owner_id"`
	Owner       *User      `json:"owner,omitempty"`
	Shared      bool       `json:"shared"`
	Color       string     `json:"color"`
	Recurrence  string     `json:"recurrence,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title,omitempty"`
	Content   string    `json:"content"`
	NoteType  string    `json:"note_type"`
	IsPinned  bool      `json:"is_pinned"`
	CreatedBy int       `json:"created_by"`
	Creator   *User     `json:"creator,omitempty"`
	Shared    bool      `json:"shared"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Favorite struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	ItemType   string    `json:"item_type"`
	ItemName   string    `json:"item_name"`
	IsFavorite bool      `json:"is_favorite"`
	Notes      string    `json:"notes,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type MediaSession struct {
	ID                int       `json:"id"`
	MediaType         string    `json:"media_type"`
	MediaURL          string    `json:"media_url"`
	MediaTitle        string    `json:"media_title,omitempty"`
	HostUserID        int       `json:"host_user_id"`
	Host              *User     `json:"host,omitempty"`
	CurrentPositionMs int64     `json:"current_position_ms"`
	IsPlaying         bool      `json:"is_playing"`
	Participants      []User    `json:"participants,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
