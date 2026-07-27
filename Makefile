APP_NAME=app
COMPOSE=docker compose
ENV_FILE=.env

-include $(ENV_FILE)
export

# читаем DEBUG из .env
DEBUG := $(shell grep -E '^DEBUG=' $(ENV_FILE) 2>/dev/null | cut -d '=' -f2 | tr '[:upper:]' '[:lower:]')

# выбираем compose файлы
ifeq ($(DEBUG),true)
	COMPOSE_FILES=-f docker-compose.yml -f docker-compose.dev.yml
	MODE=DEV
else
	COMPOSE_FILES=-f docker-compose.yml -f docker-compose.prod.yml
	MODE=PROD
endif

.PHONY: help up down build restart logs ps app postgres clean prune install mode test test-v test-cover \
        migrate-up migrate-down migrate-version migrate-force migrate-create migrate-drop migrate-action seed install-stack

help:
	@echo ""
	@echo "Mode: $(MODE)"
	@echo ""
	@echo "Available commands:"
	@echo "  make install          🚀 Deploy project (нужен готовый .env)"
	@echo "  make start            🚀 Start containers"
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
	@echo ""
	@echo "  make migrate-up                ⬆️  Apply all pending migrations"
	@echo "  make migrate-down              ⬇️  Rollback ONE migration"
	@echo "  make migrate-version           🔢 Show current migration version"
	@echo "  make migrate-force V=<n>       🛠  Force version (fix dirty state)"
	@echo "  make migrate-create NAME=<x>   📝 Create new migration files"
	@echo "  make migrate-drop              💥 Drop everything (DANGEROUS)"
	@echo ""
	@echo "  make seed             🌱 Seed database with demo houses"
	@echo ""
	@echo "  make clean            🧹 Remove containers + volumes"
	@echo "  make prune            💣 Docker system prune"
	@echo ""

mode:
	@echo "Running in $(MODE) mode (DEBUG=$(DEBUG))"


install:
	@if [ ! -f $(ENV_FILE) ]; then \
		echo "❌ Нет $(ENV_FILE). Скопируйте .env.example в .env и заполните значения вручную."; \
		exit 1; \
	fi
	@$(MAKE) install-stack

install-stack:
	$(COMPOSE) $(COMPOSE_FILES) build
	$(COMPOSE) $(COMPOSE_FILES) up -d postgres_db
	$(MAKE) migrate-up
	$(COMPOSE) $(COMPOSE_FILES) up -d

start:
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

# ----- migrations (через отдельный контейнер `migrate`) -----

migrate-up:
	@$(MAKE) migrate-action action=up

migrate-down:
	@$(MAKE) migrate-action action="down 1"

migrate-version:
	@$(MAKE) migrate-action action=version

migrate-force:
	@if [ -z "$(V)" ]; then echo "Usage: make migrate-force V=<version>"; exit 1; fi
	@$(MAKE) migrate-action action="force $(V)"

migrate-drop:
	@echo "⚠️  This will DROP all tables. Press Ctrl+C in 5s to cancel..."
	@sleep 5
	@$(MAKE) migrate-action action="drop -f"

migrate-create:
	@if [ -z "$(NAME)" ]; then echo "Usage: make migrate-create NAME=add_something"; exit 1; fi
	@$(COMPOSE) $(COMPOSE_FILES) run --rm --no-deps migrate \
		create -ext sql -dir /migrations -seq $(NAME)

migrate-action:
	@$(COMPOSE) $(COMPOSE_FILES) run --rm migrate \
		-path /migrations \
		-database postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):5432/$(POSTGRES_DB)?sslmode=disable \
		$(action)

seed:
ifeq ($(DEBUG),true)
	@$(COMPOSE) $(COMPOSE_FILES) exec app go run ./cmd/seed
else
	@echo "❌ Seed доступен только в DEV режиме (DEBUG=true)."
	@exit 1
endif

clean:
	$(COMPOSE) $(COMPOSE_FILES) down -v --remove-orphans

prune:
	@docker system prune -af --volumes
