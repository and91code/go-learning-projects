import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react'
import { API_URL } from '../utils/constants'
import { storage } from '../utils/storage'
import { listNotifications } from '../services/notificationService'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [user, setUser] = useState(() => (storage.getToken() ? storage.getUser() : null))
  const [notificationVersion, setNotificationVersion] = useState(0)
  const [resourceVersions, setResourceVersions] = useState({ list: 0, task: 0, tag: 0 })
  const [notifications, setNotifications] = useState([])
  const [notificationLoadError, setNotificationLoadError] = useState("")

  useEffect(() => {
    const token = storage.getToken()
    if (!user || !token) return undefined
    let active = true
    listNotifications().then((result) => {
      if (!active) return
      const incoming = result?.notifications || []
      setNotifications((current) => {
        const merged = new Map(incoming.map((item) => [item.id, item]))
        current.forEach((item) => merged.set(item.id, item))
        return [...merged.values()].sort((a, b) => new Date(b.timestamp) - new Date(a.timestamp))
      })
    }).catch((error) => {
      if (active) setNotificationLoadError(error.message)
    })

    const events = new EventSource(`${API_URL}/events?token=${encodeURIComponent(token)}`)
    events.addEventListener('taskflow', (event) => {
      let notification
      try { notification = JSON.parse(event.data) } catch { return }
      if (!notification?.id) return
      setNotifications((current) => current.some((item) => item.id === notification.id) ? current : [notification, ...current])
      setNotificationVersion((current) => current + 1)
      setResourceVersions((current) => ({ ...current, [notification.resource_type]: (current[notification.resource_type] || 0) + 1 }))
    })
    return () => { active = false; events.close() }
  }, [user])

  const signIn = useCallback(async (action, credentials) => {
    const result = await action(credentials)
    const data = result?.data || result
    if (!data?.token || !data?.user) {
      throw new Error('Resposta de autenticação inválida.')
    }
    storage.setSession(data)
    setUser(data.user)
    setNotifications([])
    setNotificationLoadError("")
  }, [])

  const updateUser = useCallback((updatedUser) => {
    storage.setUser(updatedUser)
    setUser(updatedUser)
  }, [])

  const logout = useCallback(() => {
    storage.clearSession()
    setUser(null)
    setNotificationVersion(0)
    setNotifications([])
    setNotificationLoadError("")
    setResourceVersions({ list: 0, task: 0, tag: 0 })
  }, [])

  const setNotification = useCallback((notification) => setNotifications((current) => current.map((item) => item.id === notification.id ? notification : item)), [])
  const removeNotification = useCallback((id) => setNotifications((current) => current.filter((item) => item.id !== id)), [])
  const hasUnreadNotification = notifications.some((item) => item.unread)
  const value = useMemo(() => ({
    user,
    signIn,
    updateUser,
    logout,
    hasUnreadNotification,
    notificationVersion,
    resourceVersions,
    notifications,
    notificationLoadError,
    setNotification,
    removeNotification,
  }), [user, signIn, updateUser, logout, hasUnreadNotification, notificationVersion, resourceVersions, notifications, notificationLoadError, setNotification, removeNotification])

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) throw new Error('useAuth deve ser usado dentro de AuthProvider.')
  return context
}
