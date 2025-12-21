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
	0: "Буфер",
	1: "Срочно",
	2: "ASAP",
	3: "Надо обсудить",
	4: "Скоро",
	5: "Когда-нибудь",
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
		h.sendText(msg.Chat.ID, "Неизвестная команда. Используй /help")
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
	case "Темы", "Topics":
		h.cmdTopics(msg)
	case "Календарь", "Calendar":
		h.cmdCalendar(msg)
	case "Заметки", "Notes":
		h.cmdNotes(msg)
	case "Приложение", "Open App":
		h.cmdApp(msg)
	default:
		h.sendText(msg.Chat.ID, "Используй меню ниже или /help")
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
		text := `Привет! Это DNA Relations.

Бот для управления:
- Темами для обсуждения
- Общим календарём
- Заметками и правилами

Войди через /login`
		h.sendText(msg.Chat.ID, text)
		return
	}

	text := fmt.Sprintf("Привет, %s!", state.Username)
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ReplyMarkup = keyboards.MainMenuWithWebApp(h.webappURL)
	h.bot.Send(reply)
}

func (h *Handler) cmdLogin(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)

	if state.LoggedIn {
		h.sendText(msg.Chat.ID, fmt.Sprintf("Уже залогинен как %s. Используй /logout", state.Username))
		return
	}

	state.Action = "awaiting_login"
	h.sendText(msg.Chat.ID, "Введи логин и пароль через пробел:\n\nПример: denis mypassword")
}

func (h *Handler) cmdLogout(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)

	if !state.LoggedIn {
		h.sendText(msg.Chat.ID, "Ты не залогинен")
		return
	}

	state.LoggedIn = false
	state.UserID = 0
	state.Username = ""
	state.Action = ""
	state.Data = make(map[string]interface{})

	h.sendText(msg.Chat.ID, "Вышел. Используй /login чтобы войти снова")
}

func (h *Handler) handleLoginInput(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	parts := strings.Fields(msg.Text)

	if len(parts) != 2 {
		h.sendText(msg.Chat.ID, "Неверный формат. Введи: логин пароль")
		return
	}

	username, password := parts[0], parts[1]
	resp, err := h.apiClient.Login(username, password)
	if err != nil {
		h.sendText(msg.Chat.ID, fmt.Sprintf("Ошибка входа: %v", err))
		log.Printf("Login error for %s: %v", username, err)
		state.Action = ""
		return
	}

	state.LoggedIn = true
	state.UserID = resp.User.ID
	state.Username = resp.User.Username
	state.Action = ""

	text := fmt.Sprintf("Привет, %s!", state.Username)
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ReplyMarkup = keyboards.MainMenuWithWebApp(h.webappURL)
	h.bot.Send(reply)
}

func (h *Handler) cmdTopics(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	if !state.LoggedIn {
		h.sendText(msg.Chat.ID, "Сначала /login")
		return
	}

	h.showTopicsList(msg.Chat.ID, "pending")
}

func (h *Handler) cmdNewTopic(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	if !state.LoggedIn {
		h.sendText(msg.Chat.ID, "Сначала /login")
		return
	}
	h.startAddTopic(msg.Chat.ID)
}

func (h *Handler) showTopicsList(chatID int64, status string) {
	topics, err := h.apiClient.GetTopics(status, 10)
	if err != nil {
		h.sendText(chatID, "Ошибка: "+err.Error())
		return
	}

	var text string
	statusRu := "ожидающие"
	if status == "discussed" {
		statusRu = "обсуждённые"
	}

	if len(topics) == 0 {
		text = fmt.Sprintf("Нет %s тем", statusRu)
	} else {
		text = fmt.Sprintf("<b>Темы (%s):</b>\n\n", statusRu)
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
	h.sendText(chatID, "Введи название темы (и описание на новой строке):\n\nПример:\nОбсудить отпуск\nНадо выбрать даты и место")
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

	text := "Выбери приоритет (0-5):\n"
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
		h.sendText(msg.Chat.ID, "Неверный приоритет. Введи 0-5:")
		return
	}

	title := state.Data["topic_title"].(string)
	description := ""
	if desc, ok := state.Data["topic_description"].(string); ok {
		description = desc
	}

	topic, err := h.apiClient.CreateTopic(title, description, priority)
	if err != nil {
		h.sendText(msg.Chat.ID, "Ошибка: "+err.Error())
		state.Action = ""
		return
	}

	state.Action = ""
	state.Data = make(map[string]interface{})
	h.sendText(msg.Chat.ID, fmt.Sprintf("Тема создана: <b>%s</b>\nПриоритет: %s", topic.Title, priorityLabels[topic.Priority]))
}

func (h *Handler) markTopicDiscussed(chatID int64, idStr string) {
	var id int
	fmt.Sscanf(idStr, "%d", &id)

	topic, err := h.apiClient.MarkTopicDiscussed(id)
	if err != nil {
		h.sendText(chatID, "Ошибка: "+err.Error())
		return
	}

	h.sendText(chatID, fmt.Sprintf("Обсуждено: <b>%s</b>", topic.Title))
	h.showTopicsList(chatID, "pending")
}

func (h *Handler) cmdCalendar(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	if !state.LoggedIn {
		h.sendText(msg.Chat.ID, "Сначала /login")
		return
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, "Календарь:")
	reply.ReplyMarkup = keyboards.CalendarMenu()
	h.bot.Send(reply)
}

func (h *Handler) cmdToday(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	if !state.LoggedIn {
		h.sendText(msg.Chat.ID, "Сначала /login")
		return
	}
	h.showTodayEvents(msg.Chat.ID)
}

