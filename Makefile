DATABASE_URL = postgresql://postgres:postgres@localhost:5432/restapi_todo?sslmode=disable
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
	@echo "Применение миграций..."
	@migrate -database "$(DATABASE_URL)" -path $(MIGRATION_DIR) up
	@echo "✓ Миграции применены"

.PHONY: down
down:
	@echo "Откат последней миграции..."
	@migrate -database "$(DATABASE_URL)" -path $(MIGRATION_DIR) down 1
	@echo "✓ Миграция откачена"

.PHONY: status
status:
	@migrate -database "$(DATABASE_URL)" -path $(MIGRATION_DIR) version

.PHONY: build
build:
	go build ./cmd/server/main.go

.PHONY: test
test:
	go test -v