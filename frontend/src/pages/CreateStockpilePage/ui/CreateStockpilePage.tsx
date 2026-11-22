import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import type { FormEvent } from 'react'
import PageHeader from '../../../shared/PageHeader/ui/PageHeader.tsx'
import AppCard from '../../../shared/AppCard/ui/AppCard.tsx'
import AppText from '../../../shared/AppText/ui/AppText.tsx'
import AppTextInput from '../../../shared/AppTextInput/ui/AppTextInput.tsx'
import AppNumericInput from '../../../shared/AppNumericInput/ui/AppNumericInput.tsx'
import AppButton from '../../../shared/AppButton/ui/AppButton.tsx'
import { useCreateStockpile } from '../hooks/useCreateStockpile.ts'
import styles from './CreateStockpilePage.module.css'

/** Идентификатор формы, чтобы кнопка «Создать» из слота хэдера могла её сабмитить. */
const FORM_ID = 'create-stockpile-form'

/** Начальное состояние полей формы (координаты — строки из числового ввода). */
const initialFields = { name: '', dimension: 'overworld', x: '', y: '', z: '' }
type Fields = typeof initialFields

export function CreateStockpilePage() {
  const navigate = useNavigate()
  const { mutate, isPending, isError, error } = useCreateStockpile()
  const [fields, setFields] = useState<Fields>(initialFields)

  const setField = (key: keyof Fields) => (value: string) =>
    setFields((prev) => ({ ...prev, [key]: value }))

  const handleSubmit = (event: FormEvent) => {
    event.preventDefault()
    mutate(
      {
        name: fields.name,
        dimension: fields.dimension,
        icon_name: '', // поле не выведено в форму; на бэкенде допускает пустую строку
        x: Number(fields.x),
        y: Number(fields.y),
        z: Number(fields.z),
      },
      { onSuccess: () => navigate('/stockpiles') },
    )
  }

  return (
    <section>
      <PageHeader
        breadcrumbs="Главная / Склады / Создание"
        title="Создание склада"
        description="Заполните параметры нового склада и его топологии хранения."
        actions={
          <>
            <AppButton primary type="submit" form={FORM_ID} disabled={isPending}>
              {'Создать'}
            </AppButton>
            <AppButton onClick={() => navigate('/stockpiles')}>{'Отмена'}</AppButton>
          </>
        }
      />

      <AppCard secondary>
        <form id={FORM_ID} className={styles.form} onSubmit={handleSubmit}>
          <div className={styles.field}>
            <AppText className={styles.label}>{'Название склада'}</AppText>
            <AppTextInput
              value={fields.name}
              placeholder="Мой склад"
              onChange={(e) => setField('name')(e.target.value)}
            />
          </div>

          <div className={styles.field}>
            <AppText className={styles.label}>{'Измерение'}</AppText>
            <AppTextInput
              value={fields.dimension}
              placeholder="overworld"
              onChange={(e) => setField('dimension')(e.target.value)}
            />
          </div>

          <div className={styles.coords}>
            {(['x', 'y', 'z'] as const).map((axis) => (
              <div className={styles.field} key={axis}>
                <AppText className={styles.label}>{axis.toUpperCase()}</AppText>
                <AppNumericInput
                  value={fields[axis]}
                  onChange={(e) => setField(axis)(e.target.value)}
                />
              </div>
            ))}
          </div>

          {isError && <AppText className={styles.error}>{`Ошибка: ${error?.message}`}</AppText>}
        </form>
      </AppCard>
    </section>
  )
}

export default CreateStockpilePage