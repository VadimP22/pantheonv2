import type { ButtonHTMLAttributes, ReactNode } from 'react'
import styles from './AppTextButton.module.css'

export interface AppTextButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  /** Вложенное содержимое (текст кнопки). */
  children?: ReactNode
  /** Визуальный вариант кнопки. 'default' — акцентный, 'danger' — красный (var(--danger)). */
  variant?: 'default' | 'danger'
}

export function AppTextButton({ className, children, variant = 'default', ...rest }: AppTextButtonProps) {
  const variantClass = variant === 'danger' ? styles.danger : undefined
  const rootClassName = [styles.root, className, variantClass].filter(Boolean).join(' ')
  return (
    <button type="button" className={rootClassName} {...rest}>
      {children}
    </button>
  )
}

export default AppTextButton