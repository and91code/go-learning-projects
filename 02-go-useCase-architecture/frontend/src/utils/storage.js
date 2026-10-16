import { STORAGE_KEYS } from './constants'

export const storage = {
  getToken: () => localStorage.getItem(STORAGE_KEYS.token),
  getUser: () => {
    const user = localStorage.getItem(STORAGE_KEYS.user)
    return user ? JSON.parse(user) : null
  },
  setUser: (user) => {
    localStorage.setItem(STORAGE_KEYS.user, JSON.stringify(user))
  },
  setSession: ({ token, user }) => {
    localStorage.setItem(STORAGE_KEYS.token, token)
    localStorage.setItem(STORAGE_KEYS.user, JSON.stringify(user))
  },
  clearSession: () => {
    localStorage.removeItem(STORAGE_KEYS.token)
    localStorage.removeItem(STORAGE_KEYS.user)
  },
  clear: () => {
    storage.clearSession()
  },
}
