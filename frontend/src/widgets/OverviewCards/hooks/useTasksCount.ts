import { useQuery } from '@tanstack/react-query'
import { listTasks } from '../../../shared/api/tasks.ts'
import { taskKeys } from '../../../shared/api/keys.ts'

/** Хук: общее количество задач (для счётчика на карточке обзора). */
export function useTasksCount() {
  return useQuery({
    queryKey: taskKeys.list({ page: 1, pageSize: 1 }),
    queryFn: () => listTasks({ page: 1, pageSize: 1 }),
    select: (data) => data.total,
  })
}
