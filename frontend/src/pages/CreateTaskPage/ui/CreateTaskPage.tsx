import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import type { FormEvent } from 'react'
import PageHeader from '../../../shared/PageHeader/ui/PageHeader.tsx'
import AppCard from '../../../shared/AppCard/ui/AppCard.tsx'
import AppText from '../../../shared/AppText/ui/AppText.tsx'
import AppTextInput from '../../../shared/AppTextInput/ui/AppTextInput.tsx'
import AppNumericInput from '../../../shared/AppNumericInput/ui/AppNumericInput.tsx'
import AppDropDown from '../../../shared/AppDropDown/ui/AppDropDown.tsx'
import AppHorizontalLine from '../../../shared/AppHorizontalLine/ui/AppHorizontalLine.tsx'
import AppButton from '../../../shared/AppButton/ui/AppButton.tsx'
import { useCreateTask } from '../hooks/useCreateTask.ts'
import { useStockpileOptions } from '../../../widgets/StockpileList/hooks/useStockpileOptions.ts'
import type { StockpileDTO } from '../../../shared/generated/dto.ts'
import styles from './CreateTaskPage.module.css'

/** Идентификатор формы, чтобы кнопка «Создать» из слота хэдера могла её сабмитить. */
const FORM_ID = 'create-task-form'

/** Начальное состояние полей формы. */
const initialFields = {
  name: '',
  description: '',
  iconItemName: '',
  isCheckbox: 'false',
  itemId: '',
  itemCount: '',
  targetItemCount: '',
  stockpileId: '',
}
type Fields = typeof initialFields

/** Страница создания новой задачи. */
export function CreateTaskPage() {
  const navigate = useNavigate()
  const { mutate, isPending, isError, error } = useCreateTask()
  const { data: stockpileData, isPending: stockpilesLoading } = useStockpileOptions()
  const stockpiles: StockpileDTO[] = stockpileData?.items ?? []
  const [fields, setFields] = useState<Fields>(initialFields)

  const setField = (key: keyof Fields) => (value: string) =>
    setFields((prev) => ({ ...prev, [key]: value }))

  const handleSubmit = (event: FormEvent) => {
    event.preventDefault()

    const isCheckbox = fields.isCheckbox === 'true'
    if (isCheckbox && fields.iconItemName.trim() === '') {
      window.alert('Укажите иконку (item name)')
      return
    }
    if (!isCheckbox && fields.itemId.trim() === '') {
      window.alert('Укажите ID предмета')
      return
    }

    mutate(
      {
        name: fields.name,
        description: fields.description,
        icon_item_name: fields.iconItemName,
        is_checkbox: isCheckbox,
        item_id: isCheckbox ? '' : fields.itemId,
        item_count: isCheckbox ? 0 : Number(fields.itemCount),
        target_item_count: isCheckbox ? 0 : Number(fields.targetItemCount),
        stockpile_id: fields.stockpileId === '' ? undefined : Number(fields.stockpileId),
      },
      { onSuccess: () => navigate('/tasks') },
    )
  }

  return (
    <section>
      <PageHeader
        breadcrumbs="Главная / Задачи / Создать"
        title="Создать задачу"
        description="Заполните параметры новой задачи."
        actions={
          <>
            <AppButton primary type="submit" form={FORM_ID} disabled={isPending}>
              {'Создать'}
            </AppButton>
            <AppButton onClick={() => navigate('/tasks')}>{'Отмена'}</AppButton>
          </>
        }
      />

      <AppCard secondary>
        <form id={FORM_ID} className={styles.form} onSubmit={handleSubmit}>
          <div className={styles.field}>
            <AppText className={styles.label}>{('Тип задачи')}</AppText>
            <AppDropDown
              value={fields.isCheckbox}
              onChange={(e) => setField('isCheckbox')(e.target.value)}
            >
              <option value="false">{'ресурсы'}</option>
              <option value="true">{'чекбокс'}</option>
            </AppDropDown>
          </div>

          <AppHorizontalLine />

          <div className={styles.field}>
            <AppText className={styles.label}>{'Название задачи'}</AppText>
            <AppTextInput
              value={fields.name}
              placeholder="Название"
              onChange={(e) => setField('name')(e.target.value)}
            />
          </div>

          <div className={styles.field}>
            <AppText className={styles.label}>{'Описание'}</AppText>
            <AppTextInput
              value={fields.description}
              placeholder="Описание"
              onChange={(e) => setField('description')(e.target.value)}
            />
          </div>

          {fields.isCheckbox === 'true' && (
          <div className={styles.field}>
            <AppText className={styles.label}>{'Иконка (item name)'}</AppText>
            <AppTextInput
              value={fields.iconItemName}
              placeholder="minecraft:barrier"
              onChange={(e) => setField('iconItemName')(e.target.value)}
            />
          </div>
        )}

        {fields.isCheckbox !== 'true' && (
          <div className={styles.field}>
            <AppText className={styles.label}>{'ID предмета'}</AppText>
            <AppTextInput
              value={fields.itemId}
              placeholder="minecraft:iron_ingot"
              onChange={(e) => setField('itemId')(e.target.value)}
            />
          </div>
        )}

        {fields.isCheckbox !== 'true' && (
          <>
            <div className={styles.field}>
              <AppText className={styles.label}>{('Текущее количество')}</AppText>
              <AppNumericInput
                value={fields.itemCount}
                onChange={(e) => setField('itemCount')(e.target.value)}
              />
            </div>

            <div className={styles.field}>
              <AppText className={styles.label}>{('Целевое количество')}</AppText>
              <AppNumericInput
                value={fields.targetItemCount}
                onChange={(e) => setField('targetItemCount')(e.target.value)}
              />
            </div>
          </>
        )}

          <div className={styles.field}>
            <AppText className={styles.label}>{fields.isCheckbox === 'true' ? 'Целевой склад (необязательно)' : 'Целевой склад'}</AppText>
            <AppDropDown
              placeholder={stockpilesLoading ? 'Загрузка складов…' : 'Выберите склад…'}
              value={fields.stockpileId}
              onChange={(e) => setField('stockpileId')(e.target.value)}
              disabled={stockpilesLoading}
            >
              <option value="">{'—'}</option>
              {stockpiles.map((s) => (
                <option key={s.id} value={String(s.id)}>{s.name}</option>
              ))}
            </AppDropDown>
          </div>

          {isError && <AppText className={styles.error}>{`Ошибка: ${error?.message}`}</AppText>}
        </form>
      </AppCard>
    </section>
  )
}

export default CreateTaskPage