package handlers

import (
	"alpha_future_fredurov/apps/backend/internal/domain"
	"alpha_future_fredurov/apps/backend/internal/transport/http/dto"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type ChatsHandler struct {
	chatRepo domain.ChatRepo
}

func NewChatsHandler(chatRepo domain.ChatRepo) *ChatsHandler {
	return &ChatsHandler{
		chatRepo: chatRepo,
	}
}

// GetChats возвращает список чатов текущего пользователя
func (h *ChatsHandler) GetChats(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	chats, err := h.chatRepo.ListByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get chats", http.StatusInternalServerError)
		return
	}

	response := dto.ChatsListResponse{
		Chats: make([]dto.ChatResponse, len(chats)),
	}

	for i, chat := range chats {
		response.Chats[i] = dto.ChatResponse{
			ID:        chat.ID.String(),
			Title:     chat.Title,
			CreatedAt: chat.CreatedAt,
			UpdatedAt: chat.UpdatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateChat создает новый чат
func (h *ChatsHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.CreateChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	chat := &domain.Chat{
		ID:     uuid.New(),
		Title:  req.Title,
		UserID: userID,
		Status: domain.ChatActive,
	}

	if err := h.chatRepo.Create(r.Context(), chat); err != nil {
		http.Error(w, "failed to create chat", http.StatusInternalServerError)
		return
	}

	response := dto.CreateChatResponse{
		ID:    chat.ID.String(),
		Title: chat.Title,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
