# pantheonv2

Full-stack приложение «запасы/склады»: складские позиции и их топология. Реализовано
как React SPA (frontend) + Go-сервис с SQLite (backend).

## Стек

**Backend** (`backend/`)

- Go 1.25.7
- SQLite через `modernc.org/sqlite` (чистый Go, без CGO)
- Миграции — `golang-migrate/migrate/v4` (применяются автоматически при старте)

**Frontend** (`frontend/`)

- React 19, Vite, TypeScript
- Feature-Sliced Design (`entities/`, `pages/`, `shared/`, `widgets/`)
- TanStack Query (серверное состояние) + Zustand (клиентское)
- Иконки — `@imumd/mcicons`
- `tygo` — генерация TypeScript-типов из Go-DTO

## Структура

```
backend/   — Go-сервис (SQLite + HTTP API)
frontend/  — React SPA (Vite + TypeScript)
tygo.yaml  — конфигурация генерации TS-типов из Go-DTO
```

## Запуск

**Backend** (из `backend/`):

```sh
go run .
```

Переменные окружения: `DATABASE_PATH` (дефолт `backend.db`),
`MIGRATIONS_DIR` (дефолт `./migrations`), `HTTP_ADDR` (дефолт `:8080`).

**Frontend** (из `frontend/`):

```sh
npm install
npm run dev
```

После изменений в `backend/dto/` перегенерировать TS-типы:

```sh
npm run generate
```

## API

Все маршруты имеют префикс `/api`.

| Метод | Путь | Описание |
| --- | --- | --- |
| `GET` | `/api/stockpiles` | Список складов (пагинация) |
| `POST` | `/api/stockpiles` | Создать склад |
| `GET` | `/api/stockpiles/{id}` | Получить склад |
| `PUT` | `/api/stockpiles/{id}` | Обновить склад |
| `DELETE` | `/api/stockpiles/{id}` | Удалить склад |
| `GET` | `/api/stockpiles/{id}/resources` | Список ресурсов склада |
| `POST` | `/api/stockpiles/{id}/resources/deposit` | Внести ресурс |
| `POST` | `/api/stockpiles/{id}/resources/withdraw` | Списать ресурс |

## Команды

**Backend:**

```sh
go build ./...
go vet ./...
go test ./...
go fmt ./...
```

**Frontend:**

```sh
npm run dev       # dev-сервер
npm run build     # tsc -b && vite build
npm run lint      # eslint
npm run generate  # tygo generate — перегенерировать TS-типы из Go-DTO
```