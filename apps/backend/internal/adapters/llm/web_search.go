package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// WebSearch выполняет веб-поиск через DuckDuckGo Instant Answer API
type WebSearch struct {
	client *http.Client
}

func NewWebSearch(client *http.Client) *WebSearch {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &WebSearch{client: client}
}

// Search выполняет поиск в интернете
func (w *WebSearch) Search(ctx context.Context, query string, maxResults int) ([]SearchResult, error) {
	if maxResults <= 0 {
		maxResults = 5
	}

	// Используем DuckDuckGo Instant Answer API
	// Также можем использовать их HTML API для более полных результатов
	results, err := w.searchDuckDuckGo(ctx, query, maxResults)
	if err != nil {
		// Если DuckDuckGo не работает, попробуем альтернативные источники
		return w.searchAlternative(ctx, query, maxResults)
	}

	return results, nil
}

// searchDuckDuckGo выполняет поиск через DuckDuckGo API
func (w *WebSearch) searchDuckDuckGo(ctx context.Context, query string, maxResults int) ([]SearchResult, error) {
	// DuckDuckGo HTML API (не требует API ключа)
	searchURL := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("duckduckgo returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Парсим HTML ответ (упрощенная версия)
	results := w.parseDuckDuckGoHTML(string(body), maxResults)

	// Также пробуем Instant Answer API для более структурированных данных
	iaResults, _ := w.searchDuckDuckGoInstant(ctx, query)
	results = append(results, iaResults...)

	if len(results) > maxResults {
		results = results[:maxResults]
	}

	return results, nil
}

// searchDuckDuckGoInstant использует Instant Answer API
func (w *WebSearch) searchDuckDuckGoInstant(ctx context.Context, query string) ([]SearchResult, error) {
	apiURL := fmt.Sprintf("https://api.duckduckgo.com/?q=%s&format=json&no_html=1&skip_disambig=1", url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		Abstract     string `json:"Abstract"`
		AbstractText string `json:"AbstractText"`
		AbstractURL  string `json:"AbstractURL"`
		Answer       string `json:"Answer"`
		RelatedTopics []struct {
			Text string `json:"Text"`
			FirstURL string `json:"FirstURL"`
		} `json:"RelatedTopics"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var results []SearchResult

	if data.AbstractText != "" {
		results = append(results, SearchResult{
			Title:   data.Abstract,
			URL:     data.AbstractURL,
			Snippet: data.AbstractText,
		})
	}

	if data.Answer != "" {
		results = append(results, SearchResult{
			Title:   "Instant Answer",
			Snippet: data.Answer,
		})
	}

	for _, topic := range data.RelatedTopics {
		if topic.Text != "" && len(results) < 5 {
			results = append(results, SearchResult{
				Title: topic.Text,
				URL:   topic.FirstURL,
			})
		}
	}

	return results, nil
}

// parseDuckDuckGoHTML парсит HTML ответ от DuckDuckGo (упрощенная версия)
func (w *WebSearch) parseDuckDuckGoHTML(html string, maxResults int) []SearchResult {
	// Упрощенный парсинг - в реальности лучше использовать html парсер
	var results []SearchResult

	// Ищем ссылки и заголовки в HTML (базовая реализация)
	// В реальном проекте лучше использовать goquery или подобную библиотеку
	parts := strings.Split(html, `<a class="result__a"`)
	
	for i, part := range parts {
		if i == 0 || len(results) >= maxResults {
			continue
		}

		// Извлекаем URL
		urlStart := strings.Index(part, `href="`)
		if urlStart == -1 {
			continue
		}
		urlStart += 6
		urlEnd := strings.Index(part[urlStart:], `"`)
		if urlEnd == -1 {
			continue
		}
		resultURL := part[urlStart : urlStart+urlEnd]

		// Извлекаем заголовок
		titleStart := strings.Index(part, `>`) + 1
		titleEnd := strings.Index(part[titleStart:], `</a>`)
		if titleEnd == -1 {
			continue
		}
		title := strings.TrimSpace(part[titleStart : titleStart+titleEnd])

		// Извлекаем сниппет (описание)
		snippetStart := strings.Index(part, `<a class="result__snippet"`)
		snippet := ""
		if snippetStart != -1 {
			snippetStart = strings.Index(part[snippetStart:], `>`) + snippetStart + 1
			snippetEnd := strings.Index(part[snippetStart:], `</a>`)
			if snippetEnd != -1 {
				snippet = strings.TrimSpace(part[snippetStart : snippetStart+snippetEnd])
			}
		}

		results = append(results, SearchResult{
			Title:   title,
			URL:     resultURL,
			Snippet: snippet,
		})
	}

	return results
}

// searchAlternative альтернативный метод поиска (через другой API)
func (w *WebSearch) searchAlternative(ctx context.Context, query string, maxResults int) ([]SearchResult, error) {
	// Можно добавить другие источники, например:
	// - SerpAPI
	// - Google Custom Search API
	// - Bing Search API
	return []SearchResult{
		{
			Title:   "Web search not available",
			Snippet: "Please try again later or reformulate your query",
		},
	}, nil
}

// SearchResult представляет результат веб-поиска
type SearchResult struct {
	Title   string
	URL     string
	Snippet string
}

// Format форматирует результат для включения в промпт
func (sr SearchResult) Format() string {
	var parts []string
	if sr.Title != "" {
		parts = append(parts, fmt.Sprintf("Title: %s", sr.Title))
	}
	if sr.Snippet != "" {
		parts = append(parts, fmt.Sprintf("Content: %s", sr.Snippet))
	}
	if sr.URL != "" {
		parts = append(parts, fmt.Sprintf("URL: %s", sr.URL))
	}
	return strings.Join(parts, "\n")
}

