APP_NAME=app
COMPOSE=docker compose
ENV_FILE=.env

# читаем DEBUG из .env
DEBUG := $(shell grep -E '^DEBUG=' $(ENV_FILE) | cut -d '=' -f2 | tr '[:upper:]' '[:lower:]')

# выбираем compose файлы
ifeq ($(DEBUG),true)
	COMPOSE_FILES=-f docker-compose.yml -f docker-compose.dev.yml
	MODE=DEV
else
	COMPOSE_FILES=-f docker-compose.yml
	MODE=PROD
endif

.PHONY: help up down build restart logs ps app postgres clean prune install mode test test-v test-cover

help:
	@echo ""
	@echo "Mode: $(MODE)"
	@echo ""
	@echo "Available commands:"
	@echo "  make install          🚀 Deploy project"
	@echo "  make up               🚀 Start containers"
	@echo "  make down             🛑 Stop containers"
	@echo "  make build            🔨 Build containers"
	@echo "  make restart          🔄 Restart containers"
	@echo "  make logs             📜 Show logs"
	@echo "  make ps               📦 Show containers"
	@echo "  make app              🐹 Enter app container"
	@echo "  make postgres         🐘 Enter postgres container"
	@echo "  make test             🧪 Run tests"
	@echo "  make test-v           🧪 Run tests (verbose)"
	@echo "  make test-cover       🧪 Run tests with coverage"
	@echo "  make clean            🧹 Remove containers + volumes"
	@echo "  make prune            💣 Docker system prune"
	@echo ""

mode:
	@echo "Running in $(MODE) mode (DEBUG=$(DEBUG))"

install: build up

up:
	$(COMPOSE) $(COMPOSE_FILES) up -d

down:
	$(COMPOSE) $(COMPOSE_FILES) down

build:
	$(COMPOSE) $(COMPOSE_FILES) build --no-cache

restart:
	$(COMPOSE) $(COMPOSE_FILES) down
	$(COMPOSE) $(COMPOSE_FILES) up -d

logs:
	$(COMPOSE) $(COMPOSE_FILES) logs -f

ps:
	$(COMPOSE) $(COMPOSE_FILES) ps

app:
	$(COMPOSE) $(COMPOSE_FILES) exec $(APP_NAME) sh

postgres:
	$(COMPOSE) $(COMPOSE_FILES) exec postgres_db psql -U $$POSTGRES_USER -d $$POSTGRES_DB

test:
	cd app && go test ./test/... -count=1

test-v:
	cd app && go test ./test/... -count=1 -v

test-cover:
	cd app && go test ./test/... -count=1 -coverprofile=coverage.out && go tool cover -func=coverage.out

clean:
	$(COMPOSE) $(COMPOSE_FILES) down -v --remove-orphans

prune:
	docker system prune -af --volumes
