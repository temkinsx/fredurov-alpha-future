package handlers

import (
	"alpha_future_fredurov/apps/backend/internal/transport/http/dto"
	"encoding/json"
	"net/http"
)

type ScenariosHandler struct{}

func NewScenariosHandler() *ScenariosHandler {
	return &ScenariosHandler{}
}

// GetScenarios возвращает список доступных сценариев
func (h *ScenariosHandler) GetScenarios(w http.ResponseWriter, r *http.Request) {
	response := dto.ScenariosResponse{
		Scenarios: []dto.Scenario{
			{
				Code:        "contract_helper",
				Title:       "Помощь с договорами",
				Description: "Объяснить условия, выделить риски, подготовить формулировки",
			},
			{
				Code:        "marketing",
				Title:       "Маркетинг",
				Description: "Посты, акции, тексты для микробизнеса",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
