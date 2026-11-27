import type { ButtonHTMLAttributes, ReactNode } from 'react'
import styles from './AppTextButton.module.css'

export interface AppTextButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  /** Вложенное содержимое (текст кнопки). */
  children?: ReactNode
}

export function AppTextButton({ className, children, ...rest }: AppTextButtonProps) {
  const rootClassName = className ? `${styles.root} ${className}` : styles.root

  return (
    <button type="button" className={rootClassName} {...rest}>
      {children}
    </button>
  )
}

export default AppTextButton