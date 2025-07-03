# Guess The Flag API Documentation

## Обзор

Это API документация для игры "Угадай флаг" - веб-приложения, где пользователи могут проверить свои знания флагов стран мира.

## Архитектура API

API построено на архитектуре REST и использует JSON для обмена данными. Аутентификация реализована через JWT токены.

### Основные компоненты:

1. **Аутентификация** (`/auth/*`) - регистрация и вход пользователей
2. **Игровая механика** (`/game/*`) - управление игровым процессом

## Просмотр документации

### Swagger UI

Чтобы просмотреть интерактивную документацию в Swagger UI:

1. Откройте [Swagger Editor](https://editor.swagger.io/)
2. Скопируйте содержимое файла `swagger.yaml` в редактор
3. Или импортируйте файл напрямую

### Локальный просмотр

Если у вас установлен Swagger UI локально:

```bash
# Установка swagger-ui через npm
npm install -g swagger-ui-dist

# Или используйте Docker
docker run -p 8081:8080 -e SWAGGER_JSON=/swagger.yaml -v $(pwd)/docs:/swagger swaggerapi/swagger-ui
```

## Эндпоинты API

### Аутентификация

| Метод | Путь | Описание | Аутентификация |
|-------|------|----------|----------------|
| POST | `/auth/register` | Регистрация пользователя | Нет |
| POST | `/auth/login` | Вход в систему | Нет |

### Игра

| Метод | Путь | Описание | Аутентификация |
|-------|------|----------|----------------|
| POST | `/game/start` | Начать новую игру | JWT |
| POST | `/game/question` | Получить вопрос | JWT |
| POST | `/game/answer` | Ответить на вопрос | JWT |
| POST | `/game/end` | Завершить игру | JWT |

## Аутентификация

API использует Bearer токены для аутентификации. После успешной регистрации или входа вы получите JWT токен, который нужно передавать в заголовке:

```
Authorization: Bearer <your-jwt-token>
```

## Игровой процесс

1. **Регистрация/Вход**: Создайте аккаунт или войдите в существующий
2. **Начало игры**: Вызовите `POST /game/start` для создания новой игровой сессии
3. **Получение вопросов**: Используйте `POST /game/question` для получения флага страны
4. **Ответы**: Отправляйте ответы через `POST /game/answer`
5. **Завершение**: Вызовите `POST /game/end` для получения результатов

## Примеры использования

### Регистрация пользователя

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "securePassword123"
  }'
```

### Начало игры

```bash
curl -X POST http://localhost:8080/game/start \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json"
```

### Получение вопроса

```bash
curl -X POST http://localhost:8080/game/question \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "gameId": "e6fdfe58-4aba-4454-8283-6f5a4a18f980",
    "questionNum": 1
  }'
```

## Коды ответов

- `200` - Успешный запрос
- `400` - Неверный формат запроса
- `401` - Пользователь не аутентифицирован
- `403` - Доступ запрещен (неверные учетные данные)
- `500` - Внутренняя ошибка сервера

## Типы данных

### UUID формат
Все идентификаторы игр и пользователей используют UUID формат:
```
e6fdfe58-4aba-4454-8283-6f5a4a18f980
```

### JWT токен
Токены имеют формат JWT (JSON Web Token):
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzNDU2NzgtOTBhYi1jZGVmLTEyMzQtNTY3ODkwYWJjZGVmIiwiZXhwIjoxNjk5OTk5OTk5fQ.signature
```

## Обработка ошибок

Все ошибки возвращаются в следующем формате:

```json
{
  "error": "Описание ошибки"
}
```

## Версионирование

Текущая версия API: `v0.1.0`

API следует принципам семантического версионирования (SemVer).

## Поддержка

Если у вас есть вопросы по API, создайте issue в репозитории проекта. 