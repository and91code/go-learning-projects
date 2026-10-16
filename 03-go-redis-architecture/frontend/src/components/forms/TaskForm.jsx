import { useEffect, useRef, useState } from 'react'
import { Check, Plus, Tag, X } from 'lucide-react'
import { getTagsByList } from '../../services/tagService'
import Button from '../ui/Button'
import Input from '../ui/Input'

export default function TaskForm({ listId, task, onSubmit, onCancel }) {
  const [title, setTitle] = useState(task?.title || '')
  const [description, setDescription] = useState(task?.description || '')
  const [priority, setPriority] = useState(task?.priority || 'medium')
  const [availableTags, setAvailableTags] = useState([])
  const [selectedTagIds, setSelectedTagIds] = useState(() => task?.tags?.map((tag) => tag.id) || [])
  const [error, setError] = useState('')
  const previousListId = useRef(listId)

  useEffect(() => {
    let ignore = false
    if (previousListId.current !== listId) {
      previousListId.current = listId
      setSelectedTagIds([])
    }
    getTagsByList(listId)
      .then((result) => {
        if (ignore) return
        setAvailableTags(Array.isArray(result) ? result : result?.data || result?.tags || [])
      })
      .catch((err) => {
        if (!ignore) setError(err.message)
      })
    return () => {
      ignore = true
    }
  }, [listId])

  function toggleTag(tagId) {
    setSelectedTagIds((current) => current.includes(tagId) ? current.filter((id) => id !== tagId) : [...current, tagId])
  }

  async function submit(event) {
    event.preventDefault()
    setError('')
    try {
      const data = {
        title,
        description,
        priority,
        tag_ids: selectedTagIds,
      }
      if (task) data.status = task.status
      else data.list_id = Number(listId)
      await onSubmit(data)
      setTitle('')
      setDescription('')
      setSelectedTagIds([])
    } catch (err) {
      setError(err.message)
    }
  }
  return <form onSubmit={submit} className="grid gap-4">
    <Input label="Título" value={title} onChange={(event) => setTitle(event.target.value)} required minLength={3} maxLength={150} />
    <label className="grid gap-1.5 text-sm font-medium text-slate-700">Descrição
      <textarea className="rounded-lg border border-slate-300 px-3 py-2 outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-100" rows="3" value={description} onChange={(event) => setDescription(event.target.value)} />
    </label>
    <label className="grid gap-1.5 text-sm font-medium text-slate-700">Prioridade
      <select className="rounded-lg border border-slate-300 px-3 py-2" value={priority} onChange={(event) => setPriority(event.target.value)}>
        <option value="low">Baixa</option><option value="medium">Média</option><option value="high">Alta</option>
      </select>
    </label>
    <fieldset className="grid gap-2 text-sm font-medium text-slate-700">
      <legend className="mb-1.5 flex items-center gap-2"><Tag size={16} />Tags</legend>
      {availableTags.length === 0
        ? <p className="text-sm font-normal text-slate-500">Nenhuma tag cadastrada nesta lista</p>
        : <div className="flex flex-wrap gap-2">
          {availableTags.map((tag) => {
            const selected = selectedTagIds.includes(tag.id)
            return <button key={tag.id} type="button" aria-pressed={selected} onClick={() => toggleTag(tag.id)} className={`inline-flex items-center gap-2 rounded-full border px-3 py-1 text-sm transition-colors ${selected ? 'border-indigo-500 bg-indigo-50 text-indigo-700' : 'border-slate-300 text-slate-600 hover:bg-slate-50'}`}>
              <span className="h-2.5 w-2.5 rounded-full" style={{ backgroundColor: tag.color }} aria-hidden="true" />
              {tag.name}
            </button>
          })}
        </div>}
    </fieldset>
    {error && <p className="text-sm text-red-600" role="alert">{error}</p>}
    <div className="flex justify-end gap-2">
      {onCancel && <Button type="button" variant="secondary" onClick={onCancel}><X size={16} className="mr-2" />Cancelar</Button>}
      <Button type="submit">{task ? <Check size={16} className="mr-2" /> : <Plus size={16} className="mr-2" />}{task ? 'Salvar alterações' : 'Criar tarefa'}</Button>
    </div>
  </form>
}
