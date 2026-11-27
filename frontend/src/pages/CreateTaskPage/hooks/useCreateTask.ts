import { useMutation, useQueryClient } from '@tanstack/react-query'
import type { CreateTaskInputDTO, TaskDTO } from '../../../shared/generated/dto.ts'
import { createTask } from '../../../shared/api/tasks.ts'
import { taskKeys } from '../../../shared/api/keys.ts'

/** Мутация создания новой задачи с инвалидацией списка после успеха. */
export function useCreateTask() {
  const queryClient = useQueryClient()
  return useMutation<TaskDTO, Error, CreateTaskInputDTO>({
    mutationFn: (input: CreateTaskInputDTO) => createTask(input),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: taskKeys.all })
    },
  })
}