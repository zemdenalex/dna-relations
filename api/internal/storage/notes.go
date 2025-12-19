package storage

import (
	"context"
	"fmt"

	"github.com/zemdenalex/dna-relations/api/internal/models"
)

func CreateNote(ctx context.Context, n *models.Note) error {
	return DB.QueryRow(ctx, `
		INSERT INTO notes (title, content, note_type, is_pinned, created_by, shared)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`, n.Title, n.Content, n.NoteType, n.IsPinned, n.CreatedBy, n.Shared).
		Scan(&n.ID, &n.CreatedAt, &n.UpdatedAt)
}

func GetNotes(ctx context.Context, noteType *string, pinnedOnly bool, limit, offset int) ([]models.Note, int, error) {
	baseWhere := "1=1"
	args := []interface{}{}
	argIdx := 1

	if noteType != nil && *noteType != "" {
		baseWhere += fmt.Sprintf(" AND note_type = $%d", argIdx)
		args = append(args, *noteType)
		argIdx++
	}

	if pinnedOnly {
		baseWhere += " AND is_pinned = true"
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM notes WHERE " + baseWhere
	if err := DB.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT id, title, content, note_type, is_pinned, created_by, shared, created_at, updated_at
		FROM notes WHERE %s
		ORDER BY is_pinned DESC, updated_at DESC
		LIMIT $%d OFFSET $%d
	`, baseWhere, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := DB.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var n models.Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.NoteType, &n.IsPinned, &n.CreatedBy, &n.Shared, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, 0, err
		}
		notes = append(notes, n)
	}

	return notes, total, nil
}

func GetNote(ctx context.Context, id int) (*models.Note, error) {
	var n models.Note
	err := DB.QueryRow(ctx, `
		SELECT id, title, content, note_type, is_pinned, created_by, shared, created_at, updated_at
		FROM notes WHERE id = $1
	`, id).Scan(&n.ID, &n.Title, &n.Content, &n.NoteType, &n.IsPinned, &n.CreatedBy, &n.Shared, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func UpdateNote(ctx context.Context, id int, title, content, noteType *string, isPinned, shared *bool) (*models.Note, error) {
	_, err := DB.Exec(ctx, `
		UPDATE notes SET
			title = COALESCE($2, title),
			content = COALESCE($3, content),
			note_type = COALESCE($4, note_type),
			is_pinned = COALESCE($5, is_pinned),
			shared = COALESCE($6, shared),
			updated_at = NOW()
		WHERE id = $1
	`, id, title, content, noteType, isPinned, shared)
	if err != nil {
		return nil, err
	}
	return GetNote(ctx, id)
}

func ToggleNotePin(ctx context.Context, id int) (*models.Note, error) {
	_, err := DB.Exec(ctx, `
		UPDATE notes SET is_pinned = NOT is_pinned, updated_at = NOW() WHERE id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	return GetNote(ctx, id)
}

func DeleteNote(ctx context.Context, id int) error {
	_, err := DB.Exec(ctx, `DELETE FROM notes WHERE id = $1`, id)
	return err
}
