import { AlertTriangle, Calendar, CheckCircle2, ChevronDown, Clock, Loader, Minus, Pencil, Tag, Trash2 } from 'lucide-react'

const statusLabels = { pending: 'Pendente', in_progress: 'Em andamento', completed: 'Concluída' }

const statusIcons = {
  pending: <Clock size={16} className="mr-1 text-amber-500" />,
  in_progress: <Loader size={16} className="mr-1 animate-spin-slow text-blue-500" />,
  completed: <CheckCircle2 size={16} className="mr-1 text-emerald-500" />,
}
const priorityIcons = {
  high: <AlertTriangle size={16} className="mr-1 text-red-500" />,
  medium: <Minus size={16} className="mr-1 text-yellow-500" />,
  low: <ChevronDown size={16} className="mr-1 text-gray-400" />,
}

const cellClass = 'px-3 py-3 text-sm text-slate-600 align-top whitespace-nowrap'

export default function TaskTable({ tasks = [], onStatusChange, onEdit, onDelete }) {
  if (tasks.length === 0) return <p className="py-8 text-center text-slate-500">Nenhuma tarefa nesta lista.</p>
  return <div className="overflow-x-auto">
    <table className="w-full text-left">
      <thead><tr className="border-b text-sm text-slate-500">
        <th className="px-3 py-3 font-medium">ID</th>
        <th className="px-3 py-3 font-medium">Lista</th>
        <th className="px-3 py-3 font-medium">Título</th>
        <th className="px-3 py-3 font-medium">Descrição</th>
        <th className="px-3 py-3 font-medium">Status</th>
        <th className="px-3 py-3 font-medium">Prioridade</th>
        <th className="px-3 py-3 font-medium">Data limite</th>
        <th className="px-3 py-3 font-medium">Tags</th>
        <th className="px-3 py-3 font-medium">Criada em</th>
        <th className="px-3 py-3 font-medium">Atualizada em</th>
        <th className="px-3 py-3 text-right font-medium">Ações</th>
      </tr></thead>
      <tbody className="divide-y divide-slate-100">
        {tasks.map((task) => <tr key={task.id} id={`task-${task.id}`}>
          <td className={cellClass}>{task.id}</td>
          <td className={cellClass}>{task.list_id}</td>
          <td className="px-3 py-3"><div className="flex items-center gap-3">
            <input type="checkbox" checked={task.status === 'completed'} onChange={() => onStatusChange(task, task.status === 'completed' ? 'pending' : 'completed')} aria-label={`Concluir ${task.title}`} />
            <span className={task.status === 'completed' ? 'text-slate-400 line-through' : 'font-medium text-slate-800'}>{task.title}</span>
          </div></td>
          <td className="max-w-xs px-3 py-3 text-sm text-slate-600"><span className="line-clamp-2" title={task.description}>{task.description || '—'}</span></td>
          <td className={cellClass}><span className="inline-flex items-center">{statusIcons[task.status]}{statusLabels[task.status] || task.status}</span></td>
          <td className={`${cellClass} capitalize`}><span className="inline-flex items-center">{priorityIcons[task.priority]}{task.priority}</span></td>
          <td className={cellClass}>{task.due_date ? <span className="inline-flex items-center"><Calendar size={14} className="mr-1 text-gray-400" />{task.due_date}</span> : '—'}</td>
          <td className="px-3 py-3"><div className="flex flex-wrap gap-1.5">
            {task.tags?.map((tag) => <span key={tag.id} className="inline-flex items-center rounded-full border px-2 py-0.5 text-xs font-medium" style={{ borderColor: tag.color, color: tag.color, backgroundColor: `color-mix(in srgb, ${tag.color} 15%, white)` }}><Tag size={12} className="mr-1" />{tag.name}</span>)}
          </div></td>
          <td className={cellClass}>{task.created_at || '—'}</td>
          <td className={cellClass}>{task.updated_at || '—'}</td>
          <td className="px-3 py-3"><div className="flex justify-end gap-3">
            <button className="inline-flex items-center gap-1 text-sm font-medium text-blue-600 hover:text-blue-800" onClick={() => onEdit(task)} aria-label={`Editar ${task.title}`}><Pencil size={16} />Editar</button>
            <button className="inline-flex items-center gap-1 text-sm font-medium text-red-600 hover:text-red-800" onClick={() => onDelete(task)} aria-label={`Excluir ${task.title}`}><Trash2 size={16} />Excluir</button>
          </div></td>
        </tr>)}
      </tbody>
    </table>
  </div>
}
