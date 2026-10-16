import { Outlet } from 'react-router-dom'
import { Header } from '../ui/Header'
import Breadcrumbs from './Breadcrumbs'

export default function AppLayout() {
  return (
    <div className="h-screen w-screen flex flex-col bg-gray-50 overflow-hidden">
      <Header />
      <div className="border-b border-gray-200 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-3">
          <Breadcrumbs />
        </div>
      </div>
      <main className="flex-1 overflow-y-auto py-6">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-full">
          <Outlet />
        </div>
      </main>
    </div>
  )
}
