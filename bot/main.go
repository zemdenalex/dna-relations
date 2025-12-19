package main

import (
	"log"
	"os"
	"time"

	"dna-bot/internal/handlers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://backend:8080"
	}

	webAppURL := os.Getenv("WEB_APP_URL")
	if webAppURL == "" {
		webAppURL = "https://dna-relations.site/miniapp"
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = os.Getenv("DEBUG") == "true"
	log.Printf("Authorized on account %s", bot.Self.UserName)

	h := handlers.New(bot, apiURL, webAppURL)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		go func(update tgbotapi.Update) {
			ctx := handlers.Context{
				Update:    &update,
				Timestamp: time.Now(),
			}
			h.Handle(ctx)
		}(update)
	}
}