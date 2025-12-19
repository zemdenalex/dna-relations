package handlers

import (
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/zemdenalex/dna-relations/bot/internal/api"
	"github.com/zemdenalex/dna-relations/bot/internal/keyboards"
)

var priorityLabels = map[int]string{
	0: "Buffer",
	1: "Emergency",
	2: "ASAP",
	3: "Must discuss",
	4: "Soon",
	5: "When possible",
}

type Handler struct {
	bot       *tgbotapi.BotAPI
	apiClient *api.Client
	webappURL string
	userState map[int64]*UserState
}

type UserState struct {
	Action   string
	Data     map[string]interface{}
	LoggedIn bool
	UserID   int
	Username string
}

func New(bot *tgbotapi.BotAPI, apiBase, webappURL string) *Handler {
	return &Handler{
		bot:       bot,
		apiClient: api.NewClient(apiBase),
		webappURL: webappURL,
		userState: make(map[int64]*UserState),
	}
}

func (h *Handler) getState(chatID int64) *UserState {
	if s, ok := h.userState[chatID]; ok {
		return s
	}
	h.userState[chatID] = &UserState{Data: make(map[string]interface{})}
	return h.userState[chatID]
}

func (h *Handler) HandleMessage(msg *tgbotapi.Message) {
	if msg.IsCommand() {
		h.handleCommand(msg)
		return
	}

	h.handleText(msg)
}

func (h *Handler) handleCommand(msg *tgbotapi.Message) {
	switch msg.Command() {
	case "start":
		h.cmdStart(msg)
	case "topics":
		h.cmdTopics(msg)
	case "calendar":
		h.cmdCalendar(msg)
	case "notes":
		h.cmdNotes(msg)
	case "app":
		h.cmdApp(msg)
	case "help":
		h.cmdHelp(msg)
	case "login":
		h.cmdLogin(msg)
	default:
		h.sendText(msg.Chat.ID, "Unknown command. Use /help to see available commands.")
	}
}

func (h *Handler) handleText(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)

	switch state.Action {
	case "awaiting_login":
		h.handleLoginInput(msg)
		return
	case "awaiting_topic_title":
		h.handleTopicTitleInput(msg)
		return
	case "awaiting_topic_priority":
		h.handleTopicPriorityInput(msg)
		return
	case "awaiting_note_content":
		h.handleNoteContentInput(msg)
		return
	}

	switch msg.Text {
	case "Topics":
		h.cmdTopics(msg)
	case "Calendar":
		h.cmdCalendar(msg)
	case "Notes":
		h.cmdNotes(msg)
	case "Open App":
		h.cmdApp(msg)
	default:
		h.sendText(msg.Chat.ID, "Use the menu below or /help for commands.")
	}
}

func (h *Handler) HandleCallback(cb *tgbotapi.CallbackQuery) {
	callback := tgbotapi.NewCallback(cb.ID, "")
	h.bot.Request(callback)

	parts := strings.Split(cb.Data, ":")
	action := parts[0]

	switch action {
	case "topics_list":
		h.showTopicsList(cb.Message.Chat.ID, "pending")
	case "topics_discussed":
		h.showTopicsList(cb.Message.Chat.ID, "discussed")
	case "topics_add":
		h.startAddTopic(cb.Message.Chat.ID)
	case "topic_discuss":
		if len(parts) > 1 {
			h.markTopicDiscussed(cb.Message.Chat.ID, parts[1])
		}
	case "calendar_today":
		h.showTodayEvents(cb.Message.Chat.ID)
	case "calendar_week":
		h.showWeekEvents(cb.Message.Chat.ID)
	case "notes_list":
		h.showNotesList(cb.Message.Chat.ID, "")
	case "notes_rules":
		h.showNotesList(cb.Message.Chat.ID, "rule")
	case "notes_add":
		h.startAddNote(cb.Message.Chat.ID)
	default:
		log.Printf("Unknown callback: %s", cb.Data)
	}
}

func (h *Handler) cmdStart(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)

	if !state.LoggedIn {
		h.sendText(msg.Chat.ID, "Welcome to DNA Relations!\n\nPlease login first. Use /login command.")
		return
	}

	text := fmt.Sprintf("Hey %s! Welcome to DNA Relations.\n\nUse the menu below to navigate.", state.Username)
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ReplyMarkup = keyboards.MainMenu()
	h.bot.Send(reply)
}

