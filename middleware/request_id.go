package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.NewString() //генерируем уникальный ID

		ctx := context.WithValue(r.Context(), requestIDKey, requestID) //кладём ID в context
		next.ServeHTTP(w, r.WithContext(ctx)) //передаём дальше, но уже с обновлённым context
	})
}