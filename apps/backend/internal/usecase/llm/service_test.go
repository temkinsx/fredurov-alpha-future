package llm

import (
	"alpha_future_fredurov/apps/backend/internal/domain"
	"reflect"
	"testing"
)

func TestNewChatService(t *testing.T) {
	type args struct {
		chatRepo domain.ChatRepo
		msgRepo  domain.MessageRepo
		llm      domain.LLM
		limits   *domain.Limits
	}
	tests := []struct {
		name    string
		args    args
		want    *Service
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewChatService(tt.args.chatRepo, tt.args.msgRepo, tt.args.llm, tt.args.limits)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewChatService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewChatService() got = %v, want %v", got, tt.want)
			}
		})
	}
}
