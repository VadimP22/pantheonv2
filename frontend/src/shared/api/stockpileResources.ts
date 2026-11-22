// Модуль fetch-обёрток для HTTP API ресурсов склада.
import type {
  ChangeStockpileResourceInputDTO,
  StockpileResourceDTO,
} from '../generated/dto.ts'
import { request } from './http.ts'

/**
 * Возвращает ресурсы склада, отсортированные по убыванию количества.
 * GET /api/stockpiles/{id}/resources
 */
export async function listStockpileResources(
  id: number,
): Promise<StockpileResourceDTO[]> {
  return request<StockpileResourceDTO[]>(`/api/stockpiles/${id}/resources`)
}

/**
 * Пополняет ресурс на складе. Возвращает актуальное состояние записи.
 * POST /api/stockpiles/{id}/resources/deposit
 */
export async function depositStockpileResource(
  id: number,
  input: ChangeStockpileResourceInputDTO,
): Promise<StockpileResourceDTO> {
  return request<StockpileResourceDTO>(`/api/stockpiles/${id}/resources/deposit`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
}

/**
 * Списывает ресурс со склада. Успех = 204 без тела.
 * POST /api/stockpiles/{id}/resources/withdraw
 */
export async function withdrawStockpileResource(
  id: number,
  input: ChangeStockpileResourceInputDTO,
): Promise<void> {
  await request<void>(`/api/stockpiles/${id}/resources/withdraw`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
}