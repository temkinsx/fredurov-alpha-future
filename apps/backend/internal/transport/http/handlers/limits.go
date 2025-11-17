package handlers

import (
	"backend/internal/domain"
	"backend/internal/transport/http/dto"
	"encoding/json"
	"net/http"
)

type LimitsHandler struct {
	limits domain.Limits
}

func NewLimitsHandler(limits domain.Limits) *LimitsHandler {
	return &LimitsHandler{
		limits: limits,
	}
}

// GetLimits возвращает конфигурацию лимитов
func (h *LimitsHandler) GetLimits(w http.ResponseWriter, r *http.Request) {
	response := dto.LimitsResponse{
		MaxFileSizeBytes: h.limits.MaxFileSizeBytes,
		MaxFileTextChars: h.limits.MaxFileTextChars,
		MaxHistoryChars:  h.limits.MaxHistoryChars,
		MaxPromptChars:   h.limits.MaxPromptChars,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
