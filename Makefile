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

# Migrations: a one-shot migrate/migrate container on the bookit_default network
# avoids depending on whether the app image happens to ship the migrate binary.
MIGRATIONS_DIR := $(CURDIR)/app/migrations
MIGRATE_NET := bookit_default
MIGRATE_IMG := migrate/migrate
MIGRATE_DB_URL := postgresql://$(shell grep -E '^POSTGRES_USER=' $(ENV_FILE) | cut -d '=' -f2):$(shell grep -E '^POSTGRES_PASSWORD=' $(ENV_FILE) | cut -d '=' -f2)@$(shell grep -E '^POSTGRES_HOST=' $(ENV_FILE) | cut -d '=' -f2):5432/$(shell grep -E '^POSTGRES_DB=' $(ENV_FILE) | cut -d '=' -f2)?sslmode=disable
MIGRATE_RUN := docker run --rm --network $(MIGRATE_NET) -v $(MIGRATIONS_DIR):/migrations $(MIGRATE_IMG) \
	-path /migrations -database "$(MIGRATE_DB_URL)"

.PHONY: help up down build restart logs ps app postgres clean prune install mode test test-v test-cover \
        migrate-up migrate-down migrate-version migrate-force migrate-create migrate-drop

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
	@echo ""
	@echo "  make migrate-up                ⬆️  Apply all pending migrations"
	@echo "  make migrate-down              ⬇️  Rollback ONE migration"
	@echo "  make migrate-version           🔢 Show current migration version"
	@echo "  make migrate-force V=<n>       🛠  Force version (fix dirty state)"
	@echo "  make migrate-create NAME=<x>   📝 Create new migration files"
	@echo "  make migrate-drop              💥 Drop everything (DANGEROUS)"
	@echo ""
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

# ----- migrations -----

migrate-up:
	$(MIGRATE_RUN) up

migrate-down:
	$(MIGRATE_RUN) down 1

migrate-version:
	$(MIGRATE_RUN) version

migrate-force:
	@if [ -z "$(V)" ]; then echo "Usage: make migrate-force V=<version>"; exit 1; fi
	$(MIGRATE_RUN) force $(V)

migrate-create:
	@if [ -z "$(NAME)" ]; then echo "Usage: make migrate-create NAME=add_something"; exit 1; fi
	docker run --rm -v $(MIGRATIONS_DIR):/migrations $(MIGRATE_IMG) \
		create -ext sql -dir /migrations -seq $(NAME)

migrate-drop:
	@echo "⚠️  This will DROP all tables. Press Ctrl+C in 5s to cancel..."
	@sleep 5
	$(MIGRATE_RUN) drop -f

clean:
	$(COMPOSE) $(COMPOSE_FILES) down -v --remove-orphans

prune:
	docker system prune -af --volumes
