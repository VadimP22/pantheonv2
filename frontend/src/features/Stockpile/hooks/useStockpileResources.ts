import { useQuery } from '@tanstack/react-query'
import { listStockpileResources } from '../../../shared/api/stockpileResources.ts'
import { stockpileResourceKeys } from '../../../shared/api/keys.ts'

/** Хук подгрузки ресурсов склада по идентификатору. */
export function useStockpileResources(stockpileId: number) {
  return useQuery({
    queryKey: stockpileResourceKeys.list(stockpileId),
    queryFn: () => listStockpileResources(stockpileId),
    enabled: Number.isFinite(stockpileId) && stockpileId > 0,
  })
}