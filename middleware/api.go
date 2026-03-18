package middleware

import (
	"net/http"
)

// проверяет, есть ли доступ
func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")

		if key != "secret123" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		meta, ok := r.Context().Value(requestMetaKey).(*RequestMeta) //считай это указателем на RequestMeta
		if ok && meta != nil {
			meta.UserID = "123"
		}
		next.ServeHTTP(w, r)
	})
}
