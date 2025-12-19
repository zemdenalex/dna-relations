package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/zemdenalex/dna-relations/api/internal/models"
	"github.com/zemdenalex/dna-relations/api/internal/storage"
)

type CreateEventRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartAt     time.Time  `json:"start_at"`
	EndAt       *time.Time `json:"end_at"`
	AllDay      bool       `json:"all_day"`
	Shared      bool       `json:"shared"`
	Color       string     `json:"color"`
}

type UpdateEventRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	StartAt     *time.Time `json:"start_at"`
	EndAt       *time.Time `json:"end_at"`
	AllDay      *bool      `json:"all_day"`
	Shared      *bool      `json:"shared"`
	Color       *string    `json:"color"`
}

func ListEvents(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	var from, to time.Time
	var err error

	if fromStr != "" {
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			from, err = time.Parse("2006-01-02", fromStr)
		}
		if err != nil {
			jsonError(w, "invalid from date", http.StatusBadRequest)
			return
		}
	} else {
		from = time.Now().AddDate(0, 0, -7)
	}

	if toStr != "" {
		to, err = time.Parse(time.RFC3339, toStr)
		if err != nil {
			to, err = time.Parse("2006-01-02", toStr)
		}
		if err != nil {
			jsonError(w, "invalid to date", http.StatusBadRequest)
			return
		}
	} else {
		to = time.Now().AddDate(0, 1, 0)
	}

	var ownerID *int
	if o := r.URL.Query().Get("owner_id"); o != "" {
		if v, err := strconv.Atoi(o); err == nil {
			ownerID = &v
		}
	}

	sharedOnly := r.URL.Query().Get("shared_only") == "true"

	events, err := storage.GetEvents(r.Context(), from, to, ownerID, sharedOnly)
	if err != nil {
		jsonError(w, "failed to fetch events", http.StatusInternalServerError)
		return
	}

	if events == nil {
		events = []models.Event{}
	}

	jsonResponse(w, map[string]interface{}{"events": events})
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		jsonError(w, "title is required", http.StatusBadRequest)
		return
	}

	if req.StartAt.IsZero() {
		jsonError(w, "start_at is required", http.StatusBadRequest)
		return
	}

	if req.Color == "" {
		req.Color = "#6366f1"
	}

	event := &models.Event{
		Title:       req.Title,
		Description: req.Description,
		StartAt:     req.StartAt,
		EndAt:       req.EndAt,
		AllDay:      req.AllDay,
		OwnerID:     GetUserID(r.Context()),
		Shared:      req.Shared,
		Color:       req.Color,
	}

	if err := storage.CreateEvent(r.Context(), event); err != nil {
		jsonError(w, "failed to create event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	jsonResponse(w, map[string]interface{}{"event": event})
}

func GetEvent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	event, err := storage.GetEvent(r.Context(), id)
	if err != nil {
		jsonError(w, "event not found", http.StatusNotFound)
		return
	}

	jsonResponse(w, map[string]interface{}{"event": event})
}

func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	event, err := storage.UpdateEvent(r.Context(), id, req.Title, req.Description, req.StartAt, req.EndAt, req.AllDay, req.Shared, req.Color)
	if err != nil {
		jsonError(w, "failed to update event", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]interface{}{"event": event})
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := storage.DeleteEvent(r.Context(), id); err != nil {
		jsonError(w, "failed to delete event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
