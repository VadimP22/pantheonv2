import type { ButtonHTMLAttributes, ReactNode } from 'react'
import styles from './sidebarButton.module.css'

export interface SidebarButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  /** Текст кнопки. */
  label: string
  /** Слот иконки (заполняется позже из ассетов дизайна). */
  icon?: ReactNode
  /** Активное состояние пункта сайдбара. */
  selected?: boolean
  /** Внешние классы для корневого элемента. */
  className?: string
}

export function SidebarButton({
  label,
  icon,
  selected = false,
  className,
  ...rest
}: SidebarButtonProps) {
  let rootClassName = styles.root
  if (selected) rootClassName += ` ${styles.selected}`
  if (className) rootClassName += ` ${className}`

  return (
    <button
      type="button"
      className={rootClassName}
      aria-pressed={selected}
      {...rest}
    >
      {(icon != null || label != null) && (
        <span className={styles.content}>
          {icon != null && <span className={styles.icon}>{icon}</span>}
          {label}
        </span>
      )}
    </button>
  )
}

export default SidebarButton