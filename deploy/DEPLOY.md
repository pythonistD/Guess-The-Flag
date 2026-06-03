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

### 2. Путь и домен в `docker-compose.prod.yml`

Приложение открывается по **`https://yacheboksarov.ru/game`** (как `/portfolio`), без отдельного поддомена.

В labels сервиса `frontend` при необходимости замените домен:

```yaml
- traefik.http.routers.guess-the-flag.rule=Host(`yacheboksarov.ru`) && PathPrefix(`/game`)
- traefik.http.routers.guess-the-flag.priority=100
```

`priority=100` — чтобы роутер срабатывал раньше общих правил на том же хосте.

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

## HTTP 502 от Traefik на /game

502 = Traefik не достучался до `guess_the_flag_frontend:80` (не nginx внутри).

**1. Сеть Traefik** — имя должно совпадать везде:

```bash
docker network ls | grep traefik
docker inspect guess_the_flag_frontend --format '{{range $k,$v := .NetworkSettings.Networks}}{{$k}} {{end}}'
```

В `docker-compose.prod.yml`: `networks.traefik.name` и label `traefik.docker.network` (сейчас `traefik_default`).

**2. Обход Traefik** — отвечает ли nginx:

```bash
docker exec guess_the_flag_frontend wget -qS -O- http://127.0.0.1/game/ 2>&1 | head -15
```

Если тут `200` — чините Traefik (сеть/labels). Если ошибка — пересоберите frontend.

**3. Пересоздать frontend с labels:**

```bash
docker compose -f docker-compose.prod.yml up -d --build --force-recreate frontend
```

**4. Логи Traefik:**

```bash
docker logs traefik --since 5m 2>&1 | tail -30
```

## Только postgres, нет backend/frontend

Обычно `filldb` завершился с ошибкой, и backend не стартовал (старая схема).

Заполнение БД вручную (один раз после миграций):

```bash
docker compose -f docker-compose.prod.yml run --rm --entrypoint /app/cli backend database fill --config=/app/config.yml
```

```bash
docker compose -f docker-compose.prod.yml ps -a
docker logs guess_the_flag_backend
docker compose -f docker-compose.prod.yml up -d --build backend frontend
```

## Ошибка загрузки стран при старте backend

```bash
docker logs guess_the_flag_backend
```

Частые причины: нет HTTPS до `restcountries.com` / CDN флагов, неверный `database.password` в `config.docker.yml`.
