import { useState } from 'react'
import type { FormEvent } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import PageHeader from '../../../shared/PageHeader/ui/PageHeader.tsx'
import AppCard from '../../../shared/AppCard/ui/AppCard.tsx'
import AppText from '../../../shared/AppText/ui/AppText.tsx'
import AppTextInput from '../../../shared/AppTextInput/ui/AppTextInput.tsx'
import AppNumericInput from '../../../shared/AppNumericInput/ui/AppNumericInput.tsx'
import AppButton from '../../../shared/AppButton/ui/AppButton.tsx'
import InventorySlotWithItem from '../../../entities/InventoryComponents/ui/InventorySlotWithItem.tsx'
import { useStockpile } from '../../../features/Stockpile/hooks/useStockpile.ts'
import { useDepositStockpileResource } from '../../../features/Stockpile/hooks/useDepositStockpileResource.ts'
import styles from './DepositResourcePage.module.css'

/** Идентификатор формы, чтобы кнопка «Внести» из слота хэдера могла её сабмитить. */
const FORM_ID = 'deposit-resource-form'

export function DepositResourcePage() {
  const navigate = useNavigate()
  const params = useParams()
  const id = Number(params.id)
  const { data } = useStockpile(id)
  const { mutate, isPending, isError, error } = useDepositStockpileResource(id)
  const [itemId, setItemId] = useState('')
  const [countText, setCountText] = useState('1')

  const isValidId = Number.isFinite(id) && id > 0
  const title = data?.name ?? 'Внесение ресурса'
  const description = !isValidId
    ? 'Некорректный идентификатор склада.'
    : data
      ? `Внести ресурс на склад «${data.name}»`
      : 'Загрузка данных склада…'

  const count = Math.max(1, Math.floor(Number(countText) || 0))
  const isValidForm = !!itemId && count > 0

  const handleSubmit = (event: FormEvent) => {
    event.preventDefault()
    mutate(
      { item_id: itemId, count },
      { onSuccess: () => navigate(`/stockpiles/${id}`) },
    )
  }

  return (
    <section>
      <PageHeader
        breadcrumbs="Главная / Склады / Обзор / Внесение ресурса"
        title={title}
        description={description}
        actions={
          <>
            <AppButton primary type="submit" form={FORM_ID} disabled={isPending || !isValidForm}>
              {'Внести'}
            </AppButton>
            <AppButton onClick={() => navigate(`/stockpiles/${id}`)}>{'Отмена'}</AppButton>
          </>
        }
      />

      <AppCard secondary>
        <form id={FORM_ID} className={styles.form} onSubmit={handleSubmit}>
          <div className={styles.preview}>
            <InventorySlotWithItem itemId={itemId || undefined} count={count} selected={!!itemId} />
          </div>
          <div className={styles.field}>
            <AppText className={styles.label}>{'Идентификатор предмета'}</AppText>
            <AppTextInput
              value={itemId}
              placeholder="diamond"
              onChange={(e) => setItemId(e.target.value)}
            />
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

export default DepositResourcePage