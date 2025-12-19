package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/zemdenalex/dna-relations/api/internal/models"
	"github.com/zemdenalex/dna-relations/api/internal/storage"
)

type CreateTopicRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	CategoryID  *int   `json:"category_id"`
}

type UpdateTopicRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Priority    *int    `json:"priority"`
	Status      *string `json:"status"`
}

func ListTopics(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status == "" {
		status = "pending"
	}

	var priority *int
	if p := r.URL.Query().Get("priority"); p != "" {
		if v, err := strconv.Atoi(p); err == nil {
			priority = &v
		}
	}

	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}

	offset := 0
	if o := r.URL.Query().Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	topics, total, err := storage.GetTopics(r.Context(), status, priority, limit, offset)
	if err != nil {
		jsonError(w, "failed to fetch topics", http.StatusInternalServerError)
		return
	}

	if topics == nil {
		topics = []models.Topic{}
	}

	jsonResponse(w, map[string]interface{}{
		"topics": topics,
		"total":  total,
	})
}

func CreateTopic(w http.ResponseWriter, r *http.Request) {
	var req CreateTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		jsonError(w, "title is required", http.StatusBadRequest)
		return
	}

	topic := &models.Topic{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		CategoryID:  req.CategoryID,
		Status:      "pending",
		CreatedBy:   GetUserID(r.Context()),
	}

	if err := storage.CreateTopic(r.Context(), topic); err != nil {
		jsonError(w, "failed to create topic", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	jsonResponse(w, map[string]interface{}{"topic": topic})
}

func GetTopic(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	topic, err := storage.GetTopic(r.Context(), id)
	if err != nil {
		jsonError(w, "topic not found", http.StatusNotFound)
		return
	}

	jsonResponse(w, map[string]interface{}{"topic": topic})
}

func UpdateTopic(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req UpdateTopicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	topic, err := storage.UpdateTopic(r.Context(), id, req.Title, req.Description, req.Priority, req.Status)
	if err != nil {
		jsonError(w, "failed to update topic", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]interface{}{"topic": topic})
}

func DeleteTopic(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := storage.DeleteTopic(r.Context(), id); err != nil {
		jsonError(w, "failed to delete topic", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func MarkTopicDiscussed(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	topic, err := storage.MarkTopicDiscussed(r.Context(), id)
	if err != nil {
		jsonError(w, "failed to mark topic as discussed", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]interface{}{"topic": topic})
}
