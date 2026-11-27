import type { HTMLAttributes } from 'react'
import styles from './AppHorizontalLine.module.css'

export interface AppHorizontalLineProps extends HTMLAttributes<HTMLDivElement> {
  /** Внешние классы для корневой линии. */
  className?: string
}

export function AppHorizontalLine({ className, ...rest }: AppHorizontalLineProps) {
  const rootClassName = className ? `${styles.root} ${className}` : styles.root

  return <div className={rootClassName} {...rest} />
}

export default AppHorizontalLine