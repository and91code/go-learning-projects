import { get, post, put, remove } from './api'

export const getTagsByList = (listId) => get(`/tags?list_id=${encodeURIComponent(listId)}`)
export const createTag = (data) => post('/tags', data)
export const updateTag = (id, data) => put(`/tags/${id}`, data)
export const deleteTag = (id) => remove(`/tags/${id}`)
