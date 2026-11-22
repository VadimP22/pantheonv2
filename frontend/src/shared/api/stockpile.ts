// Модуль fetch-обёрток для HTTP API складских позиций (запасов).
import type {
  CreateStockpileInputDTO,
  StockpileDTO,
  StockpilePageDTO,
  UpdateStockpileInputDTO,
} from '../generated/dto.ts'
import { request } from './http.ts'

/** Максимально допустимый размер страницы (совпадает с бэкенд-ограничением). */
const maxPageSize = 100

/** Список параметров пагинации списка stockpiles. */
export interface ListStockpilesParams {
  /** Номер страницы (1-based). */
  page?: number
  /** Размер страницы (1..100). */
  pageSize?: number
}

/**
 * Возвращает страницу складских позиций.
 * GET /api/stockpiles
 */
export async function listStockpiles(
  params: ListStockpilesParams = {},
): Promise<StockpilePageDTO> {
  const { page = 1, pageSize = 6 } = params
  const qs = new URLSearchParams({
    page: String(page),
    page_size: String(Math.min(Math.max(pageSize, 1), maxPageSize)),
  })
  return request<StockpilePageDTO>(`/api/stockpiles?${qs.toString()}`)
}

/**
 * Возвращает складскую позицию по идентификатору.
 * GET /api/stockpiles/{id}
 */
export async function getStockpile(id: number): Promise<StockpileDTO> {
  return request<StockpileDTO>(`/api/stockpiles/${id}`)
}

/**
 * Создаёт новую складскую позицию.
 * POST /api/stockpiles
 */
export async function createStockpile(
  input: CreateStockpileInputDTO,
): Promise<StockpileDTO> {
  return request<StockpileDTO>('/api/stockpiles', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
}

/**
 * Частично обновляет складскую позицию по идентификатору.
 * Отсутствующие поля в input не изменяются.
 * PUT /api/stockpiles/{id}
 */
export async function updateStockpile(
  id: number,
  input: UpdateStockpileInputDTO,
): Promise<StockpileDTO> {
  return request<StockpileDTO>(`/api/stockpiles/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
}

/**
 * Удаляет складскую позицию по идентификатору. Идемпотентно:
 * отсутствие записи не является ошибкой.
 * DELETE /api/stockpiles/{id}
 */
export async function deleteStockpile(id: number): Promise<void> {
  await request<void>(`/api/stockpiles/${id}`, { method: 'DELETE' })
}