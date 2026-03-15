package mymiddleware

import (
	"fmt"
	"net/http"
)
func MyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request:", r.Method, r.URL.Path)
		w.Header().Set("X-App-Name", "TaskManager")

		next.ServeHTTP(w, r)

		fmt.Println("Done")
		
	})
}

/*
Она добавляет HTTP-заголовок в ответ сервера.
То есть клиент получит ответ с таким header:
X-App-Name: TaskManager
Это просто пример того, что middleware может менять ответ.
*/