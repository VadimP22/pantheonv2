import PageHeader from '../../../shared/PageHeader/ui/PageHeader.tsx'
import { OverviewCards } from '../../../widgets/OverviewCards/ui/OverviewCards.tsx'

export function OverviewPage() {
  return (
    <section>
      <PageHeader
        breadcrumbs="Главная / Обзор"
        title="Обзор"
        description="Сводка по наличию и движению запасов на складах."
      />
      <OverviewCards />
    </section>
  )
}

export default OverviewPage
