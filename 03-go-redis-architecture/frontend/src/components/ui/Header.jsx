import { Bell, CheckSquare, List, LogIn, LogOut, Tag, UserCog, Trash2 } from 'lucide-react'
import { NavLink, useNavigate } from 'react-router-dom'
import { useAuth } from '../../context/AuthContext'
import { useState } from 'react'
import Modal from './Modal'
import { deleteNotification, markNotificationRead } from '../../services/notificationService'

const navLinkClass = ({ isActive }) => `flex items-center gap-1.5 rounded-lg px-2 py-1.5 text-sm font-medium hover:bg-indigo-50 ${isActive ? 'text-indigo-700 bg-indigo-50' : 'text-gray-600'}`

export function Header() {
  const { user, logout, hasUnreadNotification, notifications, notificationLoadError, setNotification, removeNotification } = useAuth()
  const navigate = useNavigate()
  const [notificationsOpen, setNotificationsOpen] = useState(false)
  const [notificationError, setNotificationError] = useState("")

  async function markRead(notification) {
    try {
      if (notification.unread) setNotification(await markNotificationRead(notification.id))
      setNotificationError("")
    } catch (error) {
      setNotificationError(error.message)
    }
  }

  async function remove(notification) {
    try {
      await deleteNotification(notification.id)
      removeNotification(notification.id)
      setNotificationError("")
    } catch (error) {
      setNotificationError(error.message)
    }
  }

  function openNotification(notification) {
    const routes = {
      list: `/lists?list_id=${notification.resource_id}`,
      task: `/main?task_id=${notification.resource_id}`,
      tag: `/tags?tag_id=${notification.resource_id}`,
    }
    if (notification.unread) markRead(notification)
    setNotificationsOpen(false)
    navigate(routes[notification.resource_type] || '/main')
  }

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
                onClick={() => setNotificationsOpen(true)}
                aria-label="Notificações"
                className="relative rounded-lg p-2 text-gray-600 hover:bg-gray-100"
              >
                <Bell size={20} />
                {hasUnreadNotification && (
                  <span className="absolute right-1 top-1 h-2.5 w-2.5 animate-pulse rounded-full bg-red-500" />
                )}
              </button>
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
      {notificationsOpen && (
        <Modal title="Notificações" onClose={() => setNotificationsOpen(false)}>
          <div className="max-h-96 space-y-2 overflow-y-auto">
            {(notificationError || notificationLoadError) && <p role="alert" className="text-sm text-red-600">{notificationError || notificationLoadError}</p>}
            {notifications.length === 0 ? <p className="text-sm text-gray-500">Nenhuma notificação.</p> : notifications.map((item) => (
              <div key={item.id} className={`flex items-start gap-3 rounded-lg border p-3 ${item.unread ? 'border-indigo-200 bg-indigo-50' : 'border-gray-200'}`}>
                <button type="button" className="flex-1 text-left" onClick={() => openNotification(item)}>
                  <span className="block text-sm font-medium text-gray-800">{item.message}</span>
                  <span className="mt-1 block text-xs text-gray-500">{new Date(item.timestamp).toLocaleString('pt-BR')} · {item.resource_type} #{item.resource_id}</span>
                </button>
                {item.unread && <button type="button" className="text-xs text-indigo-700" onClick={() => markRead(item)} aria-label="Marcar como lida">Lida</button>}
                <button type="button" className="text-gray-500 hover:text-red-600" onClick={() => remove(item)} aria-label="Excluir notificação"><Trash2 size={16} /></button>
              </div>
            ))}
          </div>
        </Modal>
      )}
    </header>
  )
}
