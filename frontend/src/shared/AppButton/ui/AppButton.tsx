import type { ButtonHTMLAttributes, ReactNode } from 'react'
import styles from './AppButton.module.css'

export interface AppButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  /** Оформить кнопку в основном (зелёном) виде. */
  primary?: boolean
  /** Текст кнопки. */
  children?: ReactNode
}

export function AppButton({ primary = false, className, children, ...rest }: AppButtonProps) {
  const rootClassName = [
    styles.root,
    primary ? styles.primary : '',
    className ?? '',
  ].join(' ').trim()

  return (
    <button type="button" className={rootClassName} {...rest}>
      <span className={styles.topEdge} aria-hidden="true" />
      <span className={styles.bottomEdge} aria-hidden="true" />
      <span className={styles.content}>{children}</span>
    </button>
  )
}

export default AppButton