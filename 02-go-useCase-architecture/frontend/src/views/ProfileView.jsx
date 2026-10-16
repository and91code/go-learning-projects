import { useState } from 'react'
import { ArrowLeft, KeyRound, Lock, Save } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import Button from '../components/ui/Button'
import Input from '../components/ui/Input'
import { useAuth } from '../context/AuthContext'
import { updatePassword, updateProfile } from '../services/userService'

export default function ProfileView() {
  const { user, updateUser } = useAuth()
  const navigate = useNavigate()
  const [name, setName] = useState(user?.name || '')
  const [email, setEmail] = useState(user?.email || '')
  const [profileMessage, setProfileMessage] = useState('')
  const [profileError, setProfileError] = useState('')
  const [passwordMessage, setPasswordMessage] = useState('')
  const [passwordError, setPasswordError] = useState('')
  const [passwords, setPasswords] = useState({ currentPassword: '', newPassword: '', confirmation: '' })

  async function saveProfile(event) {
    event.preventDefault()
    setProfileMessage('')
    setProfileError('')
    try {
      const result = await updateProfile({ name, email })
      const updatedUser = result?.user || result?.data?.user || result?.data || result
      if (!updatedUser || typeof updatedUser !== 'object' || !updatedUser.name || !updatedUser.email) {
        throw new Error('Resposta inválida ao atualizar o perfil.')
      }
      updateUser(updatedUser)
      setName(updatedUser.name)
      setEmail(updatedUser.email)
      setProfileMessage('Perfil atualizado com sucesso.')
    } catch (error) {
      setProfileError(error.message)
    }
  }

  async function savePassword(event) {
    event.preventDefault()
    setPasswordMessage('')
    setPasswordError('')
    if (passwords.newPassword.length < 6) {
      setPasswordError('A nova senha deve ter pelo menos 6 caracteres.')
      return
    }
    if (passwords.newPassword !== passwords.confirmation) {
      setPasswordError('A nova senha e a confirmação não coincidem.')
      return
    }
    try {
      await updatePassword({
        currentPassword: passwords.currentPassword,
        newPassword: passwords.newPassword,
      })
      setPasswords({ currentPassword: '', newPassword: '', confirmation: '' })
      setPasswordMessage('Senha atualizada com sucesso.')
    } catch (error) {
      setPasswordError(error.message)
    }
  }

  return (
    <div>
      <div className="mx-auto grid max-w-4xl gap-6">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <h1 className="text-2xl font-bold">Meu Perfil</h1>
          <Button variant="secondary" onClick={() => navigate('/main')}><ArrowLeft size={18} className="mr-2" />Voltar</Button>
        </div>
        <section className="rounded-2xl bg-white p-6 shadow-sm">
          <h2 className="mb-5 text-lg font-semibold">Dados pessoais</h2>
          <form onSubmit={saveProfile} className="grid gap-4 sm:grid-cols-2">
            <Input label="Nome" value={name} onChange={(event) => setName(event.target.value)} required minLength={2} maxLength={100} />
            <Input label="E-mail" type="email" value={email} onChange={(event) => setEmail(event.target.value)} required />
            {(profileError || profileMessage) && <p className={`sm:col-span-2 text-sm ${profileError ? 'text-red-600' : 'text-emerald-700'}`} role="status">{profileError || profileMessage}</p>}
            <div className="sm:col-span-2">
              <Button type="submit"><Save size={16} className="mr-2" />Salvar Alterações</Button>
            </div>
          </form>
        </section>
        <section className="rounded-2xl bg-white p-6 shadow-sm">
          <h2 className="mb-5 text-lg font-semibold">Alteração de senha</h2>
          <form onSubmit={savePassword} className="grid gap-4 sm:grid-cols-2">
            <label className="grid gap-1.5 text-sm font-medium text-slate-700 sm:col-span-2">Senha Atual
              <span className="relative">
                <Lock size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
                <input className="w-full rounded-lg border border-slate-300 py-2 pl-9 pr-3 outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-100" type="password" autoComplete="current-password" value={passwords.currentPassword} onChange={(event) => setPasswords({ ...passwords, currentPassword: event.target.value })} required />
              </span>
            </label>
            <Input label="Nova Senha" type="password" autoComplete="new-password" value={passwords.newPassword} onChange={(event) => setPasswords({ ...passwords, newPassword: event.target.value })} required minLength={6} />
            <Input label="Confirmar Nova Senha" type="password" autoComplete="new-password" value={passwords.confirmation} onChange={(event) => setPasswords({ ...passwords, confirmation: event.target.value })} required minLength={6} />
            {(passwordError || passwordMessage) && <p className={`sm:col-span-2 text-sm ${passwordError ? 'text-red-600' : 'text-emerald-700'}`} role="status">{passwordError || passwordMessage}</p>}
            <div className="sm:col-span-2">
              <Button type="submit"><KeyRound size={16} className="mr-2" />Atualizar Senha</Button>
            </div>
          </form>
        </section>
      </div>
    </div>
  )
}
