package middleware
/*
import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// каждому запросу даёт уникальный ID
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		meta := &RequestMeta{
			RequestID: uuid.NewString(), //генерируем уникальный ID
		}

		ctx := context.WithValue(r.Context(), RequestMetaKey, meta) //кладём ID в context
		next.ServeHTTP(w, r.WithContext(ctx))                       //передаём дальше, но уже с обновлённым context
	})
}
*/