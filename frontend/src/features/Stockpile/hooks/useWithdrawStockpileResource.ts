import { useMutation, useQueryClient } from '@tanstack/react-query'
import type { ChangeStockpileResourceInputDTO } from '../../../shared/generated/dto.ts'
import { withdrawStockpileResource } from '../../../shared/api/stockpileResources.ts'
import { stockpileResourceKeys } from '../../../shared/api/keys.ts'

/** Мутация списания ресурса со склада; после успеха инвалидирует список ресурсов. */
export function useWithdrawStockpileResource(stockpileId: number) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (input: ChangeStockpileResourceInputDTO) =>
      withdrawStockpileResource(stockpileId, input),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: stockpileResourceKeys.list(stockpileId) })
    },
  })
}