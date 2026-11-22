import type { InputHTMLAttributes } from 'react'
import styles from './AppNumericInput.module.css'

/**
 * Свойства числового ввода.
 * Расширяет стандартные атрибуты `<input type="number">`, поддерживая `value`/`onChange`/`placeholder`/`disabled` и т.д.
 */
export type AppNumericInputProps = InputHTMLAttributes<HTMLInputElement>

export function AppNumericInput({ className, step = 'any', ...rest }: AppNumericInputProps) {
  const rootClassName = className ? `${styles.root} ${className}` : styles.root

  return (
    <span className={rootClassName}>
      <span className={styles.topEdge} aria-hidden="true" />
      <input type="number" className={styles.input} step={step} {...rest} />
      <span className={styles.bottomEdge} aria-hidden="true" />
    </span>
  )
}

export default AppNumericInput