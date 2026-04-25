# Удалите эти две строки:
# include .env
# export ${cat .env}

DATABASE_URL = postgresql://${POSTGRES_USERNAME}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DBNAME}?sslmode=disable
MIGRATION_DIR = migrations

MIGRATE := $(shell command -v migrate 2> /dev/null)

.PHONY: create
create:
	@if [ -z "$(NAME)" ]; then \
		echo "Ошибка: укажите NAME. Пример: make create NAME=create_users_table"; \
		exit 1; \
	fi
	@mkdir -p $(MIGRATION_DIR)
	@migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(NAME)

.PHONY: up
up:
	@migrate -database "$(DATABASE_URL)" -path $(MIGRATION_DIR) up

.PHONY: down
down:
	@migrate -database "$(DATABASE_URL)" -path $(MIGRATION_DIR) down 1

.PHONY: status
status:
	@migrate -database "$(DATABASE_URL)" -path $(MIGRATION_DIR) version

.PHONY: build
build:
	echo "Server building..."
	go build -o main ./cmd/server/
	echo "Server built successfully"

.PHONY: build-run
build-run: build
	@./main

.PHONY: test
test:
	go test -v