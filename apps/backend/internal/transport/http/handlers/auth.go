package handlers

import (
	"backend/internal/domain"
	"backend/internal/transport/http/dto"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ContextKey - тип для ключей контекста
type ContextKey string

const UserIDKey ContextKey = "user_id"

// AuthHandler - handler для аутентификации
type AuthHandler struct {
	userRepo domain.UserRepo
}

func NewAuthHandler(userRepo domain.UserRepo) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
	}
}

// Login обрабатывает запрос на аутентификацию пользователя
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		http.Error(w, "password is required", http.StatusBadRequest)
		return
	}

	// Получаем пользователя по email
	user, err := h.userRepo.GetByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "failed to authenticate", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if !user.IsActive {
		http.Error(w, "account is disabled", http.StatusForbidden)
		return
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Обновляем время последнего входа
	if err := h.userRepo.UpdateLastLogin(r.Context(), user.ID); err != nil {
		// Логируем ошибку, но не прерываем процесс аутентификации
		_ = err
	}

	// Генерируем простой токен (в реальности здесь должен быть JWT)
	// TODO: заменить на JWT токен
	token := generateSimpleToken(user.ID)

	response := dto.LoginResponse{
		Token: token,
	}
	response.User.ID = user.ID.String()
	response.User.Email = user.Email

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// generateSimpleToken генерирует простой токен для аутентификации
// TODO: заменить на JWT токен
func generateSimpleToken(userID uuid.UUID) string {
	tokenData := map[string]interface{}{
		"user_id": userID.String(),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	tokenJSON, _ := json.Marshal(tokenData)
	return base64.URLEncoding.EncodeToString(tokenJSON)
}

// AuthMiddleware - middleware для аутентификации
// TODO: реализовать проверку JWT токена
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: извлечь токен из заголовка Authorization
		// TODO: проверить токен и извлечь userID
		// Пока используем заглушку
		userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// getUserIDFromContext извлекает userID из контекста
func getUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, http.ErrMissingFile
	}
	return userID, nil
}
