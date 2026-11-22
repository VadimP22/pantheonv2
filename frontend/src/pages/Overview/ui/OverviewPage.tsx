import PageHeader from '../../../shared/PageHeader/ui/PageHeader.tsx'

export function OverviewPage() {
  return (
    <section>
      <PageHeader
        breadcrumbs="Главная / Обзор"
        title="Обзор"
        description="Сводка по наличию и движению запасов на складах."
      />
    </section>
  )
}

export default OverviewPage