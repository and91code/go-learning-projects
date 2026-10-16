import { CheckSquare, List, LogIn, LogOut, Tag, UserCog } from 'lucide-react'
import { NavLink, useNavigate } from 'react-router-dom'
import { useAuth } from '../../context/AuthContext'

const navLinkClass = ({ isActive }) => `flex items-center gap-1.5 rounded-lg px-2 py-1.5 text-sm font-medium hover:bg-indigo-50 ${isActive ? 'text-indigo-700 bg-indigo-50' : 'text-gray-600'}`

export function Header() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  function handleLogout() {
    logout()
    navigate('/login', { replace: true })
  }

  return (
    <header className="sticky top-0 z-50 bg-white border-b border-gray-200 shadow-sm">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between gap-4">
        <div className="flex items-center gap-6">
          <NavLink to={user ? '/main' : '/login'} className="text-xl font-bold text-indigo-600 flex items-center gap-2">
            <CheckSquare className="text-indigo-600" size={24} />
            TaskFlow
          </NavLink>
          {user && (
            <nav className="flex items-center gap-1 sm:gap-2" aria-label="Principal">
              <NavLink to="/lists" className={navLinkClass} aria-label="Listas">
                <List size={18} />
                <span className="hidden sm:inline">Listas</span>
              </NavLink>
              <NavLink to="/tags" className={navLinkClass} aria-label="Tags">
                <Tag size={18} />
                <span className="hidden sm:inline">Tags</span>
              </NavLink>
            </nav>
          )}
        </div>

        <nav className="flex items-center gap-2 sm:gap-4" aria-label="Conta">
          {user !== null ? (
            <>
              <span className="hidden md:inline text-sm font-medium text-gray-700">{user.name}</span>
              <NavLink to="/my-profile" className={navLinkClass} aria-label="Meu perfil">
                <UserCog size={18} />
                <span className="hidden sm:inline">Meu perfil</span>
              </NavLink>
              <button
                type="button"
                onClick={handleLogout}
                aria-label="Sair"
                className="flex items-center gap-1.5 rounded-lg px-2 py-1.5 text-sm font-medium text-gray-600 hover:bg-gray-100"
              >
                <LogOut size={18} />
                <span className="hidden sm:inline">Sair</span>
              </button>
            </>
          ) : (
            <NavLink to="/login" className={navLinkClass} aria-label="Entrar">
              <LogIn size={18} />
              <span className="hidden sm:inline">Entrar</span>
            </NavLink>
          )}
        </nav>
      </div>
    </header>
  )
}
