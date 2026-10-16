import { get, post, put, remove } from './api'

export const getTasks = (listId) => get(`/tasks?list_id=${encodeURIComponent(listId)}`)
export const getTaskById = (id) => get(`/tasks/${encodeURIComponent(id)}`)
export const createTask = (data) => post('/tasks', data)
export const updateTask = (id, data) => put(`/tasks/${id}`, data)
export const deleteTask = (id) => remove(`/tasks/${id}`)
