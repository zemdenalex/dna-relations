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
	case "logout":
		h.cmdLogout(msg)
	case "today":
		h.cmdToday(msg)
	case "newtopic":
		h.cmdNewTopic(msg)
	case "newnote":
		h.cmdNewNote(msg)
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
	case "awaiting_event_title":
		h.handleEventTitleInput(msg)
		return
	case "awaiting_event_datetime":
		h.handleEventDatetimeInput(msg)
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
	case "calendar_add":
		h.startAddEvent(cb.Message.Chat.ID)
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
		text := `Welcome to DNA Relations!

This bot helps you and your partner manage:
- Topics to discuss
- Shared calendar
- Notes and ideas
- Daily rituals

Please login first with /login`
		h.sendText(msg.Chat.ID, text)
		return
	}

	text := fmt.Sprintf("Hey %s! Welcome to DNA Relations.\n\nUse the menu below to navigate.", state.Username)
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ReplyMarkup = keyboards.MainMenu()
	h.bot.Send(reply)
}

func (h *Handler) cmdLogin(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	
	if state.LoggedIn {
		h.sendText(msg.Chat.ID, fmt.Sprintf("Already logged in as %s. Use /logout first.", state.Username))
		return
	}
	
	state.Action = "awaiting_login"
	h.sendText(msg.Chat.ID, "Enter your credentials in format:\nusername password\n\nExample: denis mypassword")
}

func (h *Handler) cmdLogout(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	
	if !state.LoggedIn {
		h.sendText(msg.Chat.ID, "Not logged in.")
		return
	}
	
	state.LoggedIn = false
	state.UserID = 0
	state.Username = ""
	state.Action = ""
	state.Data = make(map[string]interface{})
	
	h.sendText(msg.Chat.ID, "Logged out successfully. Use /login to login again.")
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
		h.sendText(msg.Chat.ID, fmt.Sprintf("Login failed: %v", err))
		log.Printf("Login error for %s: %v", username, err)
		state.Action = ""
		return
	}

	state.LoggedIn = true
	state.UserID = resp.User.ID
	state.Username = resp.User.Username
	state.Action = ""

	text := fmt.Sprintf("Logged in as %s!\n\nUse the menu below to navigate.", state.Username)
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ReplyMarkup = keyboards.MainMenuWithWebApp(h.webappURL)
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

