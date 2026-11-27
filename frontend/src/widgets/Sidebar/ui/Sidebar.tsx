import { Icon } from '@imumd/mcicons/react'
import type { IconDefinition } from '@imumd/mcicons/svg'
import { useLocation, useNavigate } from 'react-router-dom'
import emerald from '@imumd/mcicons/definitions/emerald'
import chestMinecart from '@imumd/mcicons/definitions/chest-minecart'
import writableBook from '@imumd/mcicons/definitions/writable-book'
import { SidebarButton } from '../../../entities/SidebarComponents/ui/sidebarButton.tsx'
import AppThemeSwitcher from '../../../features/theme/ui/AppThemeSwitcher.tsx'
import styles from './Sidebar.module.css'

export type SidebarSection = 'overview' | 'stockpiles' | 'tasks'

const sections: { id: SidebarSection; to: string; label: string; icon: IconDefinition }[] = [
  { id: 'overview', to: '/', label: 'Обзор', icon: emerald },
  { id: 'stockpiles', to: '/stockpiles', label: 'Склады', icon: chestMinecart },
  { id: 'tasks', to: '/tasks', label: 'Задачи', icon: writableBook },
]

export function Sidebar() {
  const navigate = useNavigate()
  const { pathname } = useLocation()
  const active: SidebarSection = pathname.startsWith('/stockpiles')
    ? 'stockpiles'
    : pathname.startsWith('/tasks')
      ? 'tasks'
      : 'overview'

  return (
    <nav className={styles.sidebar} aria-label="Навигация">
      {sections.map(({ id, to, label, icon }) => (
        <SidebarButton
          key={id}
          label={label}
          icon={<Icon icon={icon} size={20} />}
          selected={active === id}
          onClick={() => navigate(to)}
        />
      ))}
      <div className={styles.footer}>
        <AppThemeSwitcher />
      </div>
    </nav>
  )
}

export default Sidebar