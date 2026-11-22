import type { CSSProperties } from 'react'
import styles from './AppProgressBar.module.css'

export interface AppProgressBarProps {
  /** Текущее значение в процентах 0–100 (значения вне диапазона клампируются). */
  value: number
  /** Внешние классы для управления шириной/отступами (по умолчанию трек занимает 100% ширины). */
  className?: string
  /** Доступная подпись для role="progressbar". */
  ariaLabel?: string
}

function clampValue(value: number): number {
  if (!Number.isFinite(value)) return 0
  return Math.min(100, Math.max(0, value))
}

export function AppProgressBar({ value, className, ariaLabel = 'Прогресс' }: AppProgressBarProps) {
  const progress = clampValue(value)
  const progressWidth = `${progress}%`

  const fillStyle: CSSProperties = { width: progressWidth }
  const glintStyle: CSSProperties = { width: progressWidth }

  const rootClassName = className ? `${styles.root} ${className}` : styles.root

  return (
    <div
      className={rootClassName}
      role="progressbar"
      aria-valuemin={0}
      aria-valuemax={100}
      aria-valuenow={Math.round(progress)}
      aria-label={ariaLabel}
    >
      <div className={styles.track}>
        <div className={styles.fill} style={fillStyle} />
        <div className={styles.glint} style={glintStyle} />
      </div>
    </div>
  )
}

export default AppProgressBar
