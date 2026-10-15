import { createContext, useCallback, useContext, useMemo, useState } from 'react'
import { storage } from '../utils/storage'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [user, setUser] = useState(() => (storage.getToken() ? storage.getUser() : null))

  const signIn = useCallback(async (action, credentials) => {
    const result = await action(credentials)
    const data = result?.data || result
    if (!data?.token || !data?.user) {
      throw new Error('Resposta de autenticação inválida.')
    }
    storage.setSession(data)
    setUser(data.user)
  }, [])

  const updateUser = useCallback((updatedUser) => {
    storage.setUser(updatedUser)
    setUser(updatedUser)
  }, [])

  const logout = useCallback(() => {
    storage.clearSession()
    setUser(null)
  }, [])

  const value = useMemo(() => ({ user, signIn, updateUser, logout }), [user, signIn, updateUser, logout])

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) throw new Error('useAuth deve ser usado dentro de AuthProvider.')
  return context
}
