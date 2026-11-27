import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import TaskCard from '../../../entities/TaskComponents/ui/TaskCard.tsx'
import AppPaginator from '../../../shared/AppPaginator/ui/AppPaginator.tsx'
import AppText from '../../../shared/AppText/ui/AppText.tsx'
import { useTasks } from '../hooks/useTasks.ts'
import { useDeleteTask } from '../hooks/useDeleteTask.ts'
import styles from './TasksList.module.css'

/** Количество карточек на странице: сетка 3×2. */
const PAGE_SIZE = 6

export interface TasksListProps {
  /** Внешние классы для корневого контейнера. */
  className?: string
}

/**
 * Виджет списка задач: самодостаточная сетка 3×2 с пагинацией.
 * Данные подгружаются с сервера через TanStack Query.
 */
export function TasksList({ className }: TasksListProps) {
  const [page, setPage] = useState(1)
  const navigate = useNavigate()
  const deleteTaskMutation = useDeleteTask()
  const { data, isPending, isError, error } = useTasks({ page, pageSize: PAGE_SIZE })
  const rootClassName = className ? `${styles.root} ${className}` : styles.root

  if (isPending && !data) {
    return <AppText className={rootClassName}>{'Загрузка задач…'}</AppText>
  }

  if (isError) {
    return (
      <AppText className={rootClassName}>
        {`Не удалось загрузить задачи: ${error instanceof Error ? error.message : String(error)}`}
      </AppText>
    )
  }

  const totalPages = data?.total_pages ?? 1

  return (
    <section className={rootClassName}>
      <div className={styles.grid}>
        {(data?.items ?? []).map((task) => (
          <TaskCard
            key={task.id}
            name={task.name}
            description={task.description ?? ''}
            iconItemName={task.icon_item_name}
            isCheckbox={task.is_checkbox}
            itemCount={task.item_count}
            targetItemCount={task.target_item_count}
            done={task.done}
            onDepositClick={
              !task.is_checkbox && task.stockpile_id != null
                ? () => navigate('/tasks/' + task.id + '/deposit')
                : undefined
            }
            onCompleteClick={
              !task.is_checkbox && task.stockpile_id != null && task.item_count >= task.target_item_count
                ? () => deleteTaskMutation.mutate(task.id)
                : undefined
            }
          />
        ))}
      </div>
      <AppPaginator page={page} totalPages={totalPages} onPageChange={setPage} />
    </section>
  )
}

export default TasksList
