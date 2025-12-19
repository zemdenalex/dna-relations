package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/zemdenalex/dna-relations/api/internal/handlers"
	"github.com/zemdenalex/dna-relations/api/internal/storage"
)

func main() {
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "9000"
	}

	ctx := context.Background()
	if err := storage.Connect(ctx); err != nil {
		log.Printf("Warning: Database connection failed: %v", err)
		log.Println("API will run but database operations will fail")
	} else {
		log.Println("Connected to database")
		defer storage.Close()
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/login", handlers.Login)

		r.Group(func(r chi.Router) {
			r.Use(handlers.AuthMiddleware)

			r.Get("/auth/me", handlers.Me)

			r.Route("/topics", func(r chi.Router) {
				r.Get("/", handlers.ListTopics)
				r.Post("/", handlers.CreateTopic)
				r.Get("/{id}", handlers.GetTopic)
				r.Put("/{id}", handlers.UpdateTopic)
				r.Delete("/{id}", handlers.DeleteTopic)
				r.Post("/{id}/discuss", handlers.MarkTopicDiscussed)
			})

			r.Route("/events", func(r chi.Router) {
				r.Get("/", handlers.ListEvents)
				r.Post("/", handlers.CreateEvent)
				r.Get("/{id}", handlers.GetEvent)
				r.Put("/{id}", handlers.UpdateEvent)
				r.Delete("/{id}", handlers.DeleteEvent)
			})

			r.Route("/notes", func(r chi.Router) {
				r.Get("/", handlers.ListNotes)
				r.Post("/", handlers.CreateNote)
				r.Get("/{id}", handlers.GetNote)
				r.Put("/{id}", handlers.UpdateNote)
				r.Delete("/{id}", handlers.DeleteNote)
				r.Post("/{id}/pin", handlers.ToggleNotePin)
			})
		})
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("Starting API server on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
