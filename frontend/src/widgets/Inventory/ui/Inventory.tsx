import { useState } from 'react'
import AppCard from '../../../shared/AppCard/ui/AppCard.tsx'
import AppText from '../../../shared/AppText/ui/AppText.tsx'
import AppTitle from '../../../shared/AppTitle/ui/AppTitle.tsx'
import AppButton from '../../../shared/AppButton/ui/AppButton.tsx'
import AppMinecraftIcon from '../../../shared/AppMinecraftIcon/ui/AppMinecraftIcon.tsx'
import InventorySlotWithItem from '../../../entities/InventoryComponents/ui/InventorySlotWithItem.tsx'
import { useStockpileResources } from '../../../features/Stockpile/hooks/useStockpileResources.ts'
import styles from './Inventory.module.css'

/** Размер одного стака в штуках (стандарт Minecraft). */
const STACK_SIZE = 64

/** Предмет инвентаря: идентификатор иконки и количество. */
export interface InventoryItem {
  itemId: string
  count: number
}

export interface InventoryProps {
  /** Идентификатор склада, ресурсы которого загружаются. */
  stockpileId: number
  /** Количество слотов на странице: сетка 9×3. */
  pageSize?: number
  /** Внешние классы для корневой карточки. */
  className?: string
  /** Обработчик нажатия «внести этот ресурс»; вызывается с id выбранного предмета. */
  onDeposit?: (itemId: string) => void
}

/**
 * Виджет инвентаря: карточка 9×3 слота с пагинацией снизу.
 * Ресурсы склада подгружаются с сервера по stockpileId через useStockpileResources.
 */
export function Inventory({ stockpileId, pageSize = 27, className, onDeposit }: InventoryProps) {
  const { data } = useStockpileResources(stockpileId)
  const [page, setPage] = useState(1)
  const [selectedItemId, setSelectedItemId] = useState<string | null>(null)

  const items: InventoryItem[] = (data ?? []).map((resource) => ({
    itemId: resource.item_id,
    count: resource.item_count,
  }))
  const totalPages = Math.max(1, Math.ceil(items.length / pageSize))
  const currentPage = Math.min(page, totalPages)
  const start = (currentPage - 1) * pageSize
  const pageItems = items.slice(start, start + pageSize)
  const selectedItem = items.find((item) => item.itemId === selectedItemId) ?? null
  const selectedItemStacks = selectedItem
    ? Math.floor(selectedItem.count / STACK_SIZE)
    : 0
  const wrapperClassName = className ? `${styles.wrapper} ${className}` : styles.wrapper
  const handleDeposit = () => {
    if (selectedItem) onDeposit?.(selectedItem.itemId)
  }

  return (
    <div className={wrapperClassName}>
      <AppCard secondary className={styles.root}>
        <div className={styles.grid}>
          {Array.from({ length: pageSize }, (_, index) => {
            const item = pageItems[index]
            return item ? (
              <InventorySlotWithItem
                key={index}
                itemId={item.itemId}
                count={item.count}
                selected={item.itemId === selectedItemId}
                onClick={() =>
                  setSelectedItemId(selectedItemId === item.itemId ? null : item.itemId)
                }
              />
            ) : (
              <InventorySlotWithItem key={index} onClick={() => setSelectedItemId(null)} />
            )
          })}
        </div>
        <nav className={styles.paginator} aria-label="Пагинация инвентаря">
          <AppButton
            disabled={currentPage <= 1}
            onClick={() => setPage(Math.max(1, currentPage - 1))}
          >
            {'‹ Назад'}
          </AppButton>
          <AppText className={styles.pageIndicator}>
            {currentPage} / {totalPages}
          </AppText>
          <AppButton
            primary
            disabled={currentPage >= totalPages}
            onClick={() => setPage(Math.min(totalPages, currentPage + 1))}
          >
            {'Вперёд ›'}
          </AppButton>
        </nav>
      </AppCard>
      <AppCard className={styles.detail}>
        {selectedItem ? (
          <div className={styles.detailColumn}>
            <div className={styles.detailIcon}>
              <AppMinecraftIcon name={selectedItem.itemId} size={48} />
            </div>

            <AppTitle className={styles.detailTitle}>{selectedItem.itemId}</AppTitle>
            <AppText className={styles.detailId}>{selectedItem.itemId}</AppText>

            <AppTitle className={styles.detailQty}>{selectedItem.count}</AppTitle>
            <AppText className={styles.detailStacks}>
              {`${selectedItemStacks} стаков`}
            </AppText>

            <div className={styles.detailActions}>
              <AppButton disabled={!selectedItem} onClick={handleDeposit}>
                {'Внести'}
              </AppButton>
            </div>
          </div>
        ) : null}
      </AppCard>
    </div>
  )
}

export default Inventory
