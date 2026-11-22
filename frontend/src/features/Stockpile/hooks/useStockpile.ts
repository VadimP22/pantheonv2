import { useQuery } from '@tanstack/react-query'
import { getStockpile } from '../../../shared/api/stockpile.ts'
import { stockpileKeys } from '../../../shared/api/keys.ts'

/** Хук подгрузки одного склада по идентификатору. */
export function useStockpile(id: number) {
  return useQuery({
    queryKey: stockpileKeys.detail(id),
    queryFn: () => getStockpile(id),
    enabled: Number.isFinite(id) && id > 0,
  })
}