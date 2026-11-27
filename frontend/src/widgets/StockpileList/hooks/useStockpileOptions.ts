import { useQuery } from '@tanstack/react-query'
import { listStockpiles } from '../../../shared/api/stockpile.ts'
import { stockpileKeys } from '../../../shared/api/keys.ts'

/** Хук загрузки списка складов для выпадающего выбора (первая страница, pageSize=100). */
export function useStockpileOptions() {
  return useQuery({
    queryKey: stockpileKeys.list({ page: 1, pageSize: 100 }),
    queryFn: () => listStockpiles({ page: 1, pageSize: 100 }),
  })
}