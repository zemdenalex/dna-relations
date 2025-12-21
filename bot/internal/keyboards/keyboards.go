package keyboards

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MainMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Темы"),
			tgbotapi.NewKeyboardButton("Календарь"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Заметки"),
			tgbotapi.NewKeyboardButton("Приложение"),
		),
	)
}

func MainMenuWithWebApp(webappURL string) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				tgbotapi.NewKeyboardButton("Темы"),
				tgbotapi.NewKeyboardButton("Календарь"),
			},
			{
				tgbotapi.NewKeyboardButton("Заметки"),
				tgbotapi.NewKeyboardButton("Приложение"),
			},
		},
		ResizeKeyboard: true,
	}
}

func TopicsMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Все", "topics_list"),
			tgbotapi.NewInlineKeyboardButtonData("Добавить", "topics_add"),
		),
	)
}

type TopicInfo struct {
	ID       int
	Title    string
	Priority int
	Status   string
}

func TopicsMenuExtended(topics []TopicInfo) tgbotapi.InlineKeyboardMarkup {
	rows := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("Ожидают", "topics_list"),
			tgbotapi.NewInlineKeyboardButtonData("Обсуждённые", "topics_discussed"),
			tgbotapi.NewInlineKeyboardButtonData("+", "topics_add"),
		},
	}

	for _, t := range topics {
		if t.Status == "pending" {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("✓ %s", truncate(t.Title, 25)),
					fmt.Sprintf("topic_discuss:%d", t.ID),
				),
			))
		}
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func CalendarMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сегодня", "calendar_today"),
			tgbotapi.NewInlineKeyboardButtonData("Неделя", "calendar_week"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить", "calendar_add"),
		),
	)
}

func NotesMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Все", "notes_list"),
			tgbotapi.NewInlineKeyboardButtonData("Добавить", "notes_add"),
		),
	)
}

func NotesMenuExtended() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Все", "notes_list"),
			tgbotapi.NewInlineKeyboardButtonData("Правила", "notes_rules"),
			tgbotapi.NewInlineKeyboardButtonData("+", "notes_add"),
		),
	)
}

func WebAppButton(url string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Открыть приложение", url),
		),
	)
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-3]) + "..."
}
