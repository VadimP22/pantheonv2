import type { IconDefinition } from '@imumd/mcicons/svg'
import { useEffect, useState } from 'react'
import { getLoadedMinecraftIcon, loadMinecraftIcon } from '../lib/minecraftIconMap.ts'

/**
 * Лениво загружает определение иконки по имени.
 *
 * Пока иконка не загружена, возвращает `undefined`, чтобы вызывающий код успел
 * отрисовать плейсхолдер. Неизвестные имена в итоге дают fallback-иконку.
 */
export function useMinecraftIcon(name: string): IconDefinition | undefined {
  const [loaded, setLoaded] = useState<{ name: string; icon: IconDefinition } | undefined>(() => {
    const icon = getLoadedMinecraftIcon(name)
    return icon ? { name, icon } : undefined
  })

  useEffect(() => {
    // Иконка уже в кэше — грузить нечего, значение отдаём синхронно при рендере.
    if (getLoadedMinecraftIcon(name)) return

    let active = true
    loadMinecraftIcon(name).then((icon) => {
      if (active) setLoaded({ name, icon })
    })

    return () => {
      active = false
    }
  }, [name])

  // Свежезагруженный кэш приоритетнее состояния: так после эффекта не мелькает
  // плейсхолдер, а при смене имени не показывается иконка от предыдущего.
  const cached = getLoadedMinecraftIcon(name)
  if (cached) return cached

  return loaded?.name === name ? loaded.icon : undefined
}
