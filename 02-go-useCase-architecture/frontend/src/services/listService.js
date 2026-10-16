import { get, post, put, remove } from './api'

export const getLists = () => get('/list')
export const createList = (data) => post('/list', data)
export const updateList = (id, data) => put(`/list/${id}`, data)
export const deleteList = (id) => remove(`/list/${id}`)
