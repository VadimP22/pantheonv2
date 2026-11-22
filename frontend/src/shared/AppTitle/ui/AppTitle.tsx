import type { HTMLAttributes, ReactNode } from 'react'
import styles from './AppTitle.module.css'

export interface AppTitleProps extends HTMLAttributes<HTMLHeadingElement> {
  /** Уровень заголовка h1–h6 (по умолчанию h1). */
  level?: 1 | 2 | 3 | 4 | 5 | 6
  /** Вложенное содержимое заголовка. */
  children?: ReactNode
}

const headingElements = {
  1: 'h1',
  2: 'h2',
  3: 'h3',
  4: 'h4',
  5: 'h5',
  6: 'h6',
} as const

export function AppTitle({ level = 1, children, className, ...rest }: AppTitleProps) {
  const Heading = headingElements[level]
  const rootClassName = className ? `${styles.root} ${className}` : styles.root

  return (
    <Heading className={rootClassName} {...rest}>
      {children}
    </Heading>
  )
}

export default AppTitle