func (h *Handler) showTodayEvents(chatID int64) {
	today := time.Now().Format("2006-01-02")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	events, err := h.apiClient.GetEvents(today, tomorrow)
	if err != nil {
		h.sendText(chatID, "Ошибка: "+err.Error())
		return
	}

	if len(events) == 0 {
		h.sendText(chatID, fmt.Sprintf("<b>Сегодня (%s)</b>\n\nНет событий", time.Now().Format("02.01")))
		return
	}

	text := fmt.Sprintf("<b>Сегодня (%s)</b>\n\n", time.Now().Format("02.01"))
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
		h.sendText(chatID, "Ошибка: "+err.Error())
		return
	}

	if len(events) == 0 {
		h.sendText(chatID, "<b>Эта неделя</b>\n\nНет событий")
		return
	}

	text := "<b>Эта неделя</b>\n\n"

	eventsByDay := make(map[string][]api.Event)
	for _, e := range events {
		day := e.StartAt.Format("2006-01-02")
		eventsByDay[day] = append(eventsByDay[day], e)
	}

	weekdays := []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}
	for i := 0; i < 7; i++ {
		day := time.Now().AddDate(0, 0, i)
		dayStr := day.Format("2006-01-02")
		if dayEvents, ok := eventsByDay[dayStr]; ok {
			wd := int(day.Weekday())
			if wd == 0 {
				wd = 7
			}
			text += fmt.Sprintf("<b>%s, %s</b>\n", weekdays[wd-1], day.Format("02.01"))
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
	h.sendText(chatID, "Введи название события:")
}

func (h *Handler) handleEventTitleInput(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	state.Data["event_title"] = msg.Text
	state.Action = "awaiting_event_datetime"
	h.sendText(msg.Chat.ID, "Введи дату и время:\nФормат: ДД.ММ.ГГГГ ЧЧ:ММ или ГГГГ-ММ-ДД ЧЧ:ММ\n\nПример: 25.12.2025 19:00")
}

func (h *Handler) handleEventDatetimeInput(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)

	input := strings.TrimSpace(msg.Text)
	var startTime time.Time
	var err error

	startTime, err = time.Parse("02.01.2006 15:04", input)
	if err != nil {
		startTime, err = time.Parse("2006-01-02 15:04", input)
	}
	if err != nil {
		h.sendText(msg.Chat.ID, "Неверный формат. Используй: ДД.ММ.ГГГГ ЧЧ:ММ\nПример: 25.12.2025 19:00")
		return
	}

	title := state.Data["event_title"].(string)
	endTime := startTime.Add(1 * time.Hour)

	event, err := h.apiClient.CreateEvent(title, "", startTime, endTime)
	if err != nil {
		h.sendText(msg.Chat.ID, "Ошибка: "+err.Error())
		state.Action = ""
		return
	}

	state.Action = ""
	state.Data = make(map[string]interface{})
	h.sendText(msg.Chat.ID, fmt.Sprintf("Событие создано: <b>%s</b>\n%s", event.Title, startTime.Format("02.01 в 15:04")))
}

func (h *Handler) cmdNotes(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	if !state.LoggedIn {
		h.sendText(msg.Chat.ID, "Сначала /login")
		return
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, "Заметки:")
	reply.ReplyMarkup = keyboards.NotesMenuExtended()
	h.bot.Send(reply)
}

func (h *Handler) cmdNewNote(msg *tgbotapi.Message) {
	state := h.getState(msg.Chat.ID)
	if !state.LoggedIn {
		h.sendText(msg.Chat.ID, "Сначала /login")
		return
	}
	h.startAddNote(msg.Chat.ID)
}

func (h *Handler) showNotesList(chatID int64, noteType string) {
	notes, err := h.apiClient.GetNotes(noteType, 10)
	if err != nil {
		h.sendText(chatID, "Ошибка: "+err.Error())
		return
	}

	if len(notes) == 0 {
		h.sendText(chatID, "Нет заметок")
		return
	}

	text := "<b>Заметки</b>\n\n"
	for i, n := range notes {
		title := n.Title
		if title == "" {
			title = truncate(n.Content, 30)
		}
		pinned := ""
		if n.IsPinned {
			pinned = " [закреплено]"
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
	h.sendText(chatID, "Введи заметку (первая строка = название):\n\nПример:\nСписок покупок\n- Молоко\n- Хлеб")
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
		h.sendText(msg.Chat.ID, "Ошибка: "+err.Error())
		state.Action = ""
		return
	}

	state.Action = ""
	h.sendText(msg.Chat.ID, fmt.Sprintf("Заметка создана: <b>%s</b>", note.Title))
}

func (h *Handler) cmdApp(msg *tgbotapi.Message) {
	if h.webappURL == "" {
		h.sendText(msg.Chat.ID, "URL приложения не настроен")
		return
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, "Открыть приложение:")
	reply.ReplyMarkup = keyboards.WebAppButton(h.webappURL)
	h.bot.Send(reply)
}

func (h *Handler) cmdHelp(msg *tgbotapi.Message) {
	text := `<b>Команды:</b>

<b>Авторизация:</b>
/login - Войти
/logout - Выйти

<b>Темы:</b>
/topics - Список тем
/newtopic - Новая тема

<b>Календарь:</b>
/calendar - Меню календаря
/today - События сегодня

<b>Заметки:</b>
/notes - Список заметок
/newnote - Новая заметка

<b>Другое:</b>
/app - Открыть приложение
/help - Эта справка`

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
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-3]) + "..."
}
