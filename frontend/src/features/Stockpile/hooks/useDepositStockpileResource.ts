import { useMutation, useQueryClient } from '@tanstack/react-query'
import type { ChangeStockpileResourceInputDTO } from '../../../shared/generated/dto.ts'
import { depositStockpileResource } from '../../../shared/api/stockpileResources.ts'
import { stockpileResourceKeys } from '../../../shared/api/keys.ts'

/** Мутация пополнения ресурса на складе; после успеха инвалидирует список ресурсов. */
export function useDepositStockpileResource(stockpileId: number) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (input: ChangeStockpileResourceInputDTO) =>
      depositStockpileResource(stockpileId, input),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: stockpileResourceKeys.list(stockpileId) })
    },
  })
}