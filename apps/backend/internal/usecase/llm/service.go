package llm

import (
	"alpha_future_fredurov/apps/backend/internal/domain"
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type Service struct {
	chatRepo domain.ChatRepo
	msgRepo  domain.MessageRepo
	llm      domain.LLM
	limits   domain.Limits
}

func NewChatService(chatRepo domain.ChatRepo, msgRepo domain.MessageRepo, llm domain.LLM, limits *domain.Limits) (*Service, error) {
	if chatRepo == nil {
		return nil, errors.New("chat repo should be provided")
	}

	if msgRepo == nil {
		return nil, errors.New("message repo should be provided")
	}

	if llm == nil {
		return nil, errors.New("LLM should be provided")
	}

	if limits == nil {
		return nil, errors.New("limits should be provided")
	}

	return &Service{
		chatRepo: chatRepo,
		msgRepo:  msgRepo,
		llm:      llm,
		limits:   *limits,
	}, nil
}

func (s *Service) Reply(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, userText string, documentIDs []uuid.UUID, scenarioCode *string) (*domain.Message, error) {
	chat, err := s.chatRepo.Get(ctx, chatID)
	if err != nil {
		return nil, err
	}

	if chat.UserID != userID {
		return nil, errors.New("wrong userID for this chat")
	}

	rawMsgHistory, err := s.msgRepo.GetLastN(ctx, chatID, 10)
	if err != nil {
		return nil, err
	}

	var msgHistory strings.Builder
	for _, msg := range rawMsgHistory {
		msgHistory.WriteString(msg.String())
		msgHistory.WriteString("\n")
	}

	return nil, nil
}
