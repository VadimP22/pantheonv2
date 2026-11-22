import { useNavigate } from 'react-router-dom'
import PageHeader from '../../../shared/PageHeader/ui/PageHeader.tsx'
import StockpileList from '../../../widgets/StockpileList/ui/StockpileList.tsx'
import AppButton from '../../../shared/AppButton/ui/AppButton.tsx'

export function StockpilesPage() {
  const navigate = useNavigate()
  return (
    <section>
      <PageHeader
        breadcrumbs="Главная / Склады"
        title="Склады"
        description="Список складов и структура их хранения."
        actions={
          <div style={{ display: 'flex' }}>
            <AppButton primary onClick={() => navigate('/stockpiles/create')}>Создать склад</AppButton>
          </div>
        }
      />
      <StockpileList />
    </section>
  )
}

export default StockpilesPage