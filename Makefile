CONFIG ?= config.yml
COMPOSE_PROD = docker compose -f docker-compose.prod.yml
COMPOSE_ENV = go run ./backend/cli compose-env --config=$(CONFIG)

.PHONY: prod-up prod-down prod-config compose-env

# Сгенерировать переменные для compose из config.yml (в stdout)
compose-env:
	$(COMPOSE_ENV)

# Показать итоговый compose-конфиг (проверка)
prod-config:
	$(COMPOSE_ENV) > .compose.env
	$(COMPOSE_PROD) --env-file .compose.env config
	@rm -f .compose.env

prod-up:
	$(COMPOSE_ENV) > .compose.env
	$(COMPOSE_PROD) --env-file .compose.env up -d --build
	@rm -f .compose.env

prod-down:
	@if [ -f .compose.env ]; then \
		$(COMPOSE_PROD) --env-file .compose.env down; \
		rm -f .compose.env; \
	else \
		$(COMPOSE_ENV) > .compose.env; \
		$(COMPOSE_PROD) --env-file .compose.env down; \
		rm -f .compose.env; \
	fi
