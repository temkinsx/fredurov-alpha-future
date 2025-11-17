package http

import (
	"alpha_future_fredurov/apps/backend/internal/domain"
	"alpha_future_fredurov/apps/backend/internal/transport/http/handlers"
	"alpha_future_fredurov/apps/backend/internal/usecase/llm"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	chatRepo      domain.ChatRepo
	msgRepo       domain.MessageRepo
	docRepo       domain.DocumentRepo
	llmService    *llm.Service
	docTextGetter llm.DocumentTextGetter
	limits        domain.Limits
}

func NewRouter(
	chatRepo domain.ChatRepo,
	msgRepo domain.MessageRepo,
	docRepo domain.DocumentRepo,
	llmService *llm.Service,
	docTextGetter llm.DocumentTextGetter,
	limits domain.Limits,
) *Router {
	return &Router{
		chatRepo:      chatRepo,
		msgRepo:       msgRepo,
		docRepo:       docRepo,
		llmService:    llmService,
		docTextGetter: docTextGetter,
		limits:        limits,
	}
}

func (r *Router) SetupRoutes() *chi.Mux {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(handlers.AuthMiddleware)

	// Handlers
	healthHandler := handlers.NewHealthHandler()
	chatsHandler := handlers.NewChatsHandler(r.chatRepo)
	messagesHandler := handlers.NewMessagesHandler(r.msgRepo, r.chatRepo, r.llmService, r.docTextGetter)
	documentsHandler := handlers.NewDocumentsHandler(r.docRepo, r.limits)
	scenariosHandler := handlers.NewScenariosHandler()
	limitsHandler := handlers.NewLimitsHandler(r.limits)
	ragHandler := handlers.NewRAGHandler()

	// Routes
	router.Get("/health", healthHandler.Health)

	// Chats
	router.Get("/chats", chatsHandler.GetChats)
	router.Post("/chats", chatsHandler.CreateChat)

	// Messages
	router.Get("/chats/{chat_id}/messages", messagesHandler.GetMessages)
	router.Post("/chats/{chat_id}/messages", messagesHandler.SendMessage)

	// Documents
	router.Post("/documents", documentsHandler.CreateDocument)
	router.Get("/documents", documentsHandler.GetDocuments)
	router.Get("/documents/{id}", documentsHandler.GetDocument)

	// Scenarios
	router.Get("/scenarios", scenariosHandler.GetScenarios)

	// Config
	router.Get("/config/limits", limitsHandler.GetLimits)

	// RAG (опциональный)
	router.Post("/rag/search", ragHandler.Search)

	return router
}
