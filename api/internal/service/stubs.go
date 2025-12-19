package service

import (
	"context"
	"errors"
	"time"

	"github.com/zemdenalex/dna-relations/api/contracts"
)

var ErrNotImplemented = errors.New("not implemented")

type AuthServiceImpl struct{}

func (s *AuthServiceImpl) Login(ctx context.Context, username string, password string) (token string, user contracts.User, err error) {
	return "", contracts.User{}, ErrNotImplemented
}

func (s *AuthServiceImpl) Me(ctx context.Context) (user contracts.User, err error) {
	return contracts.User{}, ErrNotImplemented
}

type TopicServiceImpl struct{}

func (s *TopicServiceImpl) Create(ctx context.Context, title string, description string, priority int, categoryID *int, tags []string) (topic contracts.Topic, err error) {
	return contracts.Topic{}, ErrNotImplemented
}

func (s *TopicServiceImpl) List(ctx context.Context, status string, categoryID *int, priority *int, limit int, offset int) (topics []contracts.Topic, total int, err error) {
	return nil, 0, ErrNotImplemented
}

func (s *TopicServiceImpl) Get(ctx context.Context, id int) (topic contracts.Topic, err error) {
	return contracts.Topic{}, ErrNotImplemented
}

func (s *TopicServiceImpl) Update(ctx context.Context, id int, title *string, description *string, priority *int, categoryID *int, status *string) (topic contracts.Topic, err error) {
	return contracts.Topic{}, ErrNotImplemented
}

func (s *TopicServiceImpl) Delete(ctx context.Context, id int) (err error) {
	return ErrNotImplemented
}

func (s *TopicServiceImpl) MarkDiscussed(ctx context.Context, id int) (topic contracts.Topic, err error) {
	return contracts.Topic{}, ErrNotImplemented
}

type CategoryServiceImpl struct{}

func (s *CategoryServiceImpl) Create(ctx context.Context, name string, color string) (category contracts.Category, err error) {
	return contracts.Category{}, ErrNotImplemented
}

func (s *CategoryServiceImpl) List(ctx context.Context) (categories []contracts.Category, err error) {
	return nil, ErrNotImplemented
}

func (s *CategoryServiceImpl) Delete(ctx context.Context, id int) (err error) {
	return ErrNotImplemented
}

type EventServiceImpl struct{}

func (s *EventServiceImpl) Create(ctx context.Context, title string, description string, startAt time.Time, endAt *time.Time, allDay bool, shared bool, color string) (event contracts.Event, err error) {
	return contracts.Event{}, ErrNotImplemented
}

func (s *EventServiceImpl) List(ctx context.Context, from time.Time, to time.Time, ownerID *int, sharedOnly bool) (events []contracts.Event, err error) {
	return nil, ErrNotImplemented
}

func (s *EventServiceImpl) Get(ctx context.Context, id int) (event contracts.Event, err error) {
	return contracts.Event{}, ErrNotImplemented
}

func (s *EventServiceImpl) Update(ctx context.Context, id int, title *string, description *string, startAt *time.Time, endAt *time.Time, allDay *bool, shared *bool, color *string) (event contracts.Event, err error) {
	return contracts.Event{}, ErrNotImplemented
}

func (s *EventServiceImpl) Delete(ctx context.Context, id int) (err error) {
	return ErrNotImplemented
}

type NoteServiceImpl struct{}

func (s *NoteServiceImpl) Create(ctx context.Context, title string, content string, noteType string, shared bool) (note contracts.Note, err error) {
	return contracts.Note{}, ErrNotImplemented
}

func (s *NoteServiceImpl) List(ctx context.Context, noteType *string, pinnedOnly bool, limit int, offset int) (notes []contracts.Note, total int, err error) {
	return nil, 0, ErrNotImplemented
}

func (s *NoteServiceImpl) Get(ctx context.Context, id int) (note contracts.Note, err error) {
	return contracts.Note{}, ErrNotImplemented
}

func (s *NoteServiceImpl) Update(ctx context.Context, id int, title *string, content *string, noteType *string, isPinned *bool, shared *bool) (note contracts.Note, err error) {
	return contracts.Note{}, ErrNotImplemented
}

func (s *NoteServiceImpl) Delete(ctx context.Context, id int) (err error) {
	return ErrNotImplemented
}

func (s *NoteServiceImpl) TogglePin(ctx context.Context, id int) (note contracts.Note, err error) {
	return contracts.Note{}, ErrNotImplemented
}

type FavoriteServiceImpl struct{}

func (s *FavoriteServiceImpl) Create(ctx context.Context, itemType string, itemName string, isFavorite bool, notes string) (favorite contracts.Favorite, err error) {
	return contracts.Favorite{}, ErrNotImplemented
}

func (s *FavoriteServiceImpl) List(ctx context.Context, itemType *string, favoritesOnly bool) (favorites []contracts.Favorite, err error) {
	return nil, ErrNotImplemented
}

func (s *FavoriteServiceImpl) Delete(ctx context.Context, id int) (err error) {
	return ErrNotImplemented
}

type MediaServiceImpl struct{}

func (s *MediaServiceImpl) CreateSession(ctx context.Context, mediaType string, mediaURL string, mediaTitle string) (session contracts.MediaSession, err error) {
	return contracts.MediaSession{}, ErrNotImplemented
}

func (s *MediaServiceImpl) GetActiveSession(ctx context.Context) (session *contracts.MediaSession, err error) {
	return nil, ErrNotImplemented
}

func (s *MediaServiceImpl) JoinSession(ctx context.Context, id int) (session contracts.MediaSession, err error) {
	return contracts.MediaSession{}, ErrNotImplemented
}

func (s *MediaServiceImpl) SyncPosition(ctx context.Context, id int, positionMs int64, isPlaying bool) (err error) {
	return ErrNotImplemented
}

func (s *MediaServiceImpl) EndSession(ctx context.Context, id int) (err error) {
	return ErrNotImplemented
}

type HealthServiceImpl struct{}

func (s *HealthServiceImpl) Check(ctx context.Context) (status string, err error) {
	return "ok", nil
}
