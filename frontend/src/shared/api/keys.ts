import type { ListStockpilesParams } from './stockpile.ts'
import type { ListTasksParams } from './tasks.ts'

/** Ключи TanStack Query для запросов к складам. */
export const stockpileKeys = {
  all: ['stockpiles'] as const,
  list: (params: ListStockpilesParams) => ['stockpiles', 'list', params] as const,
  detail: (id: number) => ['stockpiles', 'detail', id] as const,
}

/** Ключи TanStack Query для запросов ресурсов складов. */
export const stockpileResourceKeys = {
  all: ['stockpileResources'] as const,
  list: (stockpileId: number) => ['stockpileResources', 'list', stockpileId] as const,
}

/** Ключи TanStack Query для запросов задач. */
export const taskKeys = {
  all: ['tasks'] as const,
  list: (params: ListTasksParams) => ['tasks', 'list', params] as const,
  detail: (id: number) => ['tasks', 'detail', id] as const,
  delete: (id: number) => ['tasks', 'delete', id] as const,
}

/** Ключи TanStack Query для запросов локализаций предметов. */
export const itemLocalizationKeys = {
  all: ['itemLocalizations'] as const,
  detail: (lang: string, itemId: string) => ['itemLocalizations', 'detail', lang, itemId] as const,
}