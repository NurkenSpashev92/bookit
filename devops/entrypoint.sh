#!/bin/sh

set -e

# Загружаем переменные окружения из .env
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

export POSTGRES_URL="postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable"

# Ожидаем, пока PostgreSQL будет доступен
echo "⏳ Ожидание запуска PostgreSQL..."
until nc -z -v -w30 $POSTGRES_HOST $POSTGRES_PORT; do
  echo "⌛ PostgreSQL ещё не доступен, ждём..."
  sleep 2
done

echo "✅ PostgreSQL доступен, запускаем миграции..."

# Запуск миграций
echo "🚀 Запуск миграций базы данных..."
migrate -path /app/migrations -database "$POSTGRES_URL" up || echo "⚠️ Миграции отсутствуют или уже применены"

echo "🎉 Миграции успешно завершены, запускаем сервис..."
exec "$@"
