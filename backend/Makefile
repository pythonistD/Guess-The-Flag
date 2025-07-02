# 3. Запуск основного приложения (локально)
run: ## Запустить основное приложение локально
	@echo "$(GREEN)Запуск основного приложения...$(NC)"
	cd cmd && go run server/main.go

lines: ## Подсчитать строки кода в проекте
	@echo "$(GREEN)Статистика проекта:$(NC)"
	@find . -name "*.go" -not -path "./vendor/*" | xargs wc -l 2>/dev/null | tail -1 | awk '{print "Всего строк Go кода: " $$1}'
	@echo "Go файлов: $$(find . -name "*.go" -not -path "./vendor/*" | wc -l | tr -d ' ')"

