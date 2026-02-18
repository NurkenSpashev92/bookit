FROM golang:1.26-alpine AS builder
WORKDIR /app

RUN apk update && apk add --no-cache git bash build-base

# Установим migrate в /go/bin/migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY app/go.mod app/go.sum ./
RUN go mod tidy

COPY app .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main .


# ---------------- runtime ----------------
FROM alpine:latest
RUN apk add --no-cache netcat-openbsd ca-certificates

WORKDIR /app

# копируем приложение
COPY --from=builder /app/main /app/main
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# копируем миграции
COPY app/migrations /app/migrations

# entrypoint
COPY devops/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/app/main"]
