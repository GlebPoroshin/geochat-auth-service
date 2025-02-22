FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /auth-service ./cmd/auth-service

FROM alpine:latest
WORKDIR /app
COPY --from=builder /auth-service .
EXPOSE 8081
CMD ["./auth-service"] 