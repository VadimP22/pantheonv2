import { useQuery } from '@tanstack/react-query'
import { getStockpile } from '../../../shared/api/stockpile.ts'
import { stockpileKeys } from '../../../shared/api/keys.ts'

/**
 * Хук подгрузки склада, к которому привязана задача.
 * Возвращает null-данные, если склад не задан (id = 0).
 */
export function useTaskStockpile(stockpileId: number | undefined) {
  const parsed = stockpileId ?? 0
  return useQuery({
    queryKey: stockpileKeys.detail(parsed),
    queryFn: () => getStockpile(parsed),
    enabled: Number.isFinite(parsed) && parsed > 0,
  })
}

export default useTaskStockpile
