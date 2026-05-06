MIGRATION_DIR=migrations
# Драйвер базы данных
DB_DRIVER=postgres
# Строка подключения (должна быть задана в .env как DATABASE_URL=...)
DB_URL ?= "postgres://postgres:postgres@localhost:5432/guess_the_flag?sslmode=disable"

.PHONY: migrate-up migrate-down migrate-status migrate-create

migrate-up:
	@echo "Running goose up..."
	goose -dir $(MIGRATION_DIR) $(DB_DRIVER) $(DB_URL) up

migrate-down:
	@echo "Running goose down..."
	goose -dir $(MIGRATION_DIR) $(DB_DRIVER) $(DB_URL) down

migrate-status:
	@echo "Migration status:"
	goose -dir $(MIGRATION_DIR) $(DB_DRIVER) $(DB_URL) status


run: ## Запустить основное приложение локально
	@echo "$(GREEN)Запуск основного приложения...$(NC)"
	cd cmd && go run server/main.go

lines: ## Подсчитать строки кода в проекте
	@echo "$(GREEN)Статистика проекта:$(NC)"
	@find . -name "*.go" -not -path "./vendor/*" | xargs wc -l 2>/dev/null | tail -1 | awk '{print "Всего строк Go кода: " $$1}'
	@echo "Go файлов: $$(find . -name "*.go" -not -path "./vendor/*" | wc -l | tr -d ' ')"

