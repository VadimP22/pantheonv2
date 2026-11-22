import { Route, Routes } from 'react-router-dom'
import AppLayout from './shared/AppLayout/ui/AppLayout.tsx'
import Sidebar from './widgets/Sidebar/ui/Sidebar.tsx'
import OverviewPage from './pages/Overview/ui/OverviewPage.tsx'
import StockpilesPage from './pages/StockpilesPage/ui/StockpilesPage.tsx'
import CreateStockpilePage from './pages/CreateStockpilePage/ui/CreateStockpilePage.tsx'
import StockpileOverviewPage from './pages/StockpileOverviewPage/ui/StockpileOverviewPage.tsx'
import DepositResourcePage from './pages/DepositResourcePage/ui/DepositResourcePage.tsx'
import DepositSpecifiedResourcePage from './pages/DepositSpecifiedResourcePage/ui/DepositSpecifiedResourcePage.tsx'

function App() {
  return (
    <Routes>
      <Route element={<AppLayout sidebar={<Sidebar />} />}>
        <Route index element={<OverviewPage />} />
        <Route path="stockpiles" element={<StockpilesPage />} />
        <Route path="stockpiles/create" element={<CreateStockpilePage />} />
        <Route path="stockpiles/:id" element={<StockpileOverviewPage />} />
        <Route path="stockpiles/:id/deposit" element={<DepositResourcePage />} />
        <Route
          path="stockpiles/:id/deposit/:item_id"
          element={<DepositSpecifiedResourcePage />}
        />
      </Route>
    </Routes>
  )
}

export default App