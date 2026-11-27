import { useQuery } from '@tanstack/react-query'
import type { ListTasksParams } from '../../../shared/api/tasks.ts'
import { listTasks } from '../../../shared/api/tasks.ts'
import { taskKeys } from '../../../shared/api/keys.ts'

/** Хук подгрузки страницы задач с серверной пагинацией. */
export function useTasks(params: ListTasksParams) {
  return useQuery({
    queryKey: taskKeys.list(params),
    queryFn: () => listTasks(params),
  })
}
