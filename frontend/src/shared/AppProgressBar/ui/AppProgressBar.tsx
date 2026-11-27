import type { CSSProperties } from 'react'
import styles from './AppProgressBar.module.css'

export interface AppProgressBarProps {
  /**
   * Текущее значение в процентах 0–100.
   * Значения вне диапазона клампируются по ширине; значение больше 100 считается
   * переполнением — бар красится токенами `--progress-overflow-*` (золото в светлой
   * теме, фиолетовый в тёмной). Внутренний край заливки всегда остаётся внутри трека.
   */
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
  const overflow = value > 100
  const innerWidth = `calc(${progress}% - 4px)`

  const fillStyle: CSSProperties = { width: innerWidth }
  const glintStyle: CSSProperties = { width: innerWidth }

  const rootClassName = [styles.root, overflow ? styles.overflow : '', className ?? ''].join(' ').trim()

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
