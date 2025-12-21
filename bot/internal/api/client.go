package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type Topic struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    int        `json:"priority"`
	Status      string     `json:"status"`
	CreatedBy   int        `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	DiscussedAt *time.Time `json:"discussed_at"`
}

type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	AllDay      bool      `json:"all_day"`
	Color       string    `json:"color"`
	Shared      bool      `json:"shared"`
	OwnerID     int       `json:"owner_id"`
}

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	NoteType  string    `json:"note_type"`
	IsPinned  bool      `json:"is_pinned"`
	Shared    bool      `json:"shared"`
	CreatedBy int       `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) Login(username, password string) (*LoginResponse, error) {
	body := map[string]string{
		"username": username,
		"password": password,
	}

	var resp struct {
		Token string `json:"token"`
		User  User   `json:"user"`
	}

	err := c.post("/api/v1/auth/login", body, &resp)
	if err != nil {
		return nil, err
	}

	c.token = resp.Token
	return &LoginResponse{Token: resp.Token, User: resp.User}, nil
}

func (c *Client) GetTopics(status string, limit int) ([]Topic, error) {
	url := fmt.Sprintf("/api/v1/topics?status=%s&limit=%d", status, limit)

	var resp struct {
		Topics []Topic `json:"topics"`
	}

	err := c.get(url, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Topics, nil
}

func (c *Client) CreateTopic(title, description string, priority int) (*Topic, error) {
	body := map[string]interface{}{
		"title":       title,
		"description": description,
		"priority":    priority,
	}

	var resp struct {
		Topic Topic `json:"topic"`
	}

	err := c.post("/api/v1/topics", body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Topic, nil
}

func (c *Client) MarkTopicDiscussed(id int) (*Topic, error) {
	var resp struct {
		Topic Topic `json:"topic"`
	}

	err := c.post(fmt.Sprintf("/api/v1/topics/%d/discuss", id), nil, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Topic, nil
}

func (c *Client) GetEvents(from, to string) ([]Event, error) {
	url := fmt.Sprintf("/api/v1/events?from=%s&to=%s", from, to)

	var resp struct {
		Events []Event `json:"events"`
	}

	err := c.get(url, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Events, nil
}

func (c *Client) CreateEvent(title, description string, startAt, endAt time.Time) (*Event, error) {
	body := map[string]interface{}{
		"title":       title,
		"description": description,
		"start_at":    startAt.Format(time.RFC3339),
		"end_at":      endAt.Format(time.RFC3339),
		"shared":      true,
	}

	var resp struct {
		Event Event `json:"event"`
	}

	err := c.post("/api/v1/events", body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Event, nil
}

func (c *Client) GetNotes(noteType string, limit int) ([]Note, error) {
	url := fmt.Sprintf("/api/v1/notes?limit=%d", limit)
	if noteType != "" {
		url += "&note_type=" + noteType
	}

	var resp struct {
		Notes []Note `json:"notes"`
	}

	err := c.get(url, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Notes, nil
}

func (c *Client) CreateNote(title, content, noteType string) (*Note, error) {
	body := map[string]interface{}{
		"title":     title,
		"content":   content,
		"note_type": noteType,
		"shared":    true,
	}

	var resp struct {
		Note Note `json:"note"`
	}

	err := c.post("/api/v1/notes", body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.Note, nil
}

func (c *Client) get(path string, result interface{}) error {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

func (c *Client) post(path string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest("POST", c.baseURL+path, reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	if result != nil && resp.StatusCode != 204 {
		return json.NewDecoder(resp.Body).Decode(result)
	}

	return nil
}
