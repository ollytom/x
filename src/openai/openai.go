package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant      = "assistant"
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

type Chat struct {
	Messages       []Message `json:"messages"`
	Model          string    `json:"model"`
	ResponseFormat *struct {
		Type string `json:"type"`
	} `json:"response_format,omitempty"`
}

type Model struct {
	ID               string
	Name             string
	Description      string
	Aliases          []string
	MaxContextLength int
	Deprecation      time.Time
}

type Client struct {
	*http.Client
	Token   string
	BaseURL string
}

type apiError struct {
	Message struct {
		Detail []struct {
			Msg string
		}
	}
	Type string
}

func (e apiError) Error() string {
	messages := make([]string, len(e.Message.Detail))
	for i := range e.Message.Detail {
		messages[i] = e.Message.Detail[i].Msg
	}
	return fmt.Sprintf("%s: %s", e.Type, strings.Join(messages, ", "))
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	if c.Client == nil {
		c.Client = http.DefaultClient
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	req.Header.Set("Accept", "application/json")
	if req.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.Do(req)
}

type completeResponse struct {
	Choices []struct {
		Message Message
	}
}

func (c *Client) Complete(chat *Chat) (*Message, error) {
	b, err := json.Marshal(chat)
	if err != nil {
		return nil, fmt.Errorf("encode messages: %w", err)
	}
	u := c.BaseURL + "/v1/chat/completions"
	req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 && resp.StatusCode <= 499 {
		var aerr apiError
		if err := json.NewDecoder(resp.Body).Decode(&aerr); err != nil {
			return nil, fmt.Errorf(resp.Status)
		}
		return nil, aerr
	} else if resp.StatusCode >= 500 {
		return nil, fmt.Errorf(resp.Status)
	}

	var cresp completeResponse
	if err := json.NewDecoder(resp.Body).Decode(&cresp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if len(cresp.Choices) == 0 {
		return nil, fmt.Errorf("no completions in response")
	}
	return &cresp.Choices[0].Message, nil
}

func (c *Client) Models() ([]Model, error) {
	u := c.BaseURL + "/v1/models"
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var aerr apiError
		if err := json.NewDecoder(resp.Body).Decode(&aerr); err != nil {
			return nil, fmt.Errorf(resp.Status)
		}
		return nil, aerr
	}
	v := struct {
		Data []Model
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return v.Data, nil
}
