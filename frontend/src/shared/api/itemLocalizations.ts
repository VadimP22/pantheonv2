// Модуль fetch-обёрток для HTTP API локализаций предметов.
import type { ItemLocalizationDTO } from '../generated/dto.ts'
import { request } from './http.ts'

/**
 * Возвращает локализованное название предмета.
 * GET /api/loc/{lang}/{item}
 */
export async function getItemLocalization(
  lang: string,
  itemId: string,
): Promise<ItemLocalizationDTO> {
  return request<ItemLocalizationDTO>(
    `/api/loc/${encodeURIComponent(lang)}/${encodeURIComponent(itemId)}`,
  )
}