func (h *Handler) cmdNewTopic(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	if !state.LoggedIn {
		h.sendText(msg.Chat.ID, "Please /login first")
		return
	}
	h.startAddTopic(msg.Chat.ID)
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
		text = fmt.Sprintf("<b>%s topics:</b>\n\n", strings.Title(status))
		for i, t := range topics {
			priorityLabel := priorityLabels[t.Priority]
			text += fmt.Sprintf("%d. <b>[%s]</b> %s\n", i+1, priorityLabel, t.Title)
			if t.Description != "" {
				text += fmt.Sprintf("   <i>%s</i>\n", truncate(t.Description, 50))
			}
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
	reply.ParseMode = "HTML"
	reply.ReplyMarkup = keyboards.TopicsMenuExtended(topicInfos)
	h.bot.Send(reply)
}

func (h *Handler) startAddTopic(chatID int64) {
	state := h.getState(chatID)
	state.Action = "awaiting_topic_title"
	h.sendText(chatID, "Enter topic title (and optionally description on new line):\n\nExample:\nDiscuss vacation plans\nWe need to decide dates and destination")
}

func (h *Handler) handleTopicTitleInput(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	
	lines := strings.SplitN(msg.Text, "\n", 2)
	title := strings.TrimSpace(lines[0])
	description := ""
	if len(lines) > 1 {
		description = strings.TrimSpace(lines[1])
	}
	
	state.Data["topic_title"] = title
	state.Data["topic_description"] = description
	state.Action = "awaiting_topic_priority"
	
	text := "Enter priority (0-5):\n"
	for i := 0; i <= 5; i++ {
		text += fmt.Sprintf("%d - %s\n", i, priorityLabels[i])
	}
	h.sendText(msg.Chat.ID, text)
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
	description := ""
	if desc, ok := state.Data["topic_description"].(string); ok {
		description = desc
	}
	
	topic, err := h.apiClient.CreateTopic(title, description, priority)
	if err != nil {
		h.sendText(msg.Chat.ID, "Failed to create topic: "+err.Error())
		state.Action = ""
		return
	}

	state.Action = ""
	state.Data = make(map[string]interface{})
	h.sendText(msg.Chat.ID, fmt.Sprintf("Topic created: <b>%s</b>\nPriority: %s", topic.Title, priorityLabels[topic.Priority]))
}

func (h *Handler) markTopicDiscussed(chatID int64, idStr string) {
	var id int
	fmt.Sscanf(idStr, "%d", &id)

	topic, err := h.apiClient.MarkTopicDiscussed(id)
	if err != nil {
		h.sendText(chatID, "Failed to mark topic: "+err.Error())
		return
	}

	h.sendText(chatID, fmt.Sprintf("Marked as discussed: <b>%s</b>", topic.Title))
	h.showTopicsList(chatID, "pending")
}

func (h *Handler) cmdCalendar(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	if !state.LoggedIn {
		h.sendText(msg.Chat.ID, "Please /login first")
		return
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, "Calendar:")
	reply.ReplyMarkup = keyboards.CalendarMenu()
	h.bot.Send(reply)
}

func (h *Handler) cmdToday(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	if !state.LoggedIn {
		h.sendText(msg.Chat.ID, "Please /login first")
		return
	}
	h.showTodayEvents(msg.Chat.ID)
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
		h.sendText(chatID, fmt.Sprintf("<b>Today (%s)</b>\n\nNo events scheduled.", time.Now().Format("Monday, January 2")))
		return
	}

	text := fmt.Sprintf("<b>Today (%s)</b>\n\n", time.Now().Format("Monday, January 2"))
	for _, e := range events {
		timeStr := e.StartAt.Format("15:04")
		text += fmt.Sprintf("• <b>%s</b> - %s\n", timeStr, e.Title)
		if e.Description != "" {
			text += fmt.Sprintf("  <i>%s</i>\n", truncate(e.Description, 50))
		}
	}
	
	reply := tgbotapi.NewMessage(chatID, text)
	reply.ParseMode = "HTML"
	h.bot.Send(reply)
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
		h.sendText(chatID, "<b>This Week</b>\n\nNo events scheduled.")
		return
	}

	text := "<b>This Week's Events</b>\n\n"
	
	eventsByDay := make(map[string][]api.Event)
	for _, e := range events {
		day := e.StartAt.Format("2006-01-02")
		eventsByDay[day] = append(eventsByDay[day], e)
	}
	
	for i := 0; i < 7; i++ {
		day := time.Now().AddDate(0, 0, i)
		dayStr := day.Format("2006-01-02")
		if dayEvents, ok := eventsByDay[dayStr]; ok {
			text += fmt.Sprintf("<b>%s</b>\n", day.Format("Mon, Jan 2"))
			for _, e := range dayEvents {
				text += fmt.Sprintf("  %s - %s\n", e.StartAt.Format("15:04"), e.Title)
			}
			text += "\n"
		}
	}
	
	reply := tgbotapi.NewMessage(chatID, text)
	reply.ParseMode = "HTML"
	h.bot.Send(reply)
}

func (h *Handler) startAddEvent(chatID int64) {
	state := h.getState(chatID)
	state.Action = "awaiting_event_title"
	h.sendText(chatID, "Enter event title:")
}

func (h *Handler) handleEventTitleInput(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	state.Data["event_title"] = msg.Text
	state.Action = "awaiting_event_datetime"
	h.sendText(msg.Chat.ID, "Enter date and time:\nFormat: YYYY-MM-DD HH:MM\n\nExample: 2025-12-25 19:00")
}

