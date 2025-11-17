package llm

import (
	"backend/internal/domain"
	"errors"
	"testing"
)

func TestNewChatService(t *testing.T) {
	limits := &domain.Limits{
		MaxPromptChars:  1000,
		MaxHistoryChars: 500,
		MaxRequestChars: 200,
	}

	tests := []struct {
		name      string
		chatRepo  domain.ChatRepo
		msgRepo   domain.MessageRepo
		llmClient domain.LLM
		limits    *domain.Limits
		wantErr   error
	}{
		{
			name:      "nil chatRepo",
			chatRepo:  nil,
			msgRepo:   nil,
			llmClient: nil,
			limits:    limits,
			wantErr:   errors.New("chat repo should be provided"),
		},
		{
			name:      "nil msgRepo",
			chatRepo:  struct{ domain.ChatRepo }{}, // или замени на свою заглушку
			msgRepo:   nil,
			llmClient: nil,
			limits:    limits,
			wantErr:   errors.New("message repo should be provided"),
		},
		{
			name:      "nil llm",
			chatRepo:  struct{ domain.ChatRepo }{},
			msgRepo:   struct{ domain.MessageRepo }{},
			llmClient: nil,
			limits:    limits,
			wantErr:   errors.New("LLM should be provided"),
		},
		{
			name:      "nil limits",
			chatRepo:  struct{ domain.ChatRepo }{},
			msgRepo:   struct{ domain.MessageRepo }{},
			llmClient: struct{ domain.LLM }{},
			limits:    nil,
			wantErr:   errors.New("limits should be provided"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewChatService(tt.chatRepo, tt.msgRepo, tt.llmClient, tt.limits)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if err.Error() != tt.wantErr.Error() {
				t.Fatalf("error = %q, want %q", err.Error(), tt.wantErr.Error())
			}
		})
	}
}
