package http

import (
	"backend/internal/domain"
	"backend/internal/transport/http/handlers"
	"backend/internal/usecase/llm"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	chatRepo      domain.ChatRepo
	msgRepo       domain.MessageRepo
	docRepo       domain.DocumentRepo
	userRepo      domain.UserRepo
	llmService    *llm.Service
	docTextGetter llm.DocumentTextGetter
	limits        domain.Limits
}

func NewRouter(
	chatRepo domain.ChatRepo,
	msgRepo domain.MessageRepo,
	docRepo domain.DocumentRepo,
	userRepo domain.UserRepo,
	llmService *llm.Service,
	docTextGetter llm.DocumentTextGetter,
	limits domain.Limits,
) *Router {
	return &Router{
		chatRepo:      chatRepo,
		msgRepo:       msgRepo,
		docRepo:       docRepo,
		userRepo:      userRepo,
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

	// Handlers
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(r.userRepo)
	chatsHandler := handlers.NewChatsHandler(r.chatRepo)
	messagesHandler := handlers.NewMessagesHandler(r.msgRepo, r.chatRepo, r.llmService, r.docTextGetter)
	documentsHandler := handlers.NewDocumentsHandler(r.docRepo, r.limits)
	scenariosHandler := handlers.NewScenariosHandler()
	limitsHandler := handlers.NewLimitsHandler(r.limits)
	ragHandler := handlers.NewRAGHandler()

	// Public routes (без аутентификации)
	router.Get("/health", healthHandler.Health)
	router.Post("/login", authHandler.Login)

	// Protected routes (с аутентификацией)
	router.Group(func(r chi.Router) {
		r.Use(handlers.AuthMiddleware)

		// Chats
		r.Get("/chats", chatsHandler.GetChats)
		r.Post("/chats", chatsHandler.CreateChat)

		// Messages
		r.Get("/chats/{chat_id}/messages", messagesHandler.GetMessages)
		r.Post("/chats/{chat_id}/messages", messagesHandler.SendMessage)

		// Documents
		r.Post("/documents", documentsHandler.CreateDocument)
		r.Get("/documents", documentsHandler.GetDocuments)
		r.Get("/documents/{id}", documentsHandler.GetDocument)

		// Scenarios
		r.Get("/scenarios", scenariosHandler.GetScenarios)

		// Config
		r.Get("/config/limits", limitsHandler.GetLimits)

		// RAG (опциональный)
		r.Post("/rag/search", ragHandler.Search)
	})

	return router
}
