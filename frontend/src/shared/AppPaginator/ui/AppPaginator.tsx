import AppButton from '../../AppButton/ui/AppButton.tsx'
import AppText from '../../AppText/ui/AppText.tsx'
import styles from './AppPaginator.module.css'

export interface AppPaginatorProps {
  /** Текущая страница (с 1). */
  page: number
  /** Общее количество страниц (>= 1). */
  totalPages: number
  /** Callback при смене страницы. */
  onPageChange: (page: number) => void
  /** Внешние классы для корневого nav. */
  className?: string
}

/** Общий пагинатор: кнопки «Назад/Вперёд» и индикатор «N / M». */
export function AppPaginator({ page, totalPages, onPageChange, className }: AppPaginatorProps) {
  const rootClassName = className ? `${styles.paginator} ${className}` : styles.paginator

  return (
    <nav className={rootClassName} aria-label="Пагинация">
      <AppButton disabled={page <= 1} onClick={() => onPageChange(Math.max(1, page - 1))}>
        {'‹ Назад'}
      </AppButton>
      <AppText className={styles.pageIndicator}>
        {page} / {totalPages}
      </AppText>
      <AppButton
        primary
        disabled={page >= totalPages}
        onClick={() => onPageChange(Math.min(totalPages, page + 1))}
      >
        {'Вперёд ›'}
      </AppButton>
    </nav>
  )
}

export default AppPaginator
