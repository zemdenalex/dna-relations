package keyboards

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MainMenu() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Topics"),
			tgbotapi.NewKeyboardButton("Calendar"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Notes"),
			tgbotapi.NewKeyboardButton("Open App"),
		),
	)
}

func MainMenuWithWebApp(webappURL string) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				tgbotapi.NewKeyboardButton("Topics"),
				tgbotapi.NewKeyboardButton("Calendar"),
			},
			{
				tgbotapi.NewKeyboardButton("Notes"),
				tgbotapi.NewKeyboardButton("Open App"),
			},
		},
		ResizeKeyboard: true,
	}
}

func TopicsMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("View All", "topics_list"),
			tgbotapi.NewInlineKeyboardButtonData("Add New", "topics_add"),
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
			tgbotapi.NewInlineKeyboardButtonData("Pending", "topics_list"),
			tgbotapi.NewInlineKeyboardButtonData("Discussed", "topics_discussed"),
			tgbotapi.NewInlineKeyboardButtonData("Add", "topics_add"),
		},
	}
	for _, t := range topics {
		if t.Status == "pending" {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("Mark discussed: %s", truncate(t.Title, 20)),
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
			tgbotapi.NewInlineKeyboardButtonData("Today", "calendar_today"),
			tgbotapi.NewInlineKeyboardButtonData("This Week", "calendar_week"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Add Event", "calendar_add"),
		),
	)
}

func NotesMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("View All", "notes_list"),
			tgbotapi.NewInlineKeyboardButtonData("Add New", "notes_add"),
		),
	)
}

func NotesMenuExtended() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("All", "notes_list"),
			tgbotapi.NewInlineKeyboardButtonData("Rules", "notes_rules"),
			tgbotapi.NewInlineKeyboardButtonData("Add", "notes_add"),
		),
	)
}

func WebAppButton(url string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Open DNA App", url),
		),
	)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
