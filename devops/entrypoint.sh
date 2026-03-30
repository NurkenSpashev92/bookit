#!/bin/sh

set -e

# Загружаем переменные окружения из .env
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

POSTGRES_URL="postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"

echo "Waiting for PostgreSQL at ${POSTGRES_HOST}:${POSTGRES_PORT}..."

retries=0
max_retries=30
until migrate -path /app/migrations -database "$POSTGRES_URL" version > /dev/null 2>&1; do
  retries=$((retries + 1))
  if [ "$retries" -ge "$max_retries" ]; then
    echo "ERROR: PostgreSQL is not available after ${max_retries} attempts"
    exit 1
  fi

  echo "PostgreSQL not ready, retrying ($retries/$max_retries)..."
  sleep 2
done

echo "PostgreSQL is available, running migrations..."
migrate -path /app/migrations -database "$POSTGRES_URL" up || echo "Migrations already applied or none pending"

echo "Starting application..."
exec "$@"
