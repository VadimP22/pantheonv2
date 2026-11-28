import styles from './SidebarSectionLabel.module.css'

export interface SidebarSectionLabelProps {
  /** Текст метки секции (например, «УПРАВЛЕНИЕ»). */
  children: string
}

/** Заголовок-метка секции сайдбара («УПРАВЛЕНИЕ», «СИСТЕМА»). */
export function SidebarSectionLabel({ children }: SidebarSectionLabelProps) {
  return <div className={styles.label}>{children}</div>
}

export default SidebarSectionLabel
