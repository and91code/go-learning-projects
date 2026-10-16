import { useState } from 'react'
import { Check, X } from 'lucide-react'
import Button from '../ui/Button'
import Input from '../ui/Input'

export default function ListForm({ list, onSubmit, onCancel }) {
  const [title, setTitle] = useState(list?.title || '')
  const [description, setDescription] = useState(list?.description || '')
  const [error, setError] = useState('')
  async function submit(event) {
    event.preventDefault()
    setError('')
    try {
      await onSubmit({ title, description })
      if (!list) {
        setTitle('')
        setDescription('')
      }
    } catch (err) {
      setError(err.message)
    }
  }
  return (
    <form onSubmit={submit} className="grid gap-4">
      <Input label="Título" value={title} onChange={(event) => setTitle(event.target.value)} required minLength={3} maxLength={100} />
      <Input label="Descrição" value={description} onChange={(event) => setDescription(event.target.value)} maxLength={255} />
      {error && <p className="w-full text-sm text-red-600" role="alert">{error}</p>}
      <div className="flex justify-end gap-2">
        {onCancel && <Button type="button" variant="secondary" onClick={onCancel}><X size={16} className="mr-2" />Cancelar</Button>}
        <Button type="submit"><Check size={16} className="mr-2" />{list ? 'Salvar alterações' : 'Criar lista'}</Button>
      </div>
    </form>
  )
}