func (h *Handler) cmdLogin(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	state.Action = "awaiting_login"
	h.sendText(msg.Chat.ID, "Enter your credentials in format:\nusername password\n\nExample: denis mypassword")
}

func (h *Handler) handleLoginInput(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	parts := strings.Fields(msg.Text)

	if len(parts) != 2 {
		h.sendText(msg.Chat.ID, "Invalid format. Please enter: username password")
		return
	}

	username, password := parts[0], parts[1]
	resp, err := h.apiClient.Login(username, password)
	if err != nil {
		h.sendText(msg.Chat.ID, "Login failed. Check your credentials.")
		log.Printf("Login error: %v", err)
		return
	}

	state.LoggedIn = true
	state.UserID = resp.User.ID
	state.Username = resp.User.Username
	state.Action = ""

	text := fmt.Sprintf("Logged in as %s!\n\nUse the menu below to navigate.", state.Username)
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ReplyMarkup = keyboards.MainMenu()
	h.bot.Send(reply)
}

func (h *Handler) cmdTopics(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	if !state.LoggedIn {
		h.sendText(msg.Chat.ID, "Please /login first")
		return
	}

	h.showTopicsList(msg.Chat.ID, "pending")
}

func (h *Handler) showTopicsList(chatID int64, status string) {
	topics, err := h.apiClient.GetTopics(status, 10)
	if err != nil {
		h.sendText(chatID, "Failed to fetch topics: "+err.Error())
		return
	}

	var text string
	if len(topics) == 0 {
		text = fmt.Sprintf("No %s topics.\n\nUse buttons below:", status)
	} else {
		text = fmt.Sprintf("%s topics:\n\n", strings.Title(status))
		for i, t := range topics {
			text += fmt.Sprintf("%d. [P%d] %s\n", i+1, t.Priority, t.Title)
		}
	}

	topicInfos := make([]keyboards.TopicInfo, len(topics))
	for i, t := range topics {
		topicInfos[i] = keyboards.TopicInfo{
			ID:       t.ID,
			Title:    t.Title,
			Priority: t.Priority,
			Status:   t.Status,
		}
	}

	reply := tgbotapi.NewMessage(chatID, text)
	reply.ReplyMarkup = keyboards.TopicsMenuExtended(topicInfos)
	h.bot.Send(reply)
}

func (h *Handler) startAddTopic(chatID int64) {
	state := h.getState(chatID)
	state.Action = "awaiting_topic_title"
	h.sendText(chatID, "Enter topic title:")
}

func (h *Handler) handleTopicTitleInput(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	state.Data["topic_title"] = msg.Text
	state.Action = "awaiting_topic_priority"
	h.sendText(msg.Chat.ID, "Enter priority (0-5):\n0 - Buffer\n1 - Emergency\n2 - ASAP\n3 - Must discuss\n4 - Soon\n5 - When possible")
}

func (h *Handler) handleTopicPriorityInput(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)

	var priority int
	fmt.Sscanf(msg.Text, "%d", &priority)
	if priority < 0 || priority > 5 {
		h.sendText(msg.Chat.ID, "Invalid priority. Enter 0-5:")
		return
	}

	title := state.Data["topic_title"].(string)
	topic, err := h.apiClient.CreateTopic(title, "", priority)
	if err != nil {
		h.sendText(msg.Chat.ID, "Failed to create topic: "+err.Error())
		state.Action = ""
		return
	}

	state.Action = ""
	state.Data = make(map[string]interface{})
	h.sendText(msg.Chat.ID, fmt.Sprintf("✅ Topic created: %s (P%d)", topic.Title, topic.Priority))
}

func (h *Handler) markTopicDiscussed(chatID int64, idStr string) {
	var id int
	fmt.Sscanf(idStr, "%d", &id)

	topic, err := h.apiClient.MarkTopicDiscussed(id)
	if err != nil {
		h.sendText(chatID, "Failed to mark topic: "+err.Error())
		return
	}

	h.sendText(chatID, fmt.Sprintf("✅ Marked as discussed: %s", topic.Title))
	h.showTopicsList(chatID, "pending")
}

