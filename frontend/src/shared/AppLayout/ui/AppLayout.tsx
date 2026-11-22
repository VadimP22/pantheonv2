import type { ReactNode } from 'react'
import { Outlet } from 'react-router-dom'
import styles from './AppLayout.module.css'

export interface AppLayoutProps {
  /** Содержимое тёмного левого сайдбара (навигация). */
  sidebar?: ReactNode
  /** Содержимое верхней панели правой части. */
  header?: ReactNode
  /** Внешние классы для корневого элемента. */
  className?: string
}

export function AppLayout({ sidebar, header, className }: AppLayoutProps) {
  const rootClassName = className ? `${styles.root} ${className}` : styles.root

  return (
    <div className={rootClassName}>
      {sidebar != null && <aside className={styles.sidebar}>{sidebar}</aside>}
      <div className={styles.mainPane}>
        {header != null && <header className={styles.header}>{header}</header>}
        <main className={styles.content}>
          <Outlet />
        </main>
      </div>
    </div>
  )
}

export default AppLayout
