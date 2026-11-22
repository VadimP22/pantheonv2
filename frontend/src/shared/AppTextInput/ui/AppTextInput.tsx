import type { InputHTMLAttributes } from 'react'
import styles from './AppTextInput.module.css'

/**
 * Свойства однострочного текстового ввода.
 * Расширяет стандартные атрибуты `<input>`, поддерживая `value`/`onChange`/`placeholder`/`disabled` и т.д.
 */
export type AppTextInputProps = InputHTMLAttributes<HTMLInputElement>

export function AppTextInput({ className, type = 'text', ...rest }: AppTextInputProps) {
  const rootClassName = className ? `${styles.root} ${className}` : styles.root

  return (
    <span className={rootClassName}>
      <span className={styles.topEdge} aria-hidden="true" />
      <input type={type} className={styles.input} {...rest} />
      <span className={styles.bottomEdge} aria-hidden="true" />
    </span>
  )
}

export default AppTextInput