func (h *Handler) cmdCalendar(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	if !state.LoggedIn {
		h.sendText(msg.Chat.ID, "Please /login first")
		return
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, "📅 Calendar:")
	reply.ReplyMarkup = keyboards.CalendarMenu()
	h.bot.Send(reply)
}

func (h *Handler) showTodayEvents(chatID int64) {
	today := time.Now().Format("2006-01-02")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	events, err := h.apiClient.GetEvents(today, tomorrow)
	if err != nil {
		h.sendText(chatID, "Failed to fetch events: "+err.Error())
		return
	}

	if len(events) == 0 {
		h.sendText(chatID, "No events today. 🎉")
		return
	}

	text := "📅 Today's events:\n\n"
	for _, e := range events {
		timeStr := e.StartAt.Format("15:04")
		text += fmt.Sprintf("• %s - %s\n", timeStr, e.Title)
	}
	h.sendText(chatID, text)
}

func (h *Handler) showWeekEvents(chatID int64) {
	from := time.Now().Format("2006-01-02")
	to := time.Now().AddDate(0, 0, 7).Format("2006-01-02")

	events, err := h.apiClient.GetEvents(from, to)
	if err != nil {
		h.sendText(chatID, "Failed to fetch events: "+err.Error())
		return
	}

	if len(events) == 0 {
		h.sendText(chatID, "No events this week.")
		return
	}

	text := "📅 This week's events:\n\n"
	for _, e := range events {
		dateStr := e.StartAt.Format("Mon 02 Jan 15:04")
		text += fmt.Sprintf("• %s - %s\n", dateStr, e.Title)
	}
	h.sendText(chatID, text)
}

func (h *Handler) cmdNotes(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	if !state.LoggedIn {
		h.sendText(msg.Chat.ID, "Please /login first")
		return
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, "📝 Notes:")
	reply.ReplyMarkup = keyboards.NotesMenuExtended()
	h.bot.Send(reply)
}

func (h *Handler) showNotesList(chatID int64, noteType string) {
	notes, err := h.apiClient.GetNotes(noteType, 10)
	if err != nil {
		h.sendText(chatID, "Failed to fetch notes: "+err.Error())
		return
	}

	if len(notes) == 0 {
		h.sendText(chatID, "No notes found.")
		return
	}

	text := "📝 Notes:\n\n"
	for i, n := range notes {
		title := n.Title
		if title == "" {
			title = truncate(n.Content, 30)
		}
		pinned := ""
		if n.IsPinned {
			pinned = "📌 "
		}
		text += fmt.Sprintf("%d. %s[%s] %s\n", i+1, pinned, n.NoteType, title)
	}
	h.sendText(chatID, text)
}

func (h *Handler) startAddNote(chatID int64) {
	state := h.getState(chatID)
	state.Action = "awaiting_note_content"
	h.sendText(chatID, "Enter note content:")
}

func (h *Handler) handleNoteContentInput(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)

	note, err := h.apiClient.CreateNote("", msg.Text, "general")
	if err != nil {
		h.sendText(msg.Chat.ID, "Failed to create note: "+err.Error())
		state.Action = ""
		return
	}

	state.Action = ""
	h.sendText(msg.Chat.ID, fmt.Sprintf("✅ Note created (ID: %d)", note.ID))
}

func (h *Handler) cmdApp(msg *tgbotapi.Message) {
	if h.webappURL == "" {
		h.sendText(msg.Chat.ID, "Web app URL is not configured.")
		return
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, "Open the mini app:")
	reply.ReplyMarkup = keyboards.WebAppButton(h.webappURL)
	h.bot.Send(reply)
}

func (h *Handler) cmdHelp(msg *tgbotapi.Message) {
	text := `Available commands:

/login - Login to your account
/start - Show main menu
/topics - View discussion topics
/calendar - View calendar
/notes - View notes
/app - Open mini app
/help - Show this help`

	h.sendText(msg.Chat.ID, text)
}

func (h *Handler) sendText(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	h.bot.Send(msg)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
