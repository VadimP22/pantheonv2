import type { SelectHTMLAttributes, ReactNode } from 'react'
import styles from './AppDropDown.module.css'

/** Свойства выпадающего списка. */
export interface AppDropDownProps extends SelectHTMLAttributes<HTMLSelectElement> {
  /** Подпись над селектом. */
  label?: ReactNode
  /** Плейсхолдер — пустой <option> в начале списка. */
  placeholder?: string
}

export function AppDropDown({ label, placeholder, className, children, ...rest }: AppDropDownProps) {
  const rootClassName = className ? `${styles.root} ${className}` : styles.root

  return (
    <span className={rootClassName}>
      {label ? <span className={styles.label}>{label}</span> : null}
      <span className={styles.selectWrapper}>
        <span className={styles.topEdge} aria-hidden="true" />
        <select className={styles.select} {...rest}>
          {placeholder ? <option value="">{placeholder}</option> : null}
          {children}
        </select>
        <span className={styles.bottomEdge} aria-hidden="true" />
      </span>
    </span>
  )
}

export default AppDropDown
