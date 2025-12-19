package models

import "time"

type User struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	TelegramID *int64    `json:"telegram_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type Topic struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Priority    int        `json:"priority"`
	CategoryID  *int       `json:"category_id,omitempty"`
	Status      string     `json:"status"`
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

type Event struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	StartAt     time.Time  `json:"start_at"`
	EndAt       *time.Time `json:"end_at,omitempty"`
	AllDay      bool       `json:"all_day"`
	OwnerID     int        `json:"owner_id"`
	Shared      bool       `json:"shared"`
	Color       string     `json:"color"`
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
	Shared    bool      `json:"shared"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
