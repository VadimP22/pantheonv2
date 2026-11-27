import { useQuery } from '@tanstack/react-query'
import { getTask } from '../../../shared/api/tasks.ts'
import { taskKeys } from '../../../shared/api/keys.ts'

/** Хук подгрузки одной задачи по идентификатору. */
export function useTask(id: number) {
  return useQuery({
    queryKey: taskKeys.detail(id),
    queryFn: () => getTask(id),
    enabled: Number.isFinite(id) && id > 0,
  })
}
