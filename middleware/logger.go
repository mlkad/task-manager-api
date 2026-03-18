package middleware

import (
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

// логирует каждый запрос
func Logger(log zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()                                        //измерить latency (время ответа)
			ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor) //создание обёртки над http.ResponseWriter, которая узнает HTTP статус и размер ответа
			//ww  = wrapped writer

			next.ServeHTTP(ww, r) //выполняется handler
			event := log.Info().Str("method", r.Method).Str("path", r.URL.Path).Int("status", ww.Status()).Dur("duration", time.Since(start))

			//достаём ID из context
			meta, ok := r.Context().Value(requestMetaKey).(*RequestMeta)
			if ok && meta != nil {
				if meta.RequestID != "" {
					event = event.Str("request_id", meta.RequestID)
				}

				//достаём пользователя
				if meta.UserID != "" {
					event = event.Str("user_id", meta.UserID)
				}
			}
			event.Msg("request")
		})
	}
}

/*
{
  "level": "info",
  "method": "GET",
  "path": "/tasks",
  "status": 200,
  "duration": 1234567,
  "message": "request"
}
*/

// Logger принимает логгер и возвращает функцию, которая принимает handler и возвращает handler
//верни функцию (middleware), которая принимает следующий handler и возвращает новый handler (обёрнутый)

//log := zerolog.New(os.Stdout).With().Timestamp().Logger()

// zerolog.New(os.Stdout)
// создаёт логгер
// пишет в stdout (консоль / Docker / system logs)

// .With()
// включает builder для добавления полей

// .Timestamp()
// добавляет время к каждому логу

// .Logger()
// завершает сборку

// нет зависимостей → простой middleware
// есть зависимости → middleware factory
