package storage

import (
	"context"

	"github.com/zemdenalex/dna-relations/api/internal/models"
)

func GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var u models.User
	err := DB.QueryRow(ctx, `
		SELECT id, username, telegram_id, created_at FROM users WHERE id = $1
	`, id).Scan(&u.ID, &u.Username, &u.TelegramID, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var u models.User
	err := DB.QueryRow(ctx, `
		SELECT id, username, telegram_id, created_at FROM users WHERE username = $1
	`, username).Scan(&u.ID, &u.Username, &u.TelegramID, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
	var u models.User
	err := DB.QueryRow(ctx, `
		SELECT id, username, telegram_id, created_at FROM users WHERE telegram_id = $1
	`, telegramID).Scan(&u.ID, &u.Username, &u.TelegramID, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func CreateUser(ctx context.Context, username string, passwordHash string, telegramID *int64) (*models.User, error) {
	var u models.User
	err := DB.QueryRow(ctx, `
		INSERT INTO users (username, password_hash, telegram_id)
		VALUES ($1, $2, $3)
		RETURNING id, username, telegram_id, created_at
	`, username, passwordHash, telegramID).Scan(&u.ID, &u.Username, &u.TelegramID, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func UpdateUserTelegramID(ctx context.Context, userID int, telegramID int64) error {
	_, err := DB.Exec(ctx, `UPDATE users SET telegram_id = $2 WHERE id = $1`, userID, telegramID)
	return err
}
