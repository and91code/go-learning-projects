import { get, patch, remove } from './api'

export const listNotifications = () => get('/notifications')
export const markNotificationRead = (id) => patch(`/notifications/${encodeURIComponent(id)}/read`)
export const deleteNotification = (id) => remove(`/notifications/${encodeURIComponent(id)}`)
