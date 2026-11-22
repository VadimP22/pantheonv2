/** Форматирует координаты склада в строку вида «X / Y / Z     111 / 222 / 333». */
export function formatCoordinates(x: number, y: number, z: number): string {
  return `X / Y / Z\u00A0\u00A0${x} / ${y} / ${z}`
}
