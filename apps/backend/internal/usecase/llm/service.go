package llm

import (
	"alpha_future_fredurov/apps/backend/internal/domain"
	"errors"
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
