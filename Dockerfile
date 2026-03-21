FROM golang:1.25-alpine AS builder
WORKDIR /app

# Отдельный слой для зависимостей — кэшируется если go.mod не менялся
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s' -o server .

# Финальный образ без Go (~10MB вместо ~300MB)
FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8083
CMD ["./server"]


