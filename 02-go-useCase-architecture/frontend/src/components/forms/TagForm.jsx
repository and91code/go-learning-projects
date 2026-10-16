import { useState } from 'react'
import { Check, X } from 'lucide-react'
import Button from '../ui/Button'
import Input from '../ui/Input'

export default function TagForm({ listId, tag, onSubmit, onCancel }) {
  const [name, setName] = useState(tag?.name || '')
  const [color, setColor] = useState(tag?.color || '#6366f1')
  const [error, setError] = useState('')
  async function submit(event) {
    event.preventDefault()
    setError('')
    try {
      await onSubmit({ list_id: Number(listId), name, color })
      if (!tag) setName('')
    } catch (err) {
      setError(err.message)
    }
  }
  return <form className="grid gap-4" onSubmit={submit}>
    <Input label="Nome da tag" value={name} onChange={(event) => setName(event.target.value)} required maxLength={50} />
    <label className="grid gap-2 text-sm font-medium text-slate-700">Cor da tag
      <input aria-label="Cor da tag" type="color" value={color} onChange={(event) => setColor(event.target.value)} className="h-10 w-full" />
    </label>
    {error && <p className="text-sm text-red-600" role="alert">{error}</p>}
    <div className="flex justify-end gap-2">
      {onCancel && <Button type="button" variant="secondary" onClick={onCancel}><X size={16} className="mr-2" />Cancelar</Button>}
      <Button type="submit"><Check size={16} className="mr-2" />{tag ? 'Salvar alterações' : 'Criar tag'}</Button>
    </div>
  </form>
}
