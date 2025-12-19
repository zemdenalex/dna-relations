package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/zemdenalex/dna-relations/api/internal/models"
	"github.com/zemdenalex/dna-relations/api/internal/storage"
)

type CreateNoteRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	NoteType string `json:"note_type"`
	Shared   bool   `json:"shared"`
}

type UpdateNoteRequest struct {
	Title    *string `json:"title"`
	Content  *string `json:"content"`
	NoteType *string `json:"note_type"`
	IsPinned *bool   `json:"is_pinned"`
	Shared   *bool   `json:"shared"`
}

func ListNotes(w http.ResponseWriter, r *http.Request) {
	var noteType *string
	if nt := r.URL.Query().Get("note_type"); nt != "" {
		noteType = &nt
	}

	pinnedOnly := r.URL.Query().Get("pinned_only") == "true"

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

	notes, total, err := storage.GetNotes(r.Context(), noteType, pinnedOnly, limit, offset)
	if err != nil {
		jsonError(w, "failed to fetch notes", http.StatusInternalServerError)
		return
	}

	if notes == nil {
		notes = []models.Note{}
	}

	jsonResponse(w, map[string]interface{}{
		"notes": notes,
		"total": total,
	})
}

func CreateNote(w http.ResponseWriter, r *http.Request) {
	var req CreateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		jsonError(w, "content is required", http.StatusBadRequest)
		return
	}

	if req.NoteType == "" {
		req.NoteType = "general"
	}

	note := &models.Note{
		Title:     req.Title,
		Content:   req.Content,
		NoteType:  req.NoteType,
		IsPinned:  false,
		CreatedBy: GetUserID(r.Context()),
		Shared:    req.Shared,
	}

	if err := storage.CreateNote(r.Context(), note); err != nil {
		jsonError(w, "failed to create note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	jsonResponse(w, map[string]interface{}{"note": note})
}

func GetNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	note, err := storage.GetNote(r.Context(), id)
	if err != nil {
		jsonError(w, "note not found", http.StatusNotFound)
		return
	}

	jsonResponse(w, map[string]interface{}{"note": note})
}

func UpdateNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	note, err := storage.UpdateNote(r.Context(), id, req.Title, req.Content, req.NoteType, req.IsPinned, req.Shared)
	if err != nil {
		jsonError(w, "failed to update note", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]interface{}{"note": note})
}

func DeleteNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := storage.DeleteNote(r.Context(), id); err != nil {
		jsonError(w, "failed to delete note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func ToggleNotePin(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		jsonError(w, "invalid id", http.StatusBadRequest)
		return
	}

	note, err := storage.ToggleNotePin(r.Context(), id)
	if err != nil {
		jsonError(w, "failed to toggle pin", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]interface{}{"note": note})
}
