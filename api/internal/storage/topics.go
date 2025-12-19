package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/zemdenalex/dna-relations/api/internal/models"
)

func CreateTopic(ctx context.Context, t *models.Topic) error {
	return DB.QueryRow(ctx, `
		INSERT INTO topics (title, description, priority, category_id, status, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`, t.Title, t.Description, t.Priority, t.CategoryID, t.Status, t.CreatedBy).
		Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func GetTopics(ctx context.Context, status string, priority *int, limit, offset int) ([]models.Topic, int, error) {
	baseWhere := "1=1"
	args := []interface{}{}
	argIdx := 1

	if status != "" && status != "all" {
		baseWhere += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	if priority != nil {
		baseWhere += fmt.Sprintf(" AND priority = $%d", argIdx)
		args = append(args, *priority)
		argIdx++
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM topics WHERE " + baseWhere
	if err := DB.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT id, title, description, priority, category_id, status, created_by, created_at, updated_at, discussed_at
		FROM topics WHERE %s
		ORDER BY priority ASC, created_at DESC
		LIMIT $%d OFFSET $%d
	`, baseWhere, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := DB.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var topics []models.Topic
	for rows.Next() {
		var t models.Topic
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Priority, &t.CategoryID, &t.Status, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt, &t.DiscussedAt); err != nil {
			return nil, 0, err
		}
		topics = append(topics, t)
	}

	return topics, total, nil
}

func GetTopic(ctx context.Context, id int) (*models.Topic, error) {
	var t models.Topic
	err := DB.QueryRow(ctx, `
		SELECT id, title, description, priority, category_id, status, created_by, created_at, updated_at, discussed_at
		FROM topics WHERE id = $1
	`, id).Scan(&t.ID, &t.Title, &t.Description, &t.Priority, &t.CategoryID, &t.Status, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt, &t.DiscussedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func UpdateTopic(ctx context.Context, id int, title, description *string, priority *int, status *string) (*models.Topic, error) {
	_, err := DB.Exec(ctx, `
		UPDATE topics SET
			title = COALESCE($2, title),
			description = COALESCE($3, description),
			priority = COALESCE($4, priority),
			status = COALESCE($5, status),
			updated_at = NOW()
		WHERE id = $1
	`, id, title, description, priority, status)
	if err != nil {
		return nil, err
	}
	return GetTopic(ctx, id)
}

func MarkTopicDiscussed(ctx context.Context, id int) (*models.Topic, error) {
	now := time.Now()
	_, err := DB.Exec(ctx, `
		UPDATE topics SET status = 'discussed', discussed_at = $2, updated_at = NOW() WHERE id = $1
	`, id, now)
	if err != nil {
		return nil, err
	}
	return GetTopic(ctx, id)
}

func DeleteTopic(ctx context.Context, id int) error {
	_, err := DB.Exec(ctx, `DELETE FROM topics WHERE id = $1`, id)
	return err
}
