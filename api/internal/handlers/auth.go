package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zemdenalex/dna-relations/api/internal/models"
	"github.com/zemdenalex/dna-relations/api/internal/storage"
)

var jwtSecret []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-in-production"
	}
	jwtSecret = []byte(secret)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}

	denisPassword := os.Getenv("AUTH_PASSWORD_DENIS")
	nastyaPassword := os.Getenv("AUTH_PASSWORD_NASTYA")

	var userID int
	var username string

	switch {
	case req.Username == "denis" && req.Password == denisPassword:
		userID = 1
		username = "denis"
	case req.Username == "nastya" && req.Password == nastyaPassword:
		userID = 2
		username = "nastya"
	default:
		jsonError(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	user, err := storage.GetUserByID(r.Context(), userID)
	if err != nil {
		user = &models.User{ID: userID, Username: username}
	}

	token, err := generateToken(userID, username)
	if err != nil {
		jsonError(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, LoginResponse{Token: token, User: user})
}

func Me(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)
	user, err := storage.GetUserByID(r.Context(), userID)
	if err != nil {
		jsonError(w, "user not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, map[string]interface{}{"user": user})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			jsonError(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			jsonError(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		claims, err := validateToken(parts[1])
		if err != nil {
			jsonError(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = contextWithUserID(ctx, claims.UserID)
		ctx = contextWithUsername(ctx, claims.Username)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func generateToken(userID int, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func validateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
