package handlers

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// ContextKey - тип для ключей контекста
type ContextKey string

const UserIDKey ContextKey = "user_id"

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
