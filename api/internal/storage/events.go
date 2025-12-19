package storage

import (
	"context"
	"time"

	"github.com/zemdenalex/dna-relations/api/internal/models"
)

func CreateEvent(ctx context.Context, e *models.Event) error {
	return DB.QueryRow(ctx, `
		INSERT INTO events (title, description, start_at, end_at, all_day, owner_id, shared, color)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`, e.Title, e.Description, e.StartAt, e.EndAt, e.AllDay, e.OwnerID, e.Shared, e.Color).
		Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
}

func GetEvents(ctx context.Context, from, to time.Time, ownerID *int, sharedOnly bool) ([]models.Event, error) {
	query := `
		SELECT id, title, description, start_at, end_at, all_day, owner_id, shared, color, created_at, updated_at
		FROM events
		WHERE start_at >= $1 AND start_at <= $2
	`
	args := []interface{}{from, to}

	if sharedOnly {
		query += " AND shared = true"
	}

	if ownerID != nil {
		query += " AND owner_id = $3"
		args = append(args, *ownerID)
	}

	query += " ORDER BY start_at ASC"

	rows, err := DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		if err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.StartAt, &e.EndAt, &e.AllDay, &e.OwnerID, &e.Shared, &e.Color, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func GetEvent(ctx context.Context, id int) (*models.Event, error) {
	var e models.Event
	err := DB.QueryRow(ctx, `
		SELECT id, title, description, start_at, end_at, all_day, owner_id, shared, color, created_at, updated_at
		FROM events WHERE id = $1
	`, id).Scan(&e.ID, &e.Title, &e.Description, &e.StartAt, &e.EndAt, &e.AllDay, &e.OwnerID, &e.Shared, &e.Color, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func UpdateEvent(ctx context.Context, id int, title, description *string, startAt, endAt *time.Time, allDay, shared *bool, color *string) (*models.Event, error) {
	_, err := DB.Exec(ctx, `
		UPDATE events SET
			title = COALESCE($2, title),
			description = COALESCE($3, description),
			start_at = COALESCE($4, start_at),
			end_at = COALESCE($5, end_at),
			all_day = COALESCE($6, all_day),
			shared = COALESCE($7, shared),
			color = COALESCE($8, color),
			updated_at = NOW()
		WHERE id = $1
	`, id, title, description, startAt, endAt, allDay, shared, color)
	if err != nil {
		return nil, err
	}
	return GetEvent(ctx, id)
}

func DeleteEvent(ctx context.Context, id int) error {
	_, err := DB.Exec(ctx, `DELETE FROM events WHERE id = $1`, id)
	return err
}
