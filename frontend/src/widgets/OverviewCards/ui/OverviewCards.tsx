import { Icon } from '@imumd/mcicons/react'
import type { IconDefinition } from '@imumd/mcicons/svg'
import diamondPickaxe from '@imumd/mcicons/definitions/diamond-pickaxe'
import diamond from '@imumd/mcicons/definitions/diamond'
import ironIngot from '@imumd/mcicons/definitions/iron-ingot'
import book from '@imumd/mcicons/definitions/book'
import AppCard from '../../../shared/AppCard/ui/AppCard.tsx'
import AppText from '../../../shared/AppText/ui/AppText.tsx'
import AppTitle from '../../../shared/AppTitle/ui/AppTitle.tsx'
import { useTasksCount } from '../hooks/useTasksCount.ts'
import { useStockpilesCount } from '../hooks/useStockpilesCount.ts'
import styles from './OverviewCards.module.css'

interface StatCard {
  /** Иконка из mcicons. */
  icon: IconDefinition
  /** Подпись на карточке. */
  label: string
  /** Значение (число). */
  value: number | string
  /** true — значение загружается (показать «…»). */
  loading?: boolean
}

const MOCK_PROJECTS = 8
const MOCK_MEMBERS = 32

export function OverviewCards() {
  const tasksQuery = useTasksCount()
  const stockpilesQuery = useStockpilesCount()

  const cards: StatCard[] = [
    { icon: diamondPickaxe, label: 'Активных проектов', value: MOCK_PROJECTS },
    {
      icon: diamond,
      label: 'Открытых задач',
      value: tasksQuery.isPending ? '…' : (tasksQuery.data ?? 0),
      loading: tasksQuery.isPending,
    },
    {
      icon: ironIngot,
      label: 'Складов клана',
      value: stockpilesQuery.isPending ? '…' : (stockpilesQuery.data ?? 0),
      loading: stockpilesQuery.isPending,
    },
    { icon: book, label: 'Участника клана', value: MOCK_MEMBERS },
  ]

  return (
    <section className={styles.root} aria-label="Сводка клана">
      {cards.map((card) => (
        <AppCard key={card.label} className={styles.card}>
          <Icon icon={card.icon} size={40} className={styles.icon} aria-hidden="true" />
          <AppTitle level={2} className={styles.number}>
            {card.value}
          </AppTitle>
          <AppText className={styles.label}>{card.label}</AppText>
        </AppCard>
      ))}
    </section>
  )
}

export default OverviewCards
