FROM golang:1.24-alpine as builder

# Установка переменной окружения для использования прямого подключения к GitHub
ENV GOPROXY=direct
ENV GO111MODULE=on

# Установка git
RUN apk add --no-cache git

WORKDIR /app

# Копирование общих модулей
COPY geochat-shared /app/geochat-shared

# Копирование исходного кода auth-service
COPY geochat-auth-service /app/geochat-auth-service

WORKDIR /app/geochat-auth-service

# Загрузка зависимостей
RUN go mod download

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service ./cmd/main.go

# Финальный этап
FROM alpine:latest

WORKDIR /

# Копирование бинарного файла из builder
COPY --from=builder /auth-service /auth-service

# Запуск приложения
CMD ["/auth-service"]