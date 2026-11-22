import type { HTMLAttributes, ReactNode } from 'react'
import styles from './AppText.module.css'

export interface AppTextProps extends HTMLAttributes<HTMLParagraphElement> {
  /** Вложенное содержимое текста. */
  children?: ReactNode
}

export function AppText({ children, className, ...rest }: AppTextProps) {
  const rootClassName = className ? `${styles.root} ${className}` : styles.root

  return (
    <p className={rootClassName} {...rest}>
      {children}
    </p>
  )
}

export default AppText
