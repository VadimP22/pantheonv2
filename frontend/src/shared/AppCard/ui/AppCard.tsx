import type { HTMLAttributes, ReactNode } from 'react'
import styles from './AppCard.module.css'

export interface AppCardProps extends HTMLAttributes<HTMLDivElement> {
  /** Вложенное содержимое карточки. */
  children?: ReactNode
  /** Оформить карточку во вторичном (сером) виде кнопки. */
  secondary?: boolean
}

export function AppCard({ secondary = false, className, children, ...rest }: AppCardProps) {
  const rootClassName = [
    styles.root,
    secondary ? styles.secondary : '',
    className ?? '',
  ].join(' ').trim()

  return (
    <div className={rootClassName} {...rest}>
      <span className={styles.topEdge} aria-hidden="true" />
      <span className={styles.bottomEdge} aria-hidden="true" />
      <div className={styles.content}>{children}</div>
    </div>
  )
}

export default AppCard