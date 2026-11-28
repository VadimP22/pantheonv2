import { useNavigate } from 'react-router-dom'
import { AppCard } from '../../../shared/AppCard/ui/AppCard.tsx'
import { AppTitle } from '../../../shared/AppTitle/ui/AppTitle.tsx'
import { AppText } from '../../../shared/AppText/ui/AppText.tsx'
import { AppProgressBar } from '../../../shared/AppProgressBar/ui/AppProgressBar.tsx'
import AppMinecraftIcon from '../../../shared/AppMinecraftIcon/ui/AppMinecraftIcon.tsx'
import AppTextButton from '../../../shared/AppTextButton/ui/AppTextButton.tsx'
import { useTaskStockpile } from '../hooks/useTaskStockpile.ts'
import styles from './TaskCard.module.css'

export interface TaskCardProps {
  /** Название задачи. */
  name: string
  /** Описание задачи (необязательное). */
  description?: string
  /** Имя minecraft-иконки (например `minecraft:crafting_table`). */
  iconItemName: string
  /** Признак «чекбокс-задачи»: true — статус выполнено/не выполнено, false — прогресс-бар. */
  isCheckbox: boolean
  /** Текущее количество выполненных единиц (для прогресс-бара). */
  itemCount?: number
  /** Целевое количество единиц (для прогресс-бара). */
  targetItemCount?: number
  /** Признак выполнения задачи (для чекбокс-задач). */
  done?: boolean
  /** Внешние классы для корневой карточки. */
  className?: string
  /** Идентификатор целевого склада; undefined — склад не задан. */
  stockpileId?: number
  /** Обработчик нажатия на кнопку «Внести»; undefined — кнопка не рендерится */
  onDepositClick?: () => void
  /** Обработчик нажатия на кнопку «Завершить»; undefined — кнопка не рендерится */
  onCompleteClick?: () => void
}

/**
 * Карточка задачи: иконка, название, описание, под описанием — статус
 * «Выполнено / Не выполнено» или прогресс до цели; внизу — ссылка на склад.
 */
export function TaskCard({
  name,
  description,
  iconItemName,
  isCheckbox,
  itemCount,
  targetItemCount,
  done,
  className,
  stockpileId,
  onDepositClick,
  onCompleteClick,
}: TaskCardProps) {
  const rootClassName = className ? `${styles.root} ${className}` : styles.root
  const navigate = useNavigate()
  const targetStockpileId = stockpileId ?? 0
  const { data: stockpile } = useTaskStockpile(targetStockpileId)
  const progress = targetItemCount && targetItemCount > 0
    ? Math.round(((itemCount ?? 0) / targetItemCount) * 100)
    : 0

  return (
    <AppCard className={rootClassName}>
      <div className={styles.header}>
        <div className={styles.icon}>
          <AppMinecraftIcon name={iconItemName} size={40} title={name} />
        </div>
        {onDepositClick || onCompleteClick ? (
          <div className={styles.deposit}>
            {onCompleteClick && (
              <AppTextButton variant="danger" onClick={onCompleteClick}>
                {'Завершить'}
              </AppTextButton>
            )}
            {onDepositClick && (
              <AppTextButton onClick={onDepositClick}>{'Внести'}</AppTextButton>
            )}
          </div>
        ) : null}
      </div>
      <AppTitle level={2} className={styles.title}>
        {name}
      </AppTitle>
      {description && <AppText className={styles.description}>{description}</AppText>}
      <div className={styles.statusRow}>
        {isCheckbox ? (
          <AppText className={done ? styles.statusDone : styles.statusPending}>
            {done ? 'Выполнено' : 'Не выполнено'}
          </AppText>
        ) : (
          <div className={styles.progress}>
            <AppProgressBar value={progress} />
            <AppText className={styles.progressCaption}>
              {itemCount ?? 0} / {targetItemCount ?? 0}
            </AppText>
          </div>
        )}
      </div>
      <div className={styles.bottom}>
        {targetStockpileId > 0 ? (
          <AppTextButton
            className={styles.stockpileLink}
            onClick={() => navigate(`/stockpiles/${targetStockpileId}`)}
          >
            {stockpile?.name ?? `Склад #${targetStockpileId}`}
          </AppTextButton>
        ) : (
          <AppText className={styles.noStockpile}>{'Нет склада'}</AppText>
        )}
      </div>
    </AppCard>
  )
}

export default TaskCard
