# Деплой на VPS с Traefik

## Одна команда

```bash
docker compose -f docker-compose.prod.yml up -d --build
```

Go, `.env` и `make` не нужны.

## Перед первым запуском

### 1. `config.docker.yml`

Настройки backend и пароль БД (должен совпадать с `POSTGRES_PASSWORD` в `docker-compose.prod.yml`, по умолчанию `postgres`):

```yaml
secret: <длинный_jwt_секрет>
database:
  password: postgres
```

### 2. Домен в `docker-compose.prod.yml`

В сервисе `frontend` → `labels` замените:

```yaml
- traefik.http.routers.guess-the-flag.rule=Host(`flags.example.com`)
```

на ваш домен. DNS (A-запись) должен указывать на VPS.

### 3. Traefik

- Запущен Traefik в Docker
- External network `traefik` (если у вас другое имя — измените `networks.traefik.name`)
- Entrypoint `websecure` и cert resolver `letsencrypt` (как в вашем Traefik)

### 4. Смена пароля БД

Одновременно поменяйте в **трёх** местах:

- `config.docker.yml` → `database.password`
- `docker-compose.prod.yml` → `postgres.environment.POSTGRES_PASSWORD`
- `docker-compose.prod.yml` → `migrate.environment.GOOSE_DBSTRING` (пароль в строке)

## Обновление

```bash
git pull
docker compose -f docker-compose.prod.yml up -d --build
```

## Локально (без Traefik)

```bash
docker compose up -d --build
```
