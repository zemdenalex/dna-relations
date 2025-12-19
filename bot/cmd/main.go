package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"

	"github.com/zemdenalex/dna-relations/bot/internal/handlers"
)

func main() {
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is required")
	}

	apiBase := os.Getenv("API_BASE")
	if apiBase == "" {
		apiBase = "http://localhost:9000"
	}

	webappURL := os.Getenv("WEBAPP_URL")
	botPort := os.Getenv("BOT_PORT")
	if botPort == "" {
		botPort = "3002"
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	h := handlers.New(bot, apiBase, webappURL)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})
		log.Printf("Health check server on :%s", botPort)
		http.ListenAndServe(":"+botPort, nil)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("Shutting down bot...")
		bot.StopReceivingUpdates()
		os.Exit(0)
	}()

	for update := range updates {
		if update.Message != nil {
			h.HandleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			h.HandleCallback(update.CallbackQuery)
		}
	}
}
