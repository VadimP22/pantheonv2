import { Icon } from '@imumd/mcicons/react'
import chestMinecart from '@imumd/mcicons/definitions/chest-minecart'
import { AppCard } from '../../../shared/AppCard/ui/AppCard.tsx'
import { AppTitle } from '../../../shared/AppTitle/ui/AppTitle.tsx'
import { AppText } from '../../../shared/AppText/ui/AppText.tsx'
import AppButton from '../../../shared/AppButton/ui/AppButton.tsx'
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
}

export function StockpileCard({
  dimension,
  name,
  coordinates,
  className,
  onDetailsClick,
}: StockpileCardProps) {
  const rootClassName = className ? `${styles.root} ${className}` : styles.root

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
      <div className={styles.action}>
        <AppButton onClick={onDetailsClick}>{'Подробнее'}</AppButton>
      </div>
    </AppCard>
  )
}

export default StockpileCard