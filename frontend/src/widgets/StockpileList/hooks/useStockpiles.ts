import { useQuery } from '@tanstack/react-query'
import type { ListStockpilesParams } from '../../../shared/api/stockpile.ts'
import { listStockpiles } from '../../../shared/api/stockpile.ts'
import { stockpileKeys } from '../../../shared/api/keys.ts'

/** Хук подгрузки страницы складов с серверной пагинацией. */
export function useStockpiles(params: ListStockpilesParams) {
  return useQuery({
    queryKey: stockpileKeys.list(params),
    queryFn: () => listStockpiles(params),
  })
}