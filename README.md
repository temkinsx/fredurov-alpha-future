# Alfa Copilot AI (alpha_future_fredurov)

Моно-репозиторий прототипа ассистента для микробизнеса. Состоит из Go backend с чат- и документными API, React/Vite frontend с UX‑макетами чат-клиента и окружения Docker Compose (Postgres + Ollama) для быстрого запуска.

## Состав репозитория

```
apps/
├── backend/   # Go 1.25 сервис: доменные сущности, use-case слой и HTTP-хендлеры
└── frontend/  # React 19 + Vite приложение с Tailwind 4 и react-router
docker-compose.yml
```

## Технологический стек

| Область     | Технологии |
|-------------|------------|
| Backend     | Go 1.25, chi, pgx/v5, bcrypt, testcontainers, Ollama API |
| Хранение    | PostgreSQL 16 (схемы `app` и `auth`), миграции в `apps/backend/migrations` |
| LLM         | Ollama (ручной выбор модели, веб-поиск через DuckDuckGo HTML API) |
| Frontend    | React 19, Vite 7, TypeScript 5.9, Tailwind CSS v4, react-auth-kit |
| Инфраструктура | Docker/Docker Compose |

## Основные возможности

- **Аутентификация** — заглушка на основе `auth.users`, авторизация по токену (пока без JWT).
- **Чаты и сообщения** — CRUD через `ChatRepo`/`MessageRepo`, история сообщений подтягивается в use-case `llm.Service`.
- **LLM-сервис** — сборка промпта (сценарии, история, документы, веб-поиск) и отправка в Ollama, учёт лимитов (`domain.Limits`).
- **Frontend** — макет страницы чата с сайдбаром чатов, быстрыми действиями, лентой диалога и формой ввода; маршрутизация `/`, `/login`, `*`.

## Быстрый старт через Docker Compose

1. Скопируйте репозиторий и убедитесь, что Docker/Docker Compose ≥ v2 установлены.
2. При первом запуске создайте каталоги для volume (`pgdata`, `ollama`) или доверьте это docker-compose.
3. Запустите все сервисы:

   ```bash
   docker compose up --build
   ```

   Сервисы:

   - `postgres` — база с healthcheck (`localhost:5432`).
   - `ollama` — LLM runtime на `localhost:11434` (после старта загрузите модель, напр. `docker exec -it ollama ollama pull mistral`).
   - `backend` — слушает `:8080`. В контейнер нужно передать `POSTGRES_DSN`.
   - `frontend` — билд Vite + nginx на `http://localhost:3000`.

4. Для остановки: `docker compose down` (при необходимости с `-v` для очистки volume).

## Локальная разработка

### Backend (Go)

1. Требования: Go ≥ 1.25 (модульный файл указывает 1.25), PostgreSQL 16, доступ к Ollama.
2. Установите переменные окружения (пример для локальной машины):

   ```bash
   export POSTGRES_DSN="postgres://postgres:postgres@127.0.0.1:5432/app?sslmode=disable"
   export HTTP_ADDR=":8080"   # необязательно, по умолчанию :8080
   ```

3. Примените миграции (любой инструмент, пример через `psql`):

   ```bash
   psql "$POSTGRES_DSN" -f apps/backend/migrations/0001_init_schemas.up.sql
   psql "$POSTGRES_DSN" -f apps/backend/migrations/0002_init_app_tables.up.sql
   psql "$POSTGRES_DSN" -f apps/backend/migrations/0003_init_auth_tables.up.sql
   ```

4. Запустите HTTP-сервер:

   ```bash
   cd apps/backend
   go run ./cmd/main
   ```

5. Тесты:

   ```bash
   go test ./...
   ```

### Frontend (React)

1. Требования: Node.js 20, npm 10.
2. Установка зависимостей:

   ```bash
   cd apps/frontend
   npm ci
   ```

3. Режим разработки:

   ```bash
   npm run dev
   ```

   Vite поднимет dev-сервер (обычно `http://localhost:5173`). Маршруты:
   - `/` — корневая страница с заглушкой.
   - `/login` — форма входа (без привязки к API).
   - прочее — NotFound.

4. Сборка/проверки:

   ```bash
   npm run build   # production-сборка
   npm run lint    # eslint
   npm run preview # просмотр готовой сборки
   ```

### Ollama

Backend ожидает работающий Ollama API (`http://localhost:11434`). После старта контейнера (или локального демона) загрузите нужную модель и убедитесь, что параметр `config.Model` в `internal/adapters/llm/ollama.go` соответствует её имени. Опциональный веб-поиск можно отключить, установив `EnableWebSearch: false` в конфиге клиента.

## API (черновик)

| Метод | Путь | Описание | Auth |
|-------|------|----------|------|
| GET   | `/health` | Проверка состояния backend | нет |
| POST  | `/login` | Вход по email/паролю, возвращает токен-заглушку | нет |
| GET   | `/chats` | Список чатов пользователя | да |
| POST  | `/chats` | Создание чата | да |
| GET   | `/chats/{chat_id}/messages` | История сообщений | да |
| POST  | `/chats/{chat_id}/messages` | Отправка запроса и получение ответа LLM | да |
| POST  | `/documents` | Загрузка документа (multipart/form-data) | да |
| GET   | `/scenarios` | Предустановленные сценарии (contract_helper, marketing) | да |
| GET   | `/config/limits` | Возвращает активные лимиты промптов и файлов | да |

Маршруты под `/chats`, `/documents`, `/scenarios`, `/config`, `/rag` защищены middleware `AuthMiddleware`, который сейчас подставляет фиксированный `user_id`. В production необходимо заменить на реальную проверку JWT.

## Переменные окружения

| Имя | Назначение | Значение по умолчанию |
|-----|------------|-----------------------|
| `POSTGRES_DSN` | Подключение к PostgreSQL (используется в `postgres.NewPool`) | обязательна |
| `HTTP_ADDR`    | Полный адрес HTTP-сервера backend | `:8080` |
| `PORT`         | Альтернативный способ задать порт (Heroku-style) | пусто |

Конфигурация LLM (`BaseURL`, `Model`, `Temperature`, лимиты токенов, флаг веб-поиска) задаётся через структуру `llm.Config` при создании `OllamaClient`. Вынесите значения в `.env`/конфиг при интеграции.

## Известные ограничения

- `AuthMiddleware` и генерация токена в `handlers/auth.go` — демо-заглушки, их нужно заменить на полноценный JWT-пайплайн.
- Модули документов (`DocumentRepo`, RAG) возвращают статические данные и ждут реализации загрузки в постоянное хранилище.
- UI использует моковые данные (`mockChats`, `mockMessages`); интеграция с API отсутствует.
- Dockerfile backend и реальный путь до `cmd/main` расходятся — обновите target или добавьте `cmd/server`.

## Рекомендованные next steps

1. Починить `apps/backend/Dockerfile` и подключить router/LLM wiring в `cmd/main`.
2. Добавить реализацию `DocumentRepo` и интеграцию с реальным сториджем (S3/MinIO) + RAG.
3. Вынести конфигурацию Ollama в `.env`, добавить health-probes и ретраи.
4. Соединить frontend с backend API, внедрить react-query/fetcher, удалить моковые данные.
5. Заменить простые токены на JWT + refresh, покрыть критичные use-case тестами.
