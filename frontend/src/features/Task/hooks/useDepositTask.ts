import { useMutation, useQueryClient } from '@tanstack/react-query'
import type { DepositTaskInputDTO } from '../../../shared/generated/dto.ts'
import { depositTask } from '../../../shared/api/tasks.ts'
import { taskKeys } from '../../../shared/api/keys.ts'

/** Мутация внесения ресурса в задачу; после успеха инвалидирует список задач. */
export function useDepositTask(id: number) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (input: DepositTaskInputDTO) => depositTask(id, input),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: taskKeys.all })
    },
  })
}
