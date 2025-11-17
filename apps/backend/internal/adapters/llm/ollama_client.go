package llmadapter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Config struct {
	BaseURL     string
	Model       string
	Temperature float32
	TopP        float32
	MaxTokens   int
	Timeout     time.Duration
}

type OllamaClient struct {
	http        HTTPClient
	baseURL     string
	model       string
	temperature float32
	topP        float32
	maxTokens   int
}

var ErrEmptyPrompt = errors.New("llm: empty prompt")

func NewOllamaClient(httpClient HTTPClient, cfg Config) *OllamaClient {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "http://localhost:11434"
	}
	if cfg.Model == "" {
		cfg.Model = "llama3.1:8b"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: cfg.Timeout}
	}

	return &OllamaClient{
		http:        httpClient,
		baseURL:     cfg.BaseURL,
		model:       cfg.Model,
		temperature: cfg.Temperature,
		topP:        cfg.TopP,
		maxTokens:   cfg.MaxTokens,
	}
}

type ollamaChatRequest struct {
	Model    string              `json:"model"`
	Messages []ollamaChatMessage `json:"messages"`
	Stream   bool                `json:"stream"`
	Options  ollamaOptions       `json:"options,omitempty"`
}

type ollamaChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaOptions struct {
	Temperature float32 `json:"temperature,omitempty"`
	TopP        float32 `json:"top_p,omitempty"`
	NumPredict  int     `json:"num_predict,omitempty"`
}

type ollamaChatResponse struct {
	Message ollamaChatMessage `json:"message"`
}

func (c *OllamaClient) Generate(ctx context.Context, prompt []byte) (string, error) {
	if len(prompt) == 0 {
		return "", ErrEmptyPrompt
	}

	reqBody := ollamaChatRequest{
		Model: c.model,
		Messages: []ollamaChatMessage{
			{
				Role:    "user",
				Content: string(prompt),
			},
		},
		Stream: false,
		Options: ollamaOptions{
			Temperature: c.temperature,
			TopP:        c.topP,
			NumPredict:  c.maxTokens,
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("llm: marshal request: %w", err)
	}

	url := c.baseURL + "/api/chat"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("llm: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("llm: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4*1024))
		return "", fmt.Errorf("llm: bad status %d: %s", resp.StatusCode, string(b))
	}

	var r ollamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", fmt.Errorf("llm: decode response: %w", err)
	}

	if r.Message.Content == "" {
		return "", errors.New("llm: empty content in response")
	}

	return r.Message.Content, nil
}
