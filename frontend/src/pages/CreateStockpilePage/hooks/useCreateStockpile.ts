import { useMutation, useQueryClient } from '@tanstack/react-query'
import type { CreateStockpileInputDTO } from '../../../shared/generated/dto.ts'
import { createStockpile } from '../../../shared/api/stockpile.ts'
import { stockpileKeys } from '../../../shared/api/keys.ts'

/** Мутация создания нового склада с инвалидацией списка после успеха. */
export function useCreateStockpile() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (input: CreateStockpileInputDTO) => createStockpile(input),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: stockpileKeys.all })
    },
  })
}