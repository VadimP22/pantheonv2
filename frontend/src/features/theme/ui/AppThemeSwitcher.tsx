import { Icon } from '@imumd/mcicons/react'
import type { IconDefinition } from '@imumd/mcicons/svg'
import sunflower from '@imumd/mcicons/definitions/sunflower'
import torch from '@imumd/mcicons/definitions/torch'
import { useThemeStore } from '../store/themeStore.ts'
import styles from './AppThemeSwitcher.module.css'

/** Набор (иконка, текст) для каждого режима: показываем, куда переключаемся. */
const options: Record<'light' | 'dark', { icon: IconDefinition; label: string }> = {
  light: { icon: torch, label: 'Тёмная тема' },
  dark: { icon: sunflower, label: 'Светлая тема' },
}

/** Переключатель светлой/тёмной темы; ставится в нижней части сайдбара. */
export function AppThemeSwitcher() {
  const { theme, toggleTheme } = useThemeStore()
  const { icon, label } = options[theme]

  return (
    <button type="button" className={styles.root} onClick={toggleTheme}>
      <span className={styles.content}>
        <span className={styles.icon}>
          <Icon icon={icon} size={20} />
        </span>
        {label}
      </span>
    </button>
  )
}

export default AppThemeSwitcher