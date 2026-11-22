import type { ReactNode } from 'react'
import styles from './InventorySlot.module.css'

export interface InventorySlotProps {
  /** Контент слота (обычно иконка), центрируется внутри ячейки. */
  children?: ReactNode
  /** Признак выделенного слота: меняет фон на светлый акцентный. */
  selected?: boolean
  /** Обработчик клика по слоту; выделение управляется снаружи. */
  onClick?: () => void
}

export function InventorySlot({ children, selected = false, onClick }: InventorySlotProps) {
  const rootClassName = selected ? `${styles.root} ${styles.selected}` : styles.root

  return (
    <div className={rootClassName} onClick={onClick}>
      <span className={`${styles.edge} ${styles.edgeTop}`} aria-hidden="true" />
      <span className={`${styles.edge} ${styles.edgeBottom}`} aria-hidden="true" />
      <div className={styles.slot}>{children}</div>
    </div>
  )
}

export default InventorySlot
