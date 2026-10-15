import { useState } from 'react'
import { LogIn, UserPlus } from 'lucide-react'
import { Link, Navigate, useLocation } from 'react-router-dom'
import Button from '../components/ui/Button'
import Input from '../components/ui/Input'
import { useAuth } from '../context/AuthContext'
import { register } from '../services/authService'

export default function RegisterView() {
  const { user, signIn } = useAuth()
  const location = useLocation()
  const [form, setForm] = useState({ name: '', email: '', password: '' })
  const [error, setError] = useState('')
  const [busy, setBusy] = useState(false)
  async function submit(event) {
    event.preventDefault()
    setBusy(true)
    setError('')
    try { await signIn(register, form) } catch (err) { setError(err.message) } finally { setBusy(false) }
  }
  if (user) return <Navigate to={location.state?.from?.pathname || '/main'} replace />
  return <div className="grid min-h-full place-items-center p-4">
    <section className="w-full max-w-md rounded-2xl bg-white p-8 shadow-sm">
      <h1 className="mb-2 text-3xl font-bold text-slate-900">Criar conta</h1>
      <p className="mb-6 text-slate-600">Comece a organizar suas tarefas.</p>
      <form onSubmit={submit} className="grid gap-4">
        <Input label="Nome" autoComplete="name" value={form.name} onChange={(event) => setForm({ ...form, name: event.target.value })} required minLength={2} maxLength={100} />
        <Input label="E-mail" type="email" autoComplete="email" value={form.email} onChange={(event) => setForm({ ...form, email: event.target.value })} required />
        <Input label="Senha" type="password" autoComplete="new-password" value={form.password} onChange={(event) => setForm({ ...form, password: event.target.value })} required minLength={6} />
        {error && <p className="text-sm text-red-600" role="alert">{error}</p>}
        <Button disabled={busy}><UserPlus size={18} className="mr-2" />{busy ? 'Criando…' : 'Criar conta'}</Button>
      </form>
      <p className="mt-5 text-center text-sm text-slate-600">Já tem conta? <Link className="font-semibold text-indigo-600" to="/login"><LogIn size={18} className="mr-1 inline" />Entrar</Link></p>
    </section>
  </div>
}
