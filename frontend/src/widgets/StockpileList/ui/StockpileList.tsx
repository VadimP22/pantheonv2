import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import StockpileCard from '../../../entities/StockpileComponents/ui/StockpileCard.tsx'
import AppText from '../../../shared/AppText/ui/AppText.tsx'
import AppButton from '../../../shared/AppButton/ui/AppButton.tsx'
import { formatCoordinates } from '../../../shared/lib/formatCoordinates.ts'
import { useStockpiles } from '../hooks/useStockpiles.ts'
import styles from './StockpileList.module.css'

/** Количество карточек на странице: сетка 3×2. */
const PAGE_SIZE = 6

export interface StockpileListProps {
  /** Внешние классы для корневого контейнера. */
  className?: string
}

/**
 * Виджет списка складов: самодостаточная сетка 3×2 с пагинацией.
 * Данные подгружаются с сервера через TanStack Query.
 */
export function StockpileList({ className }: StockpileListProps) {
  const navigate = useNavigate()
  const [page, setPage] = useState(1)
  const { data, isPending, isError, error } = useStockpiles({ page, pageSize: PAGE_SIZE })
  const rootClassName = className ? `${styles.root} ${className}` : styles.root

  if (isPending && !data) {
    return <AppText className={rootClassName}>{'Загрузка складов…'}</AppText>
  }

  if (isError) {
    return (
      <AppText className={rootClassName}>
        {`Не удалось загрузить склады: ${error.message}`}
      </AppText>
    )
  }

  const totalPages = data?.total_pages ?? 1

  return (
    <section className={rootClassName}>
      <div className={styles.grid}>
        {(data?.items ?? []).map((stockpile) => (
          <StockpileCard
            key={stockpile.id}
            dimension={stockpile.dimension}
            name={stockpile.name}
            coordinates={formatCoordinates(stockpile.x, stockpile.y, stockpile.z)}
            onDetailsClick={() => navigate(`/stockpiles/${stockpile.id}`)}
          />
        ))}
      </div>
      <nav className={styles.paginator} aria-label="Пагинация">
        <AppButton disabled={page <= 1} onClick={() => setPage((p) => Math.max(1, p - 1))}>
          {'‹ Назад'}
        </AppButton>
        <AppText className={styles.pageIndicator}>
          {page} / {totalPages}
        </AppText>
        <AppButton
          primary
          disabled={page >= totalPages}
          onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
        >
          {'Вперёд ›'}
        </AppButton>
      </nav>
    </section>
  )
}

export default StockpileList