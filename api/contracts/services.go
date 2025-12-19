package contracts

import (
	"context"
	"time"
)

// NOTE: These contracts are designed for use with github.com/seniorGolang/tg
// To make them work with tg, move types to a separate "models" package
// and update return types to use models.Topic, models.Event, etc.
// For now, the API uses chi router with stub handlers.

// @tg version=1.0.0
// @tg title=DNA Relations API
// @tg servers=http://localhost:9000;local
// @tg security=bearer

// @tg jsonRPC-server log metrics
// @tg http-prefix=api/v1
// @tg swaggerTags=auth
type AuthService interface {
	// @tg http-method=POST
	// @tg http-path=/auth/login
	Login(ctx context.Context, username string, password string) (token string, user User, err error)

	// @tg http-method=GET
	// @tg http-path=/auth/me
	Me(ctx context.Context) (user User, err error)
}

// @tg jsonRPC-server log metrics
// @tg http-prefix=api/v1
// @tg swaggerTags=topics
type TopicService interface {
	// @tg http-method=POST
	// @tg http-path=/topics
	Create(ctx context.Context, title string, description string, priority int, categoryID *int, tags []string) (topic Topic, err error)

	// @tg http-method=GET
	// @tg http-path=/topics
	List(ctx context.Context, status string, categoryID *int, priority *int, limit int, offset int) (topics []Topic, total int, err error)

	// @tg http-method=GET
	// @tg http-path=/topics/:id
	// @tg http-args=id|id
	Get(ctx context.Context, id int) (topic Topic, err error)

	// @tg http-method=PUT
	// @tg http-path=/topics/:id
	// @tg http-args=id|id
	Update(ctx context.Context, id int, title *string, description *string, priority *int, categoryID *int, status *string) (topic Topic, err error)

	// @tg http-method=DELETE
	// @tg http-path=/topics/:id
	// @tg http-args=id|id
	Delete(ctx context.Context, id int) (err error)

	// @tg http-method=POST
	// @tg http-path=/topics/:id/discuss
	// @tg http-args=id|id
	MarkDiscussed(ctx context.Context, id int) (topic Topic, err error)
}

// @tg jsonRPC-server log metrics
// @tg http-prefix=api/v1
// @tg swaggerTags=categories
type CategoryService interface {
	// @tg http-method=POST
	// @tg http-path=/categories
	Create(ctx context.Context, name string, color string) (category Category, err error)

	// @tg http-method=GET
	// @tg http-path=/categories
	List(ctx context.Context) (categories []Category, err error)

	// @tg http-method=DELETE
	// @tg http-path=/categories/:id
	// @tg http-args=id|id
	Delete(ctx context.Context, id int) (err error)
}

// @tg jsonRPC-server log metrics
// @tg http-prefix=api/v1
// @tg swaggerTags=events
type EventService interface {
	// @tg http-method=POST
	// @tg http-path=/events
	Create(ctx context.Context, title string, description string, startAt time.Time, endAt *time.Time, allDay bool, shared bool, color string) (event Event, err error)

	// @tg http-method=GET
	// @tg http-path=/events
	List(ctx context.Context, from time.Time, to time.Time, ownerID *int, sharedOnly bool) (events []Event, err error)

	// @tg http-method=GET
	// @tg http-path=/events/:id
	// @tg http-args=id|id
	Get(ctx context.Context, id int) (event Event, err error)

	// @tg http-method=PUT
	// @tg http-path=/events/:id
	// @tg http-args=id|id
	Update(ctx context.Context, id int, title *string, description *string, startAt *time.Time, endAt *time.Time, allDay *bool, shared *bool, color *string) (event Event, err error)

	// @tg http-method=DELETE
	// @tg http-path=/events/:id
	// @tg http-args=id|id
	Delete(ctx context.Context, id int) (err error)
}

// @tg jsonRPC-server log metrics
// @tg http-prefix=api/v1
// @tg swaggerTags=notes
type NoteService interface {
	// @tg http-method=POST
	// @tg http-path=/notes
	Create(ctx context.Context, title string, content string, noteType string, shared bool) (note Note, err error)

	// @tg http-method=GET
	// @tg http-path=/notes
	List(ctx context.Context, noteType *string, pinnedOnly bool, limit int, offset int) (notes []Note, total int, err error)

	// @tg http-method=GET
	// @tg http-path=/notes/:id
	// @tg http-args=id|id
	Get(ctx context.Context, id int) (note Note, err error)

	// @tg http-method=PUT
	// @tg http-path=/notes/:id
	// @tg http-args=id|id
	Update(ctx context.Context, id int, title *string, content *string, noteType *string, isPinned *bool, shared *bool) (note Note, err error)

	// @tg http-method=DELETE
	// @tg http-path=/notes/:id
	// @tg http-args=id|id
	Delete(ctx context.Context, id int) (err error)

	// @tg http-method=POST
	// @tg http-path=/notes/:id/pin
	// @tg http-args=id|id
	TogglePin(ctx context.Context, id int) (note Note, err error)
}

// @tg jsonRPC-server log metrics
// @tg http-prefix=api/v1
// @tg swaggerTags=favorites
type FavoriteService interface {
	// @tg http-method=POST
	// @tg http-path=/favorites
	Create(ctx context.Context, itemType string, itemName string, isFavorite bool, notes string) (favorite Favorite, err error)

	// @tg http-method=GET
	// @tg http-path=/favorites
	List(ctx context.Context, itemType *string, favoritesOnly bool) (favorites []Favorite, err error)

	// @tg http-method=DELETE
	// @tg http-path=/favorites/:id
	// @tg http-args=id|id
	Delete(ctx context.Context, id int) (err error)
}

// @tg jsonRPC-server log metrics
// @tg http-prefix=api/v1
// @tg swaggerTags=media
type MediaService interface {
	// @tg http-method=POST
	// @tg http-path=/media/sessions
	CreateSession(ctx context.Context, mediaType string, mediaURL string, mediaTitle string) (session MediaSession, err error)

	// @tg http-method=GET
	// @tg http-path=/media/sessions/active
	GetActiveSession(ctx context.Context) (session *MediaSession, err error)

	// @tg http-method=POST
	// @tg http-path=/media/sessions/:id/join
	// @tg http-args=id|id
	JoinSession(ctx context.Context, id int) (session MediaSession, err error)

	// @tg http-method=POST
	// @tg http-path=/media/sessions/:id/sync
	// @tg http-args=id|id
	SyncPosition(ctx context.Context, id int, positionMs int64, isPlaying bool) (err error)

	// @tg http-method=DELETE
	// @tg http-path=/media/sessions/:id
	// @tg http-args=id|id
	EndSession(ctx context.Context, id int) (err error)
}

// @tg http-method=GET
// @tg http-path=/health
type HealthService interface {
	Check(ctx context.Context) (status string, err error)
}
