/** Длительность минуты в миллисекундах. */
const MINUTE_MS = 60 * 1000
/** Длительность часа в миллисекундах. */
const HOUR_MS = 60 * MINUTE_MS
/** Длительность суток в миллисекундах. */
const DAY_MS = 24 * HOUR_MS

/** Выбирает русскую форму слова для числа (1 час / 2 часа / 5 часов). */
function pluralizeRu(n: number, forms: [one: string, few: string, many: string]): string {
  const abs = Math.abs(n) % 100
  const last = abs % 10
  if (abs > 10 && abs < 20) return forms[2]
  if (last > 1 && last < 5) return forms[1]
  if (last === 1) return forms[0]
  return forms[2]
}

/**
 * Форматирует время последнего изменения склада:
 * «обновлено недавно» — если прошло меньше часа;
 * «изменено N часов назад» — если прошло меньше суток;
 * «изменено ДД.ММ.ГГГГ» — если прошли сутки и больше.
 * Возвращает пустую строку для некорректной или отсутствующей даты.
 */
export function formatUpdatedAt(iso: string): string {
  if (!iso) return ''
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return ''

  const diff = Date.now() - date.getTime()
  if (diff < HOUR_MS) return 'обновлено недавно'
  if (diff < DAY_MS) {
    const hours = Math.max(1, Math.floor(diff / HOUR_MS))
    return `изменено ${hours} ${pluralizeRu(hours, ['час', 'часа', 'часов'])} назад`
  }

  const formatted = new Intl.DateTimeFormat('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  }).format(date)
  return `изменено ${formatted}`
}
