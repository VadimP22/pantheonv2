import { Icon } from '@imumd/mcicons/react'
import chestMinecart from '@imumd/mcicons/definitions/chest-minecart'
import { AppCard } from '../../../shared/AppCard/ui/AppCard.tsx'
import { AppTitle } from '../../../shared/AppTitle/ui/AppTitle.tsx'
import { AppText } from '../../../shared/AppText/ui/AppText.tsx'
import AppTextButton from '../../../shared/AppTextButton/ui/AppTextButton.tsx'
import AppHorizontalLine from '../../../shared/AppHorizontalLine/ui/AppHorizontalLine.tsx'
import { formatUpdatedAt } from '../../../shared/lib/formatUpdatedAt.ts'
import styles from './StockpileCard.module.css'

export interface StockpileCardProps {
  /** Измерение (dimension) склада, отображается справа от иконки. */
  dimension: string
  /** Название склада. */
  name: string
  /** Координаты склада. */
  coordinates: string
  /** Внешние классы для корневой карточки. */
  className?: string
  /** Обработчик нажатия на кнопку «Подробнее». */
  onDetailsClick: () => void
  /** Время последнего изменения склада (ISO 8601). */
  updatedAt: string
  /** Общее количество предметов на складе. */
  itemCount: number
  /** Количество типов предметов на складе. */
  itemTypeCount: number
}

export function StockpileCard({
  dimension,
  name,
  coordinates,
  className,
  onDetailsClick,
  updatedAt,
  itemCount,
  itemTypeCount,
}: StockpileCardProps) {
  const rootClassName = className ? `${styles.root} ${className}` : styles.root
  const updatedAtLabel = formatUpdatedAt(updatedAt)

  return (
    <AppCard className={rootClassName}>
      <div className={styles.header}>
        <Icon icon={chestMinecart} size={40} className={styles.icon} />
        <AppText className={styles.dimension}>{dimension}</AppText>
      </div>
      <AppTitle level={2} className={styles.title}>
        {name}
      </AppTitle>
      <AppText className={styles.coordinates}>{coordinates}</AppText>
      <AppHorizontalLine className={styles.divider} />
      <div className={styles.counts}>
        <div className={styles.countBlock}>
          <AppTitle level={3} className={styles.countValue}>
            {itemCount}
          </AppTitle>
          <AppText className={styles.countCaption}>{'предметов'}</AppText>
        </div>
        <div className={styles.countBlockRight}>
          <AppTitle level={3} className={styles.countValue}>
            {itemTypeCount}
          </AppTitle>
          <AppText className={styles.countCaption}>{'типов'}</AppText>
        </div>
      </div>
      <div className={styles.action}>
        {updatedAtLabel && <AppText>{updatedAtLabel}</AppText>}
        <AppTextButton onClick={onDetailsClick}>{'Подробнее'}</AppTextButton>
      </div>
    </AppCard>
  )
}

export default StockpileCard