import { PrimeReactProvider } from 'primereact/api'
import React from 'react'
import { Route, Routes } from 'react-router-dom'
import { RecoilRoot } from 'recoil'
import { Layout } from './components/Layout'
import { NotFound } from './components/NotFound'
import { Building } from './pages/Building'
import BuildingsPage from './pages/Buildings'
import { ConstructBuilding } from './pages/ConstructBuilding'
import LandsPage from './pages/Lands'
import LogisticsPage from './pages/Logistics'
import { Map } from './pages/Map'
import OrdersPage from './pages/Orders'
import ResourcesPage from './pages/Resources'
import SearchBuilding from './pages/SearchBuilding'
import StoragesPage from './pages/Storages'
import { Store } from './pages/Store'
import { User } from './pages/User'
import 'primereact/resources/themes/saga-blue/theme.css'

import './i18n'

const App = () => {
  return (
    <PrimeReactProvider>
      <RecoilRoot>
        <Routes>
          <Route path="/" element={<Layout/>}>
            <Route index element={<User/>}/>
            <Route path="/map" element={<Map/>}/>
            <Route path="/buildings" element={<BuildingsPage/>}/>
            <Route path="/lands" element={<LandsPage/>}/>
            <Route path="/resources" element={<ResourcesPage/>}/>
            <Route path="/logistics" element={<LogisticsPage/>}/>
            <Route path="/storages" element={<StoragesPage/>}/>
            <Route path="/orders" element={<OrdersPage/>}/>
            <Route path="/building/:buildingID" element={<Building/>}/>
            <Route path="/store/:buildingID" element={<Store/>}/>
            <Route path="/construct_building/:xParam/:yParam" element={<ConstructBuilding/>}/>
            <Route path="/search-building" element={<SearchBuilding/>}/>
            <Route path="*" element={<NotFound/>}/>
          </Route>
        </Routes>
      </RecoilRoot>
    </PrimeReactProvider>
  )
}

export default App
