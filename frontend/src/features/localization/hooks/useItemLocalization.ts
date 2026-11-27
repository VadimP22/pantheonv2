import { useQuery } from '@tanstack/react-query'
import { getItemLocalization } from '../../../shared/api/itemLocalizations.ts'
import { itemLocalizationKeys } from '../../../shared/api/keys.ts'

/** Язык локализации по умолчанию. */
const DEFAULT_LANG = 'ru'

/** Хук подгрузки локализованного названия предмета по его item_id. */
export function useItemLocalization(itemId: string | null) {
  return useQuery({
    queryKey: itemLocalizationKeys.detail(DEFAULT_LANG, itemId ?? ''),
    queryFn: () => getItemLocalization(DEFAULT_LANG, itemId as string),
    enabled: itemId !== null,
  })
}
