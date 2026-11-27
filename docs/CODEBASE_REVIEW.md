# Codebase Review — оставшиеся находки

Дата: 2026-11-28.
Область: полное дерево `pantheonv2` (backend/ + frontend/ + корень) + состояние git:
ветка `master`, 3 коммита (базовые `498828e feat: add stockpile inventory management app`,
`303ae7d feat: add light/dark theme switcher`; WIP-фича задач и локализаций предметов
закоммичена как `0d5367b feat: add tasks and item localizations; stockpile updated_at and
aggregated counts`). Рабочее дерево содержит только правки этого PR: 5 изменённых
backend-файлов (юниты A и B) + неотслеживаемый `docs/CODEBASE_REVIEW.md` (юнит C);
других незакоммиченных изменений нет. Правки PR не закоммичены и не запушены.

Этот PR исправляет:
- (A) слоистость сервиса tasks: `backend/model/tasks/tasks.go` + `backend/services/tasks/tasks.go` + `backend/controller/tasks/tasks.go` (новый `model.CreateTaskInput`, сервисный Create принимает model, а не dto, контроллер получает конвертер `toCreateModel`; `dto/tasks.go` не меняется);
- (B) `repository/stockpile_resources` `GetByStockpileAndItem` + связывание в `services/stockpile_resources`;
- (C) этот документ.

## Исправлено в этом PR
1. Слоистость tasks: модель/сервис/контроллер разнесены, конвертер Model↔DTO в контроллере — `backend/model/tasks/tasks.go`, `backend/services/tasks/tasks.go`, `backend/controller/tasks/tasks.go`.
2. Конвертеры/DTO-граница на слое tasks (юнит A) — те же три файла.
7. `GetByStockpileAndItem` в репозитории stockpile_resources + связывание в сервисе — `backend/repository/stockpile_resources/*`, `backend/services/stockpile_resources/*`.

## Отложено решением автора (не трогаем сейчас)
3. Дрейф документации в комментариях модели — «оставь как есть».
4. Транзакции вокруг двухзаписного Deposit — «так и задумано».
5. withdraw никогда не уменьшает `tasks.item_count` — «так и задумано».
6. Сбой Touch после успешной мутации возвращает 500 при изменившемся состоянии — «оставь как есть».
8. Мёртвые `List`/`GetByID` и N+1-оболочка — MVP-first, «сначала MVP, потом рефакторинг».

## Открыто (следующие PRs)

### Frontend: пропуски инвалидации кэша
- `frontend/src/features/Stockpile/hooks/useDepositStockpileResource.ts` инвалидирует только `stockpileResourceKeys.list(stockpileId)`, но не `stockpileKeys` list/detail → устаревшие счётчики «предметов/типов» на `StockpilesPage` и устаревший заголовок на `StockpileOverviewPage`.
- `frontend/src/features/Task/hooks/useDepositTask.ts` инвалидирует только `taskKeys.all`, но не кэши stockpile-ресурсов/списка stockpile → Inventory на `StockpileOverviewPage` устаревает после депозита задачи.

### Мёртвый UI-код
- `frontend/src/features/Stockpile/hooks/useWithdrawStockpileResource.ts` — хук определён, нигде не импортируется.
- `frontend/src/shared/api/tasks.ts:62` `deleteTask` и `frontend/src/shared/api/stockpile.ts:79` `deleteStockpile` — обёртки без вызывающих.

### Пробелы фич
- `backend/controller/tasks/tasks.go` не имеет маршрута `DELETE /api/tasks/{id}` вообще (у stockpiles-контроллера DELETE есть); кнопок удаления в UI нет.
- `item_localizations` наполовину построена: таблица + миграция `000003` + эндпоинт `GET /api/loc/{lang}/{item}` + repo + service, но нет writer-эндпоинта, нет seed-строк и нет frontend-потребителя `/api/loc` нигде в `frontend/src` — на свежей БД эндпоинт всегда отвечает 404.

### Робастность
- `frontend/src/App.tsx` нет catch-all 404-маршрута (несовпавшие пути рендерят пустой main; `/stockpiles/abc` → NaN id → пустая сетка Inventory).
- `frontend/src/pages/CreateTaskPage/ui/CreateTaskPage.tsx:64` шлёт `item_count` 0 (`Number("")`=0) и target 0 без валидации.
- `widgets/StockpileList/hooks/useStockpileOptions.ts` `pageSize=100` молча обрезает список после 100 stockpile (бэкенд `max_page_size=100` в `backend/lib/response.go:17`).
- `frontend/src/features/theme/store/themeStore.ts:24-26` читает `window.localStorage` при импорте модуля (SSR-хрупкость), тогда как AGENTS.md обещает тёмную тему через `prefers-color-scheme`, а реализация использует `[data-theme='dark']` + ручной переключатель (расхождение doc/impl в `frontend/src/index.css:62`).
- `eslint.config.js` не имеет type-aware rules / react-плагина против рекомендации `frontend/README.md`.
- Нет ни одного `_test.go` и ни одного `*.test.ts` во всём дереве.
- Корневой `.gitignore` не перечисляет `backend.db` явно (покрывает только `backend/.gitignore`); AGENTS.md это признаёт.

### Мёртвый компонент
- `frontend/src/entities/InventoryComponents/ui/InventorySlot.tsx` никогда не импортируется (`Inventory.tsx:7` импортирует только `InventorySlotWithItem`, пустые слоты рендерят `InventorySlotWithItem` без `itemId` в `Inventory.tsx:73-75`).

## По дизайну (не баг, не переоткрывать)
- `ParseID` отклоняет `id<1` → 400.
- Обработка 204 в `shared/api/http.ts:26-28`.
- `max_page_size=100`.
- Агрегаты списка stockpile — коррелированные подзапросы, пересчитываемые на строку (нормально на масштабе MVP).
- gofmt-clean по построению.

## Приоритеты
- P0 (видимость): frontend-инвалидация кэша (`useDepositStockpileResource.ts`, `useDepositTask.ts`).
- P1: наполовину построенная фича item_localizations, пробел `DELETE /api/tasks/{id}`, кнопки удаления в UI.
- P2: слоистость (исправлено в этом PR), нарушение PascalCase у `entities/SidebarComponents/ui/sidebarButton.tsx`, валидация счётчика задач, 404-маршрут.
- P3: N+1 (отложенный пункт 8), SSR-риск, конфиг линтера, тесты, git-гигиена (`.gitignore`/`backend.db`).

Пункты 1, 2, 7 помечены «исправлено в этом PR».

## Верификация
- `Get-Content docs/CODEBASE_REVIEW.md` — файл непустой.
- `git status --porcelain` — ровно 5 изменённых backend-файлов юнитов A/B (`backend/controller/tasks/tasks.go`,
  `backend/model/tasks/tasks.go`, `backend/services/tasks/tasks.go`, `backend/repository/stockpile_resources/stockpile_resources.go`,
  `backend/services/stockpile_resources/stockpile_resources.go`) + неотслеживаемый `docs/CODEBASE_REVIEW.md`; других изменений нет.
- Файл содержит литералы: `Исправлено в этом PR`, `Отложено решением`, `Открыто`, `GetByStockpileAndItem`.
- Сборка: из `backend/` — `go build ./...` exit 0, `go vet ./...` exit 0,
  `gofmt -l` на пяти файлах PR пуст.
