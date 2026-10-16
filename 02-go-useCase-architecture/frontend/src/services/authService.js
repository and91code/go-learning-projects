import { post } from './api'

export const login = (credentials) => post('/auth/login', credentials)
export const register = (user) => post('/auth/register', user)
