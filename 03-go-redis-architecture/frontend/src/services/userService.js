import { put } from './api'

export const updateProfile = ({ name, email }) => put('/user/profile', { name, email })
export const updatePassword = ({ currentPassword, newPassword }) => put('/user/password', {
  current_password: currentPassword,
  new_password: newPassword,
})
