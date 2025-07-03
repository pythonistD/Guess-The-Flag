# 🏁 Guess The Flag

Полнофункциональная игра "Угадай флаг" с современным веб-интерфейсом. Проект состоит из Go backend API и React frontend приложения.

## 🎮 Описание игры

**Guess The Flag** - это увлекательная игра, где игроки должны угадывать страны по их флагам. Игра включает:

- Систему регистрации и аутентификации пользователей
- Игровые сессии с вопросами о флагах
- Подсчет правильных ответов и статистику
- Красивый современный интерфейс

## 🏗️ Архитектура проекта

```
Guess-The-Flag/
├── backend/                 # Go API сервер
│   ├── cmd/server/         # Точка входа приложения
│   ├── internal/           # Внутренняя логика
│   │   ├── api/           # HTTP обработчики и middleware
│   │   ├── db/            # Работа с базой данных
│   │   ├── service/       # Бизнес-логика
│   │   └── config/        # Конфигурация
│   └── migrations/        # SQL миграции
├── frontend/               # React приложение
│   ├── src/
│   │   ├── components/    # React компоненты
│   │   ├── services/      # API сервисы
│   │   └── types/         # TypeScript типы
│   └── public/            # Статические файлы
└── config.yml             # Конфигурация приложения
```

## 🚀 Быстрый старт

### Предварительные требования

- **Go** 1.21 или выше
- **Node.js** 14 или выше
- **PostgreSQL** 12 или выше
- **Git**

### 1. Клонирование репозитория

```bash
git clone <repository-url>
cd Guess-The-Flag
```

### 2. Настройка базы данных

Создайте базу данных PostgreSQL:

```sql
CREATE DATABASE guess_the_flag;
CREATE USER postgres WITH PASSWORD 'postgres';
GRANT ALL PRIVILEGES ON DATABASE guess_the_flag TO postgres;
```

### 3. Запуск Backend

```bash
# Переходим в директорию backend
cd backend

# Устанавливаем зависимости
go mod download

# Запускаем миграции (если необходимо)
# Добавьте данные стран в таблицу countries

# Запускаем сервер
go run cmd/server/main.go --config=../config.yml
```

Backend будет доступен на `http://localhost:8080`

### 4. Запуск Frontend

```bash
# Открываем новый терминал
cd frontend

# Устанавливаем зависимости
npm install

# Запускаем development сервер
npm start
```

Frontend будет доступен на `http://localhost:3000`

## 🔧 Конфигурация

### Backend (config.yml)

```yaml
addr: localhost:8080
secret: secret_key
token_lifetime: 24h
database:
  host: localhost
  port: 5432
  username: postgres
  password: postgres
  dbname: guess_the_flag
  max_idle_connections: 2
  max_open_connections: 100
  max_connection_lifetime: 60m
```

### Аргументы командной строки

Backend поддерживает следующие аргументы:

```bash
go run cmd/server/main.go --config=path/to/config.yml --addr=:8081 --log-level=debug
```

## 📱 Пользовательский интерфейс

### Экраны приложения

1. **🔐 Аутентификация** (`/login`)
   - Регистрация новых пользователей
   - Вход в систему
   - Валидация данных

2. **🏠 Главное меню** (`/game`)
   - Приветственный экран
   - Описание игры
   - Кнопка начала игры

3. **🎯 Игровой процесс** (`/game/play`)
   - Отображение флага страны
   - Поле ввода ответа
   - Обратная связь (правильно/неправильно)
   - Переход к следующему вопросу

4. **📊 Результаты** (`/game/results`)
   - Общая статистика игры
   - Детальный разбор ответов
   - Возможность начать новую игру

## 🔒 Аутентификация и безопасность

- **JWT токены** для аутентификации
- **Bcrypt** для хеширования паролей
- **Middleware** для защиты маршрутов
- **Автоматическое** обновление токенов
- **CORS** поддержка для cross-origin запросов

## 📊 API Endpoints

### Аутентификация
- `POST /auth/register` - Регистрация пользователя
- `POST /auth/login` - Вход в систему

### Игра
- `POST /game/start` - Начать новую игру
- `POST /game/question` - Получить вопрос
- `POST /game/answer` - Отправить ответ
- `POST /game/end` - Завершить игру

## 🗄️ База данных

### Основные таблицы

- `users` - Пользователи системы
- `games` - Игровые сессии
- `countries` - Страны и их флаги
- `questions` - Вопросы в играх
- `answers` - Ответы пользователей

### Миграции

Миграции находятся в `backend/migrations/` и должны выполняться при первом запуске.

## 🧪 Тестирование

### Backend тесты

```bash
cd backend
go test ./...
```

### Frontend тесты

```bash
cd frontend
npm test
```

## 🐛 Отладка

### Backend логирование

Backend поддерживает различные уровни логирования:

```bash
go run cmd/server/main.go --log-level=debug
```

### Frontend debugging

- React DevTools для отладки компонентов
- Network tab для мониторинга API запросов
- Console логи для ошибок

## 📦 Деплой

### Production сборка Frontend

```bash
cd frontend
npm run build
```

### Backend production

```bash
cd backend
go build -o guess-the-flag cmd/server/main.go
./guess-the-flag --config=config.yml
```

## 🛠️ Разработка

### Добавление новых стран

1. Добавьте запись в таблицу `countries`:
```sql
INSERT INTO countries (name, code, flag_url) 
VALUES ('Russia', 'RU', 'https://example.com/flags/ru.png');
```

### Расширение API

1. Добавьте новые handlers в `internal/api/handlers/`
2. Обновите router в `internal/api/v1/router.go`
3. Добавьте соответствующие service методы

### Стилизация Frontend

Проект использует Styled Components для стилизации. Основные цвета:

- Основной градиент: `#667eea` → `#764ba2`
- Успех: `#28a745`
- Ошибка: `#dc3545`

## 📋 TODO

- [ ] Добавить таблицу лидеров
- [ ] Реализовать различные уровни сложности
- [ ] Добавить звуковые эффекты
- [ ] Создать мобильное приложение
- [ ] Добавить мультиплеер режим
- [ ] Интеграция с социальными сетями
- [ ] PWA поддержка
- [ ] Темная тема
- [ ] Интернационализация

## 🤝 Участие в разработке

1. Fork проекта
2. Создайте feature branch (`git checkout -b feature/amazing-feature`)
3. Commit изменения (`git commit -m 'Add amazing feature'`)
4. Push в branch (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

## 📄 Лицензия

Этот проект лицензирован под MIT License - см. файл [LICENSE](LICENSE) для деталей.

## 👥 Авторы

- **Backend** - Go API with PostgreSQL
- **Frontend** - React with TypeScript
- **Design** - Modern gradient UI/UX

## 📞 Поддержка

Если у вас есть вопросы или проблемы:

1. Проверьте [Issues](../../issues)
2. Создайте новый Issue с подробным описанием
3. Укажите версию и платформу

---

**Наслаждайтесь игрой! 🎮🏁**