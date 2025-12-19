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

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) SetToken(token string) {
	c.token = token
}

type Topic struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    int        `json:"priority"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	DiscussedAt *time.Time `json:"discussed_at"`
}

type Event struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartAt     time.Time  `json:"start_at"`
	EndAt       *time.Time `json:"end_at"`
	AllDay      bool       `json:"all_day"`
	Shared      bool       `json:"shared"`
	Color       string     `json:"color"`
}

type Note struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	NoteType string `json:"note_type"`
	IsPinned bool   `json:"is_pinned"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
	} `json:"user"`
}

func (c *Client) Login(username, password string) (*LoginResponse, error) {
	body := map[string]string{"username": username, "password": password}
	var resp LoginResponse
	if err := c.post("/api/v1/auth/login", body, &resp); err != nil {
		return nil, err
	}
	c.token = resp.Token
	return &resp, nil
}

func (c *Client) GetTopics(status string, limit int) ([]Topic, error) {
	var resp struct {
		Topics []Topic `json:"topics"`
		Total  int     `json:"total"`
	}
	url := fmt.Sprintf("/api/v1/topics?status=%s&limit=%d", status, limit)
	if err := c.get(url, &resp); err != nil {
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
	if err := c.post("/api/v1/topics", body, &resp); err != nil {
		return nil, err
	}
	return &resp.Topic, nil
}

func (c *Client) MarkTopicDiscussed(id int) (*Topic, error) {
	var resp struct {
		Topic Topic `json:"topic"`
	}
	if err := c.post(fmt.Sprintf("/api/v1/topics/%d/discuss", id), nil, &resp); err != nil {
		return nil, err
	}
	return &resp.Topic, nil
}

func (c *Client) GetEvents(from, to string) ([]Event, error) {
	var resp struct {
		Events []Event `json:"events"`
	}
	url := fmt.Sprintf("/api/v1/events?from=%s&to=%s", from, to)
	if err := c.get(url, &resp); err != nil {
		return nil, err
	}
	return resp.Events, nil
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
	if err := c.post("/api/v1/notes", body, &resp); err != nil {
		return nil, err
	}
	return &resp.Note, nil
}

func (c *Client) GetNotes(noteType string, limit int) ([]Note, error) {
	var resp struct {
		Notes []Note `json:"notes"`
		Total int    `json:"total"`
	}
	url := fmt.Sprintf("/api/v1/notes?limit=%d", limit)
	if noteType != "" {
		url += "&note_type=" + noteType
	}
	if err := c.get(url, &resp); err != nil {
		return nil, err
	}
	return resp.Notes, nil
}

func (c *Client) get(path string, result interface{}) error {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return err
	}
	return c.do(req, result)
}

func (c *Client) post(path string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest("POST", c.baseURL+path, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.do(req, result)
}

func (c *Client) do(req *http.Request, result interface{}) error {
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

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}
	return nil
}
