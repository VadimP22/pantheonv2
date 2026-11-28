import { useQuery } from '@tanstack/react-query'
import { listStockpiles } from '../../../shared/api/stockpile.ts'
import { stockpileKeys } from '../../../shared/api/keys.ts'

/** Хук: общее количество складов (для счётчика на карточке обзора). */
export function useStockpilesCount() {
  return useQuery({
    queryKey: stockpileKeys.list({ page: 1, pageSize: 1 }),
    queryFn: () => listStockpiles({ page: 1, pageSize: 1 }),
    select: (data) => data.total,
  })
}
