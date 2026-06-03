# Деплой на VPS с Traefik

Единый файл настроек: **`config.yml`** в корне репозитория (backend, PostgreSQL, Traefik).

## config.yml перед деплоем

```yaml
addr: 0.0.0.0:8080
secret: <длинный_секрет_jwt>
database:
  host: postgres          # имя сервиса в docker-compose
  port: 5432
  username: postgres
  password: <пароль>
  dbname: guess_the_flag
deploy:
  domain: flags.example.com   # DNS → VPS
  traefik_network: traefik    # как в docker network ls
  traefik_entrypoint: websecure
  traefik_cert_resolver: letsencrypt
```

Секция `deploy` backend не читает — только для генерации переменных Docker Compose.

## Запуск

**Linux (VPS):**

```bash
docker compose -f docker-compose.prod.yml \
  --env-file <(go run ./backend/cli compose-env --config=config.yml) \
  up -d --build
```

**Через Make:**

```bash
make prod-up
```

На VPS без Go можно использовать образ backend:

```bash
docker run --rm -v "$(pwd)/config.yml:/config.yml:ro" \
  guess_the_flag_backend:latest compose-env --config=/config.yml > .compose.env
docker compose -f docker-compose.prod.yml --env-file .compose.env up -d --build
rm .compose.env
```

## Проверка

```bash
go run ./backend/cli compose-env --config=config.yml
make prod-config
```

## Локальная разработка

`docker compose up` — по-прежнему `config.docker.yml` (порты 3000/8080).

Локальный `go run` — `config.yml` с `database.host: localhost`.

## Обновление

```bash
git pull
make prod-up
```
