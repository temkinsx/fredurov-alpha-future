package chat

import (
	"alpha_future_fredurov/apps/backend/internal/domain"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	chatRepo domain.ChatRepo
	msgRepo  domain.MessageRepo
	llm      domain.LLM
}

func NewChatService(chatRepo domain.ChatRepo, msgRepo domain.MessageRepo, llm domain.LLM) (*Service, error) {
	if chatRepo == nil {
		return nil, errors.New("chat repo should be provided")
	}

	if msgRepo == nil {
		return nil, errors.New("message repo should be provided")
	}

	if llm == nil {
		return nil, errors.New("LLM should be provided")
	}

	return &Service{
		chatRepo: chatRepo,
		msgRepo:  msgRepo,
		llm:      llm,
	}, nil
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, title *string) (*domain.Chat, error) {
	chat := &domain.Chat{
		ID:        uuid.New(),
		Title:     title,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.chatRepo.Create(ctx, chat); err != nil {
		return nil, err
	}

	return chat, nil
}

func (s *Service) ListAll(ctx context.Context, userID uuid.UUID) ([]*domain.Chat, error) {
	return s.chatRepo.ListByUser(ctx, userID)
}

func (s *Service) ListMessages(ctx context.Context, chatID uuid.UUID) ([]*domain.Message, error) {

}
