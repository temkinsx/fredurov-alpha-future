package handlers

import (
	"alpha_future_fredurov/apps/backend/internal/domain"
	"alpha_future_fredurov/apps/backend/internal/transport/http/dto"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type DocumentsHandler struct {
	docRepo domain.DocumentRepo
	limits  domain.Limits
}

func NewDocumentsHandler(docRepo domain.DocumentRepo, limits domain.Limits) *DocumentsHandler {
	return &DocumentsHandler{
		docRepo: docRepo,
		limits:  limits,
	}
}

// CreateDocument загружает файл
func (h *DocumentsHandler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Парсим multipart/form-data
	if err := r.ParseMultipartForm(int64(h.limits.MaxFileSizeBytes)); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Проверяем размер файла
	if header.Size > int64(h.limits.MaxFileSizeBytes) {
		http.Error(w, fmt.Sprintf("file too large, max size: %d bytes", h.limits.MaxFileSizeBytes), http.StatusBadRequest)
		return
	}

	// Опциональный chat_id
	var chatID *string
	if chatIDStr := r.FormValue("chat_id"); chatIDStr != "" {
		if _, err := uuid.Parse(chatIDStr); err == nil {
			chatID = &chatIDStr
		}
	}

	// Читаем файл (в реальности здесь должна быть загрузка в S3 или другое хранилище)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		return
	}

	// TODO: извлечение текста из файла (PDF, DOCX и т.д.)
	// Пока используем заглушку
	_ = fileBytes // будет использовано для извлечения текста

	// Создаём документ
	doc := &domain.Document{
		ID:         uuid.New().String(),
		UserID:     userID.String(),
		ChatID:     chatID,
		Name:       header.Filename,
		MimeType:   header.Header.Get("Content-Type"),
		StorageURL: "", // TODO: URL после загрузки в хранилище
		SizeBytes:  header.Size,
	}

	if err := h.docRepo.Create(r.Context(), doc); err != nil {
		http.Error(w, "failed to create document", http.StatusInternalServerError)
		return
	}

	response := dto.CreateDocumentResponse{
		Document: dto.DocumentResponse{
			ID:       doc.ID,
			FileName: doc.Name,
			ChatID:   doc.ChatID,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetDocuments возвращает список документов пользователя
func (h *DocumentsHandler) GetDocuments(w http.ResponseWriter, r *http.Request) {
	_, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// TODO: реализовать фильтрацию по chat_id когда будет реализован ListByUser в DocumentRepo
	_ = r.URL.Query().Get("chat_id")

	// TODO: реализовать ListByUser в DocumentRepo
	// Пока возвращаем заглушку
	response := dto.DocumentsListResponse{
		Documents: []dto.DocumentResponse{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetDocument возвращает метаданные документа
func (h *DocumentsHandler) GetDocument(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	docIDStr := chi.URLParam(r, "id")
	docID, err := uuid.Parse(docIDStr)
	if err != nil {
		http.Error(w, "invalid document id", http.StatusBadRequest)
		return
	}

	doc, err := h.docRepo.GetByID(r.Context(), docID)
	if err != nil {
		http.Error(w, "document not found", http.StatusNotFound)
		return
	}

	// Проверяем права доступа
	if doc.UserID != userID.String() {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	// TODO: получить text_excerpt из хранилища или БД
	textExcerpt := "Настоящий договор подряда..." // заглушка

	response := dto.DocumentDetailResponse{
		ID:          doc.ID,
		FileName:    doc.Name,
		MimeType:    doc.MimeType,
		ChatID:      doc.ChatID,
		SizeBytes:   doc.SizeBytes,
		TextExcerpt: textExcerpt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
