package main

import (
	"backend/internal/adapters/db/postgres"
	llmadapter "backend/internal/adapters/llm"
	"backend/internal/domain"
	httptransport "backend/internal/transport/http"
	"backend/internal/usecase/llm"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()

	// Инициализация базы данных
	pool, err := postgres.NewPool(ctx)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer pool.Close()

	// Инициализация репозиториев
	chatRepo := postgres.NewChatRepo(pool)
	msgRepo := postgres.NewMessageRepo(pool)
	userRepo := postgres.NewUserRepo(pool)

	// TODO: реализовать DocumentRepo
	docRepo := &mockDocumentRepo{} // заглушка

	// Инициализация LLM клиента
	baseURL := getenv("LLM_BASE_URL", "http://ollama:11434")
	modelName := getenv("LLM_MODEL", "llama3.1:8b")

	temp := parseFloatEnv("LLM_TEMPERATURE", 0.2)
	topP := parseFloatEnv("LLM_TOP_P", 0.9)
	maxTokens := parseIntEnv("LLM_MAX_TOKENS", 512)
	enableWebSearch := getenv("LLM_ENABLE_WEB_SEARCH", "true") == "true"

	httpClient := &http.Client{Timeout: 30 * time.Second}

	llmClient := llmadapter.NewOllamaClient(httpClient, llmadapter.Config{
		BaseURL:         baseURL,
		Model:           modelName,
		Temperature:     float32(temp),
		TopP:            float32(topP),
		MaxTokens:       maxTokens,
		Timeout:         30 * time.Second,
		EnableWebSearch: enableWebSearch,
	})

	// Проверка, что llmClient реализует интерфейс domain.LLM
	var _ domain.LLM = llmClient

	// Тестовый запрос к LLM (опционально, можно убрать в production)
	if os.Getenv("LLM_TEST_ON_STARTUP") == "true" {
		resp, err := llmClient.Generate(ctx, []byte("Say hello in Russian"))
		if err != nil {
			log.Printf("LLM test failed: %v", err)
		} else {
			log.Printf("LLM test successful: %s", resp)
		}
	}

	// Инициализация лимитов
	limits := domain.Limits{
		MaxPromptChars:    parseIntEnv("LIMIT_MAX_PROMPT_CHARS", 10000),
		MaxOutputTokens:   parseIntEnv("LIMIT_MAX_OUTPUT_TOKENS", 2048),
		MaxFileSizeBytes:  parseIntEnv("LIMIT_MAX_FILE_SIZE_BYTES", 10*1024*1024), // 10MB
		MaxFileTextChars:  parseIntEnv("LIMIT_MAX_FILE_TEXT_CHARS", 100000),
		MaxHistoryChars:   parseIntEnv("LIMIT_MAX_HISTORY_CHARS", 50000),
		MaxRequestChars:   parseIntEnv("LIMIT_MAX_REQUEST_CHARS", 50000),
		MaxRequestsPerMin: parseIntEnv("LIMIT_MAX_REQUESTS_PER_MIN", 60),
		MaxConcurrentLLM:  parseIntEnv("LIMIT_MAX_CONCURRENT_LLM", 5),
	}

	// Инициализация LLM сервиса
	llmService, err := llm.NewChatService(chatRepo, msgRepo, llmClient, &limits)
	if err != nil {
		log.Fatalf("failed to create LLM service: %v", err)
	}

	// TODO: реализовать DocumentTextGetter
	docTextGetter := &mockDocumentTextGetter{} // заглушка

	// Инициализация HTTP роутера
	router := httptransport.NewRouter(
		chatRepo,
		msgRepo,
		docRepo,
		userRepo,
		llmService,
		docTextGetter,
		limits,
	)

	mux := router.SetupRoutes()

	addr := serverAddr()
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		log.Printf("HTTP server is starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server failed: %v", err)
		}
		return
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}

	log.Println("server stopped gracefully")
}

func serverAddr() string {
	if addr := os.Getenv("HTTP_ADDR"); addr != "" {
		return addr
	}

	if port := os.Getenv("PORT"); port != "" {
		if strings.HasPrefix(port, ":") {
			return port
		}
		return ":" + port
	}

	return ":8080"
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func parseFloatEnv(key string, def float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return def
	}
	return f
}

func parseIntEnv(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}

// mockDocumentRepo - заглушка для DocumentRepo
// TODO: реализовать настоящий DocumentRepo
type mockDocumentRepo struct{}

func (m *mockDocumentRepo) Create(ctx context.Context, doc *domain.Document) error {
	return errors.New("DocumentRepo not implemented")
}

func (m *mockDocumentRepo) GetByID(ctx context.Context, docID uuid.UUID) (*domain.Document, error) {
	return nil, errors.New("DocumentRepo not implemented")
}

func (m *mockDocumentRepo) Update(ctx context.Context, doc *domain.Document) (*domain.Document, error) {
	return nil, errors.New("DocumentRepo not implemented")
}

func (m *mockDocumentRepo) Delete(ctx context.Context, docID uuid.UUID) error {
	return errors.New("DocumentRepo not implemented")
}

// mockDocumentTextGetter - заглушка для DocumentTextGetter
// TODO: реализовать настоящий DocumentTextGetter
type mockDocumentTextGetter struct{}

func (m *mockDocumentTextGetter) GetDocumentText(ctx context.Context, docID uuid.UUID) (string, error) {
	return "", errors.New("DocumentTextGetter not implemented")
}
