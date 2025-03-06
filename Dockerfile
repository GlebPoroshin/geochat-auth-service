FROM golang:1.24-alpine as builder

ENV GOPROXY=direct
ENV GO111MODULE=on

RUN apk add --no-cache git

WORKDIR /app

COPY geochat-shared /app/geochat-shared

COPY geochat-auth-service /app/geochat-auth-service

WORKDIR /app/geochat-auth-service

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/auth-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/auth-service /app/auth-service

CMD ["/app/auth-service"]