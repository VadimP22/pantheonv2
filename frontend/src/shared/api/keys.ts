import type { ListStockpilesParams } from './stockpile.ts'

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