import type { IconDefinition } from '@imumd/mcicons/svg'
import barrier from '@imumd/mcicons/definitions/barrier'

// Карта ленивых загрузчиков определений иконок библиотеки @imumd/mcicons.
//
// `import.meta.glob` из Vite принимает только пути, начинающиеся с `/` или `./`,
// поэтому указываем каталог с определениями внутри node_modules относительно корня
// проекта (`frontend/`). Без `eager` каждый файл превращается в отдельный чанк,
// который подгружается по требованию — в бандл попадают только реально показанные иконки.
const iconLoaders = import.meta.glob<IconDefinition>(
  '/node_modules/@imumd/mcicons/dist/definitions/*.mjs',
  { import: 'default' },
)

// Имя иконки (kebab-case) → загрузчик. Ключом служит имя файла без расширения.
const loadersByName = new Map<string, () => Promise<IconDefinition>>(
  Object.entries(iconLoaders).map(([path, load]) => [
    path.split('/').pop()!.replace(/\.mjs$/, ''),
    load,
  ]),
)

// Кэш уже загруженных определений: повторные запросы одной иконки не грузят чанк снова.
const loadedIcons = new Map<string, IconDefinition>()

/** Иконка-заглушка для неизвестных имён. */
export const fallbackMinecraftIcon: IconDefinition = barrier

/**
 * Приводит имя иконки к канону библиотеки @imumd/mcicons.
 *
 * Снимает namespace (`minecraft:` / `minecraft_`), приводит к нижнему регистру
 * и заменяет `_` на `-`: файлы определений в библиотеке названы в kebab-case
 * (`acacia-fence-gate`), поэтому `minecraft:acacia_fence` → `acacia-fence`.
 */
export function normalizeMinecraftIconName(name: string): string {
  return name
    .trim()
    .toLowerCase()
    .replace(/^minecraft(?::|_)/, '')
    .replace(/_/g, '-')
}

/** Возвращает уже загруженное определение, если иконка успела подгрузиться. */
export function getLoadedMinecraftIcon(name: string): IconDefinition | undefined {
  return loadedIcons.get(normalizeMinecraftIconName(name))
}

/**
 * Асинхронно возвращает определение для имени.
 * Неизвестное имя (или ошибка загрузки) → fallback-иконка.
 */
export async function loadMinecraftIcon(name: string): Promise<IconDefinition> {
  const key = normalizeMinecraftIconName(name)

  const cached = loadedIcons.get(key)
  if (cached) return cached

  const load = loadersByName.get(key)
  if (!load) return fallbackMinecraftIcon

  const icon = await load()
  loadedIcons.set(key, icon)
  return icon
}
