package llm

import (
	"backend/internal/domain"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	BaseURL       string
	Model         string
	Temperature   float32
	TopP          float32
	MaxTokens     int
	Timeout       time.Duration
	EnableWebSearch bool // Включает автоматический веб-поиск
}

type OllamaClient struct {
	client    *http.Client
	config    Config
	webSearch *WebSearch // Модуль для веб-поиска
}

func NewOllamaClient(client *http.Client, config Config) *OllamaClient {
	oc := &OllamaClient{
		client: client,
		config: config,
	}

	// Инициализируем веб-поиск если он включен
	if config.EnableWebSearch {
		oc.webSearch = NewWebSearch(client)
	}

	return oc
}

// Generate реализует интерфейс domain.LLM
func (c *OllamaClient) Generate(ctx context.Context, prompt []byte) (string, error) {
	promptStr := string(prompt)

	// Если веб-поиск включен, проверяем нужен ли он
	if c.config.EnableWebSearch && c.webSearch != nil && c.needsWebSearch(promptStr) {
		webResults, err := c.performWebSearch(ctx, promptStr)
		if err == nil && len(webResults) > 0 {
			// Добавляем результаты веб-поиска в промпт
			promptStr = c.enrichPromptWithWebSearch(promptStr, webResults)
		}
	}

	url := fmt.Sprintf("%s/api/generate", c.config.BaseURL)

	requestBody := map[string]interface{}{
		"model":       c.config.Model,
		"prompt":      promptStr,
		"stream":      false,
		"temperature": c.config.Temperature,
		"top_p":       c.config.TopP,
		"num_predict": c.config.MaxTokens,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		Response string `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Response, nil
}

// needsWebSearch определяет, нужен ли веб-поиск для данного промпта
func (c *OllamaClient) needsWebSearch(prompt string) bool {
	lowerPrompt := strings.ToLower(prompt)
	
	// Ключевые слова, указывающие на необходимость актуальной информации
	searchKeywords := []string{
		"текущий", "актуальный", "сегодня", "сейчас", "последний",
		"новости", "события", "курс", "цена", "погода",
		"когда", "где находится", "адрес", "расписание",
		"что происходит", "что случилось", "когда будет",
		"проверь в интернете", "найди в интернете", "поищи",
	}

	for _, keyword := range searchKeywords {
		if strings.Contains(lowerPrompt, keyword) {
			return true
		}
	}

	// Проверяем, есть ли явная инструкция о веб-поиске
	if strings.Contains(lowerPrompt, "[web_search]") || strings.Contains(lowerPrompt, "[поиск]") {
		return true
	}

	return false
}

// performWebSearch извлекает запрос из промпта и выполняет веб-поиск
func (c *OllamaClient) performWebSearch(ctx context.Context, prompt string) ([]SearchResult, error) {
	// Извлекаем основной запрос пользователя (последняя строка после "USER:")
	parts := strings.Split(prompt, "USER:")
	query := prompt
	if len(parts) > 1 {
		query = strings.TrimSpace(parts[len(parts)-1])
	}

	// Удаляем метки веб-поиска
	query = strings.ReplaceAll(query, "[web_search]", "")
	query = strings.ReplaceAll(query, "[поиск]", "")
	query = strings.TrimSpace(query)

	if query == "" {
		return nil, fmt.Errorf("empty search query")
	}

	return c.webSearch.Search(ctx, query, 5)
}

// enrichPromptWithWebSearch добавляет результаты веб-поиска в промпт
func (c *OllamaClient) enrichPromptWithWebSearch(prompt string, results []SearchResult) string {
	var webInfo strings.Builder
	webInfo.WriteString("\n\n[WEB_SEARCH_RESULTS]\n")
	webInfo.WriteString("Я выполнил поиск в интернете по вашему запросу. Вот найденная информация:\n\n")

	for i, result := range results {
		webInfo.WriteString(fmt.Sprintf("Результат %d:\n", i+1))
		webInfo.WriteString(result.Format())
		webInfo.WriteString("\n\n")
	}

	webInfo.WriteString("Используйте эту информацию для ответа на вопрос пользователя. Если информация не найдена, используйте свои знания.\n")

	// Вставляем результаты веб-поиска перед последним запросом пользователя
	parts := strings.Split(prompt, "USER:")
	if len(parts) > 1 {
		prompt = parts[0] + webInfo.String() + "\nUSER:\n" + parts[len(parts)-1]
	} else {
		prompt = prompt + webInfo.String()
	}

	return prompt
}

// Проверка, что OllamaClient реализует интерфейс domain.LLM
var _ domain.LLM = (*OllamaClient)(nil)

