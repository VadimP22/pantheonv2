import { AppMinecraftIcon } from '../../../shared/AppMinecraftIcon/ui/AppMinecraftIcon.tsx'
import { InventorySlot } from './InventorySlot.tsx'
import styles from './InventorySlotWithItem.module.css'

export interface InventorySlotWithItemProps {
  /** ID предмета (имя иконки, например `minecraft:diamond`). Если не задано — слот пустой. */
  itemId?: string
  /** Количество предметов в слоте, отображается в правом нижнем углу. */
  count?: number
  /** Признак выделенного слота: фон меняется на акцентный `#a5bd88`. */
  selected?: boolean
  /** Обработчик клика по слоту; выделение переключает родитель. */
  onClick?: () => void
}

export function InventorySlotWithItem({
  itemId,
  count,
  selected = false,
  onClick,
}: InventorySlotWithItemProps) {
  const showCount = count != null && count > 0

  return (
    <InventorySlot selected={selected} onClick={onClick}>
      {itemId ? <AppMinecraftIcon name={itemId} size={48} /> : null}
      {showCount ? <span className={styles.count}>{count}</span> : null}
    </InventorySlot>
  )
}

export default InventorySlotWithItem