func (h *Handler) handleEventDatetimeInput(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	
	startTime, err := time.Parse("2006-01-02 15:04", strings.TrimSpace(msg.Text))
	if err != nil {
		h.sendText(msg.Chat.ID, "Invalid format. Use: YYYY-MM-DD HH:MM\nExample: 2025-12-25 19:00")
		return
	}
	
	title := state.Data["event_title"].(string)
	endTime := startTime.Add(1 * time.Hour)
	
	event, err := h.apiClient.CreateEvent(title, "", startTime, endTime)
	if err != nil {
		h.sendText(msg.Chat.ID, "Failed to create event: "+err.Error())
		state.Action = ""
		return
	}
	
	state.Action = ""
	state.Data = make(map[string]interface{})
	h.sendText(msg.Chat.ID, fmt.Sprintf("Event created: <b>%s</b>\n%s", event.Title, startTime.Format("Mon, Jan 2 at 15:04")))
}

func (h *Handler) cmdNotes(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	if !state.LoggedIn {
		h.sendText(msg.Chat.ID, "Please /login first")
		return
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, "Notes:")
	reply.ReplyMarkup = keyboards.NotesMenuExtended()
	h.bot.Send(reply)
}

func (h *Handler) cmdNewNote(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	if !state.LoggedIn {
		h.sendText(msg.Chat.ID, "Please /login first")
		return
	}
	h.startAddNote(msg.Chat.ID)
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

	text := "<b>Notes</b>\n\n"
	for i, n := range notes {
		title := n.Title
		if title == "" {
			title = truncate(n.Content, 30)
		}
		pinned := ""
		if n.IsPinned {
			pinned = " [pinned]"
		}
		text += fmt.Sprintf("%d. <b>%s</b>%s\n", i+1, title, pinned)
		if n.Content != "" && n.Title != "" {
			text += fmt.Sprintf("   <i>%s</i>\n", truncate(n.Content, 40))
		}
	}
	
	reply := tgbotapi.NewMessage(chatID, text)
	reply.ParseMode = "HTML"
	h.bot.Send(reply)
}

func (h *Handler) startAddNote(chatID int64) {
	state := h.getState(chatID)
	state.Action = "awaiting_note_content"
	h.sendText(chatID, "Enter note (first line = title, rest = content):\n\nExample:\nGrocery list\n- Milk\n- Bread\n- Eggs")
}

func (h *Handler) handleNoteContentInput(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)

	lines := strings.SplitN(msg.Text, "\n", 2)
	title := strings.TrimSpace(lines[0])
	content := ""
	if len(lines) > 1 {
		content = strings.TrimSpace(lines[1])
	}

	note, err := h.apiClient.CreateNote(title, content, "general")
	if err != nil {
		h.sendText(msg.Chat.ID, "Failed to create note: "+err.Error())
		state.Action = ""
		return
	}

	state.Action = ""
	h.sendText(msg.Chat.ID, fmt.Sprintf("Note created: <b>%s</b>", note.Title))
}

func (h *Handler) cmdApp(msg *tgbotapi.Message) {
	if h.webappURL == "" {
		h.sendText(msg.Chat.ID, "Web app URL is not configured.")
		return
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, "Open the DNA app:")
	reply.ReplyMarkup = keyboards.WebAppButton(h.webappURL)
	h.bot.Send(reply)
}

func (h *Handler) cmdHelp(msg *tgbotapi.Message) {
	text := `<b>Available commands:</b>

<b>Auth:</b>
/login - Login to your account
/logout - Logout

<b>Topics:</b>
/topics - View discussion topics
/newtopic - Create new topic

<b>Calendar:</b>
/calendar - View calendar menu
/today - Today's events

<b>Notes:</b>
/notes - View notes
/newnote - Create new note

<b>Other:</b>
/app - Open mini app
/help - Show this help`

	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = "HTML"
	h.bot.Send(reply)
}

func (h *Handler) sendText(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	h.bot.Send(msg)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}