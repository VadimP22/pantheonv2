import type { HTMLAttributes, ReactNode } from 'react'
import styles from './AppBanner.module.css'

export interface AppBannerProps extends HTMLAttributes<HTMLDivElement> {
  /** Вложенное содержимое баннера. */
  children?: ReactNode
}

export function AppBanner({ children, className, ...rest }: AppBannerProps) {
  const rootClassName = className ? `${styles.root} ${className}` : styles.root

  return (
    <div className={rootClassName} {...rest}>
      <div className={styles.topEdge} />
      <div className={styles.bottomEdge} />
      <div className={styles.content}>{children}</div>
    </div>
  )
}

export default AppBanner
