import { Icon } from '@imumd/mcicons/react'
import type { SVGProps } from 'react'
import { useMinecraftIcon } from '../hooks/useMinecraftIcon.ts'

export interface AppMinecraftIconProps {
  /** Имя иконки (например `minecraft:minecart`). Неизвестное имя вернёт fallback. */
  name: string
  /** Размер иконки. */
  size?: number | string
  /** Подпись для доступности. */
  title?: string
  /** Атрибуты SVG-иконки. */
  svgProps?: SVGProps<SVGSVGElement>
}

export function AppMinecraftIcon({ name, size, title, svgProps }: AppMinecraftIconProps) {
  const icon = useMinecraftIcon(name)

  // Пока определение иконки грузится, держим плейсхолдер того же размера,
  // чтобы не было скачка раскладки.
  if (!icon) {
    return (
      <span
        aria-hidden="true"
        style={{ display: 'inline-block', width: size ?? '1em', height: size ?? '1em' }}
      />
    )
  }

  return <Icon icon={icon} size={size} title={title} {...svgProps} />
}

export default AppMinecraftIcon
