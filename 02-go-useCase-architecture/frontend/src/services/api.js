import { API_URL } from '../utils/constants'
import { storage } from '../utils/storage'

export async function api(path, options = {}) {
  const headers = new Headers(options.headers)
  headers.set('Content-Type', 'application/json')
  const token = storage.getToken()
  if (token) headers.set('Authorization', `Bearer ${token}`)

  const response = await fetch(`${API_URL}${path}`, { ...options, headers })
  const body = await response.json().catch(() => null)
  if (!response.ok) {
    const message = body?.error || body?.message || `Erro HTTP ${response.status}`
    throw new Error(message)
  }
  return body
}

export const get = (path) => api(path)
export const post = (path, data) => api(path, { method: 'POST', body: JSON.stringify(data) })
export const put = (path, data) => api(path, { method: 'PUT', body: JSON.stringify(data) })
export const remove = (path) => api(path, { method: 'DELETE' })
