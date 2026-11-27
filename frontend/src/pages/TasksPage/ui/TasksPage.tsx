import { useNavigate } from 'react-router-dom'
import PageHeader from '../../../shared/PageHeader/ui/PageHeader.tsx'
import AppButton from '../../../shared/AppButton/ui/AppButton.tsx'
import TasksList from '../../../widgets/TasksList/ui/TasksList.tsx'

export function TasksPage() {
  const navigate = useNavigate()
  return (
    <section>
      <PageHeader
        breadcrumbs="Главная / Задачи"
        title="Задачи"
        description="Список задач: чекбоксы и задачи с прогрессом."
        actions={
          <AppButton primary onClick={() => navigate('/tasks/create')}>
            {'Создать задачу'}
          </AppButton>
        }
      />
      <TasksList />
    </section>
  )
}

export default TasksPage
