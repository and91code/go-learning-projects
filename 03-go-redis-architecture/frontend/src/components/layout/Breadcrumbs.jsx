import { ChevronRight } from 'lucide-react'
import { Link, useLocation } from 'react-router-dom'

const pageLabels = {
  '/main': 'Início',
  '/my-profile': 'Meu perfil',
  '/lists': 'Listas',
  '/tags': 'Tags',
  '/login': 'Login',
  '/register': 'Cadastro',
}

export default function Breadcrumbs() {
  const { pathname } = useLocation()
  const label = pageLabels[pathname] || pathname
  const home = { label: 'Início', to: '/main' }
  const crumbs = pathname === '/main' ? [{ label: 'Início' }] : [home, { label }]

  return <nav aria-label="Breadcrumb">
    <ol className="flex flex-wrap items-center gap-1 text-sm text-slate-500">
      {crumbs.map((crumb, index) => {
        const isLast = index === crumbs.length - 1
        return <li key={crumb.label} className="flex items-center gap-1">
          {index > 0 && <ChevronRight size={14} className="text-slate-400" aria-hidden="true" />}
          {crumb.to && !isLast
            ? <Link to={crumb.to} className="hover:text-indigo-600">{crumb.label}</Link>
            : <span className={isLast ? 'font-medium text-slate-800' : ''} aria-current={isLast ? 'page' : undefined}>{crumb.label}</span>}
        </li>
      })}
    </ol>
  </nav>
}
