FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache build-base

COPY app/go.mod app/go.sum ./
RUN go mod download

COPY app .

RUN CGO_ENABLED=1 GOOS=linux go build -o /app/main .


# ---------------- runtime ----------------
FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/main /app/main

EXPOSE 8080

CMD ["/app/main"]
