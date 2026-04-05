package middleware
/*
type contextKey string

type RequestMeta struct {
	RequestID string
	UserID    string
}

const (
	RequestMetaKey contextKey = "request_meta"
)


Мы добавили 4 логических вещи:

структуру RequestMeta

ключ для хранения этой структуры в context

middleware, который создаёт RequestMeta и кладёт его в context

код, который:
читает RequestMeta
дописывает в него UserID
потом логирует всё это
*/
