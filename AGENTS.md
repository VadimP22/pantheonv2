# AGENTS.md — pantheonv2

Guidelines for AI agents (and humans) working in this repository.

В этом файле зафиксированы конвенции проекта. Часть из них уже реализована в коде
(«существующие»), часть — целевые соглашения, которых нужно придерживаться при новой
работе (отмечены как «целевые»). Комментарии, сообщения и документация на русском.

## Обзор проекта

`pantheonv2` — full-stack приложение «запасы/склады» (складские позиции/топология):

- `backend/` — Go-сервис (SQLite + HTTP API);
- `frontend/` — React SPA (Vite + TypeScript).

Репозиторий свежий: ветка `master`, коммитов и remote пока нет.

## Стек проекта

**Backend:**

- Go 1.25.7 (module `backend`).
- SQLite через `modernc.org/sqlite` (чистый Go, без CGO).
- Миграции — `golang-migrate/migrate/v4`.

**Frontend:**

- React 19, Vite 8, TypeScript ~6, react-router-dom 7.
- Иконки — `@imumd/mcicons` (Minecraft-стиль).
- `tygo` — генерация TypeScript-типов из Go-DTO (именно поэтому пакет `dto` держится
  плоским). После изменений в `backend/dto/` обязательно выполнять
  `npm run generate`.
- TanStack Query (серверное/данные-состояние) и Zustand (клиентское/локальное
  состояние). **Целевые:** зависимости пока не установлены, но конвенции под них
  закреплены заранее (см. «Frontend конвенции»).

## Структура репозитория

```
backend/   — Go-сервис (SQLite + HTTP)
frontend/  — React SPA (Vite + TypeScript)
.gitignore — только артефакты инструментов (.codesearch.db, .serena)
```

Замечание: корневой `.gitignore` не покрывает `backend/backend.db` (файл БД) —
новые артефакты БД/сборки стоит игнорировать.

## Backend конвенции (Go)

### Архитектура и связывание

- Код разбит по фичам, слои внутри фичи: `repository/<feature>` →
  `services/<feature>` → `controller/<feature>`; рядом плоские `model/<feature>`
  и `dto`.
- Связывание вручную в `main.go`: `New(db)` → `New(repo)` → `New(svc)` →
  `Register(mux)`.
- Направление зависимостей: репозиторий импортирует модель; сервис — репозиторий и
  модель; контроллер — сервис, модель и `dto`.

### Миграции

- Файлы в `backend/migrations/` с именованием `NNNNNN_name.up.sql` / `.down.sql`.
- Миграции применяются автоматически при старте (`m.Up()`).

### Конфигурация среды

- `DATABASE_PATH` — путь к БД (дефолт `backend.db`);
- `MIGRATIONS_DIR` — каталог миграций (дефолт `./migrations`);
- `HTTP_ADDR` — адрес сервера (дефолт `:8080`).

### Модель и DTO

- Структуры модели **без JSON-тегов**: конверсия Model↔DTO происходит только в
  контроллере.
- `dto` — единый плоский пакет, один файл на ресурс (ради `tygo`). Типы называются
  `<Name>DTO`; `json`-теги в `snake_case` (`icon_name`, `page_size`).

### HTTP

- Стандартный `net/http.ServeMux`; маршруты — паттерны с методом
  (`GET /stockpiles`, `POST /stockpiles`, `PUT /stockpiles/{id}` …).
- Все HTTP-маршруты начинаются со служебного префикса `/api`:
  `GET /api/stockpiles`, `POST /api/stockpiles`, `PUT /api/stockpiles/{id}` ….
- Path-переменные — через `r.PathValue("id")`.
- Ошибки сервиса сопоставляются с HTTP-статусами в `writeServiceError`
  (`ErrNotFound` → 404, `ErrInvalidInput` → 400).

### Стиль кода и обработка ошибок

- Все комментарии (doc-комментарии пакетов, типов, методов) на русском.
- Ошибки оборачиваются `fmt.Errorf("...: %w")` (lowercase-сообщение).
- Sentinel-ошибки: `var ErrX = errors.New("stockpiles: ...")` — на сервисном слое.
- Имена колонок выносятся в константы; SELECT собирается хелпером `buildSelect`.
- Частичное обновление — поля-указатели `*T` (nil = «не изменять»);
  `ErrNoUpdateFields`, когда нечего обновлять.
- Пагинация по умолчанию: `page=1`, `page_size=6`, `max_page_size=100`.
- Изменения в БД идут через raw SQL-запросы с параметрами (`?`).

## Frontend конвенции (React/TypeScript/FSD)

### Структура (Feature-Sliced Design)

Код под `src/` разбит по слоям FSD: `entities/`, `pages/`, `shared/`, `widgets/`.
Внутри слоя/фичи — подпапки:

- **UI** → `<Name>/ui` (см. правило именования файлов);
- **Сторы (Zustand)** → `<Name>/store` *(целевое)*;
- **Хуки (обёртки над TanStack Query и пр.)** → `<Name>/hooks` *(целевое)*;
- **Прочее (утилиты и т.п.)** → `<Name>/lib` *(целевое)*.

### Именование файлов

- **UI-компоненты** (`*.tsx` и их CSS-модули) → **`ThisNaming.tsx`** (PascalCase):
  `AppCard.tsx`, `AppCard.module.css`.
- **Не-UI код** (утилиты, хуки, сторы, api-обёртки и т.п.) → **`thisNaming.ts`**
  (camelCase): `keys.ts`, `useStockpiles.ts`.

### Компоненты

- Компонент объявляется именованным `export function X()` и дополнительно
  экспортируется `export default X`.
- Интерфейс пропсов — `XProps` с JSDoc-комментариями на русском.
- Shared-примитивы идут с префиксом `App*`: `AppCard`, `AppLayout`, `AppText`,
  `AppTitle`, `AppProgressBar`.
- CSS-модули лежат рядом (`<Name>.module.css`), классы — camelCase
  (`root`, `topEdge`, `mainPane`).
- Глобальные токены — CSS-переменные в `src/index.css` (`--bg`, `--text`,
  `--border`, `--accent`) с тёмной темой через `prefers-color-scheme`.
- Типы импортируются через `import type` (`verbatimModuleSyntax`); импорты
  локальных `.tsx`-модулей с явным расширением.

### Работа с API *(целевое)*

- Обёртки `fetch` для запросов к бэкенду — в **`src/shared/api/<feature>`**
  (директория/файл на фичу).
- Ключи TanStack Query — централизованно в **`src/shared/api/keys.ts`**.

## Команды

**Backend** (`backend/`):

```sh
go run .          # миграции + HTTP-сервер
go build ./...
go vet ./...
go test ./...
go fmt ./...
```

**Frontend** (`frontend/`):

```sh
npm run dev
npm run build     # tsc -b && vite build
npm run generate  # tygo generate — перегенерировать TS-типы из Go-DTO
npm run lint
```

Пакетный менеджер — **npm**. Команды пакета вызываются через `npm run <script>`;
для установки зависимостей используется `npm install`. Артефакты pnpm
(`pnpm-lock.yaml` и пр.) не используются и не создаются; единственный lock-файл —
`package-lock.json`.