package llm

import (
	"backend/internal/domain"
	"strings"
	"testing"
)

func newServiceNoTrunc() *Service {
	return &Service{
		limits: domain.Limits{
			MaxPromptChars:    5000,
			MaxOutputTokens:   512,
			MaxFileSizeBytes:  5_000_000,
			MaxFileTextChars:  2000,
			MaxHistoryChars:   2000,
			MaxRequestChars:   1000,
			MaxRequestsPerMin: 60,
			MaxConcurrentLLM:  4,
		},
	}
}

func newServiceTightLimits() *Service {
	return &Service{
		limits: domain.Limits{
			MaxPromptChars:    200,
			MaxOutputTokens:   512,
			MaxFileSizeBytes:  5_000_000,
			MaxFileTextChars:  80,
			MaxHistoryChars:   80,
			MaxRequestChars:   60,
			MaxRequestsPerMin: 60,
			MaxConcurrentLLM:  4,
		},
	}
}

func TestPromptBudget_Take(t *testing.T) {
	type fields struct {
		MaxTotal int
		Used     int
	}
	type args struct {
		s   string
		max int
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		want     string
		wantUsed int
	}{
		{
			name:   "empty string returns empty and does not change Used",
			fields: fields{MaxTotal: 100, Used: 10},
			args: args{
				s:   "",
				max: 50,
			},
			want:     "",
			wantUsed: 10,
		},
		{
			name:   "no remaining budget returns empty",
			fields: fields{MaxTotal: 10, Used: 10},
			args: args{
				s:   "hello",
				max: 5,
			},
			want:     "",
			wantUsed: 10,
		},
		{
			name:   "string shorter than max and remaining",
			fields: fields{MaxTotal: 100, Used: 0},
			args: args{
				s:   "hello",
				max: 10,
			},
			want:     "hello",
			wantUsed: 5,
		},
		{
			name:   "string trimmed by max",
			fields: fields{MaxTotal: 100, Used: 0},
			args: args{
				s:   "helloworld",
				max: 5,
			},
			want:     "hello",
			wantUsed: 5,
		},
		{
			name:   "string trimmed by remaining budget",
			fields: fields{MaxTotal: 10, Used: 7},
			args: args{
				s:   "abcdef",
				max: 10,
			},
			want:     "abc", // remaining = 3
			wantUsed: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &promptBudget{
				MaxTotal: tt.fields.MaxTotal,
				Used:     tt.fields.Used,
			}

			got := b.Take(tt.args.s, tt.args.max)

			if got != tt.want {
				t.Fatalf("Take() = %q, want %q", got, tt.want)
			}
			if b.Used != tt.wantUsed {
				t.Fatalf("Used = %d, want %d", b.Used, tt.wantUsed)
			}
		})
	}
}

func TestBuildPrompt_BasicSectionsAndDefaultSystem(t *testing.T) {
	svc := newServiceNoTrunc()

	tests := []struct {
		name       string
		sysPrompt  string
		history    string
		documents  string
		userReq    string
		wantChecks []string
	}{
		{
			name:      "default system prompt is used when empty",
			sysPrompt: "",
			history:   "",
			documents: "",
			userReq:   "привет",
			wantChecks: []string{
				"SYSTEM:",
				"Ты — умный и осторожный ассистент для владельцев микробизнеса.",
				"USER:",
				"привет",
			},
		},
		{
			name:      "custom system prompt is used when provided",
			sysPrompt: "кастомный системный промпт",
			history:   "",
			documents: "",
			userReq:   "хай",
			wantChecks: []string{
				"SYSTEM:",
				"кастомный системный промпт",
				"USER:",
				"хай",
			},
		},
		{
			name:      "all sections appear in correct order",
			sysPrompt: "sys",
			history:   "hist",
			documents: "doc",
			userReq:   "user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.buildPrompt(tt.sysPrompt, tt.history, tt.documents, tt.userReq)

			for _, substr := range tt.wantChecks {
				if substr == "" {
					continue
				}
				if !strings.Contains(got, substr) {
					t.Fatalf("prompt does not contain %q\nprompt:\n%s", substr, got)
				}
			}

			if tt.name == "all sections appear in correct order" {
				order := []string{"SYSTEM:", "HISTORY:", "DOCUMENTS:", "USER:"}
				lastIndex := -1
				for _, marker := range order {
					idx := strings.Index(got, marker)
					if idx == -1 {
						t.Fatalf("expected section %q in prompt, but not found\nprompt:\n%s", marker, got)
					}
					if idx < lastIndex {
						t.Fatalf("section %q appears before previous section, wrong order\nprompt:\n%s", marker, got)
					}
					lastIndex = idx
				}
			}
		})
	}
}

func TestBuildPrompt_RespectsBudgetAndLimits(t *testing.T) {
	svc := newServiceTightLimits()

	longHistory := strings.Repeat("H", 200)
	longDocs := strings.Repeat("D", 200)
	longUser := strings.Repeat("U", 200)

	tests := []struct {
		name      string
		history   string
		documents string
		userReq   string
	}{
		{
			name:      "history truncated by history limit and global budget",
			history:   longHistory,
			documents: "",
			userReq:   "ok",
		},
		{
			name:      "documents truncated by docs limit and global budget",
			history:   "",
			documents: longDocs,
			userReq:   "ok",
		},
		{
			name:      "userReq truncated by user limit and global budget",
			history:   "",
			documents: "",
			userReq:   longUser,
		},
		{
			name:      "all together still within MaxPromptChars",
			history:   longHistory,
			documents: longDocs,
			userReq:   longUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := svc.buildPrompt("", tt.history, tt.documents, tt.userReq)

			if len(prompt) > svc.limits.MaxPromptChars {
				t.Fatalf("prompt length = %d, exceeds MaxPromptChars = %d",
					len(prompt), svc.limits.MaxPromptChars)
			}

			if tt.history != "" && strings.Contains(prompt, tt.history) {
				t.Fatalf("history text not truncated as expected")
			}
			if tt.documents != "" && strings.Contains(prompt, tt.documents) {
				t.Fatalf("documents text not truncated as expected")
			}
			if tt.userReq != "" && strings.Contains(prompt, tt.userReq) {
				t.Fatalf("user request text not truncated as expected")
			}
		})
	}
}
