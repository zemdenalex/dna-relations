package handlers

import (
	"context"
	"encoding/json"
	"net/http"
)

type contextKey string

const (
	userIDKey   contextKey = "user_id"
	usernameKey contextKey = "username"
)

func contextWithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func contextWithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, usernameKey, username)
}

func GetUserID(ctx context.Context) int {
	if v := ctx.Value(userIDKey); v != nil {
		return v.(int)
	}
	return 0
}

func GetUsername(ctx context.Context) string {
	if v := ctx.Value(usernameKey); v != nil {
		return v.(string)
	}
	return ""
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
