import { useState } from 'react'
import type { FormEvent } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import PageHeader from '../../../shared/PageHeader/ui/PageHeader.tsx'
import AppCard from '../../../shared/AppCard/ui/AppCard.tsx'
import AppText from '../../../shared/AppText/ui/AppText.tsx'
import AppNumericInput from '../../../shared/AppNumericInput/ui/AppNumericInput.tsx'
import AppButton from '../../../shared/AppButton/ui/AppButton.tsx'
import InventorySlotWithItem from '../../../entities/InventoryComponents/ui/InventorySlotWithItem.tsx'
import { useTask } from '../../../features/Task/hooks/useTask.ts'
import { useDepositTask } from '../../../features/Task/hooks/useDepositTask.ts'
import styles from './TaskDepositPage.module.css'

/** Идентификатор формы, чтобы кнопка «Внести» из слота хэдера могла её сабмитить. */
const FORM_ID = 'deposit-task-form'

/** Страница внесения ресурса в задачу: предмет берётся из задачи, вводится только количество. */
export function TaskDepositPage() {
  const navigate = useNavigate()
  const params = useParams()
  const id = Number(params.id)
  const { data } = useTask(id)
  const { mutate, isPending, isError, error } = useDepositTask(id)
  const [countText, setCountText] = useState('1')

  const isValidId = Number.isFinite(id) && id > 0
  const title = data?.name ?? 'Внесение ресурса'
  const description = !isValidId
    ? 'Некорректный идентификатор задачи.'
    : data
      ? `Внести ресурс «${data.item_id}» в задачу «${data.name}» — количество также добавится на склад задачи`
      : 'Загрузка данных задачи…'

  const count = Math.max(1, Math.floor(Number(countText) || 0))

  const handleSubmit = (event: FormEvent) => {
    event.preventDefault()
    mutate(
      { count },
      { onSuccess: () => navigate('/tasks') },
    )
  }

  return (
    <section>
      <PageHeader
        breadcrumbs="Главная / Задачи / Внесение ресурса"
        title={title}
        description={description}
        actions={
          <>
            <AppButton primary type="submit" form={FORM_ID} disabled={isPending || !isValidId}>
              {'Внести'}
            </AppButton>
            <AppButton onClick={() => navigate('/tasks')}>{'Отмена'}</AppButton>
          </>
        }
      />

      <AppCard secondary>
        <form id={FORM_ID} className={styles.form} onSubmit={handleSubmit}>
          <div className={styles.preview}>
            <InventorySlotWithItem itemId={data?.item_id ?? ''} count={count} selected={!!data?.item_id} />
          </div>
          <div className={styles.field}>
            <AppText className={styles.label}>{'Количество'}</AppText>
            <AppNumericInput value={countText} onChange={(e) => setCountText(e.target.value)} />
          </div>
          {isError && <AppText className={styles.error}>{`Ошибка: ${error?.message}`}</AppText>}
        </form>
      </AppCard>
    </section>
  )
}

export default TaskDepositPage
