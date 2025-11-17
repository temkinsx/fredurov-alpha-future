package handlers

import (
	"backend/internal/transport/http/dto"
	"encoding/json"
	"net/http"
)

type RAGHandler struct {
	// TODO: добавить RAG сервис когда будет реализован
}

func NewRAGHandler() *RAGHandler {
	return &RAGHandler{}
}

// Search выполняет поиск по документам через RAG
func (h *RAGHandler) Search(w http.ResponseWriter, r *http.Request) {
	var req dto.RAGSearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Query == "" {
		http.Error(w, "query is required", http.StatusBadRequest)
		return
	}

	if req.TopK <= 0 {
		req.TopK = 6 // дефолтное значение
	}

	// TODO: реализовать RAG поиск
	// Пока возвращаем заглушку
	response := dto.RAGSearchResponse{
		Chunks: []dto.RAGChunk{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
