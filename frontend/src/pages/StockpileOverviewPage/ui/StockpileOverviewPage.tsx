import { useNavigate, useParams } from 'react-router-dom'
import PageHeader from '../../../shared/PageHeader/ui/PageHeader.tsx'
import AppButton from '../../../shared/AppButton/ui/AppButton.tsx'
import { formatCoordinates } from '../../../shared/lib/formatCoordinates.ts'
import Inventory from '../../../widgets/Inventory/ui/Inventory.tsx'
import { useStockpile } from '../../../features/Stockpile/hooks/useStockpile.ts'

export function StockpileOverviewPage() {
  const navigate = useNavigate()
  const params = useParams()
  const id = Number(params.id)
  const { data, isError, error } = useStockpile(id)

  const isValidId = Number.isFinite(id) && id > 0
  const title = data?.name ?? 'Обзор склада'
  const description = !isValidId
    ? 'Некорректный идентификатор склада.'
    : isError
      ? `Не удалось загрузить склад: ${error.message}`
      : data
        ? `${data.dimension} · ${formatCoordinates(data.x, data.y, data.z)}`
        : 'Загрузка данных склада…'

  return (
    <section>
      <PageHeader
        breadcrumbs="Главная / Склады / Обзор"
        title={title}
        description={description}
        actions={
          <>
            <AppButton primary onClick={() => navigate(`/stockpiles/${id}/deposit`)}>
              {'Внести ресурс'}
            </AppButton>
            <AppButton onClick={() => navigate('/stockpiles')}>{'К складам'}</AppButton>
          </>
        }
      />
      <Inventory
        stockpileId={id}
        onDeposit={(itemId) => navigate(`/stockpiles/${id}/deposit/${itemId}`)}
      />
    </section>
  )
}

export default StockpileOverviewPage
