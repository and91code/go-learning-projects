import { Folder, PlusCircle } from "lucide-react";
import { useEffect, useState } from "react";
import TaskForm from "../components/forms/TaskForm";
import TaskTable from "../components/tasks/TaskTable";
import Button from "../components/ui/Button";
import Modal from "../components/ui/Modal";
import { getLists } from "../services/listService";
import {
  createTask,
  deleteTask,
  getTasks,
  updateTask,
} from "../services/taskService";

const secondaryLinkClass =
  "inline-flex items-center justify-center rounded-lg border border-slate-300 bg-white px-4 py-2 font-medium text-slate-700 transition hover:bg-slate-50";

function fetchTasks(listId) {
  return getTasks(listId).then((result) =>
    Array.isArray(result) ? result : result?.data || result?.tasks || [],
  );
}

export default function MainView() {
  const [lists, setLists] = useState([]);
  const [selected, setSelected] = useState(null);
  const [tasks, setTasks] = useState([]);
  const [error, setError] = useState("");
  const [modal, setModal] = useState(null);
  const [editingTask, setEditingTask] = useState(null);

  useEffect(() => {
    getLists()
      .then((result) => {
        const items = Array.isArray(result)
          ? result
          : result?.data || result?.lists || [];
        setLists(items);
        setSelected(items[0] || null);
      })
      .catch((err) => setError(err.message));
  }, []);

  useEffect(() => {
    if (!selected) return;
    fetchTasks(selected.id)
      .then(setTasks)
      .catch((err) => setError(err.message));
  }, [selected]);

  async function refreshTasks() {
    if (!selected) return;
    try {
      setTasks(await fetchTasks(selected.id));
    } catch (err) {
      setError(err.message);
    }
  }

  async function addTask(data) {
    await createTask(data);
    await refreshTasks();
    setModal(null);
  }

  async function saveTask(data) {
    await updateTask(editingTask.id, data);
    await refreshTasks();
    setEditingTask(null);
    setModal(null);
  }

  async function removeTask(task) {
    if (!window.confirm(`Excluir a tarefa "${task.title}"?`)) return;
    try {
      await deleteTask(task.id);
      setTasks((current) => current.filter((item) => item.id !== task.id));
    } catch (err) {
      setError(err.message);
    }
  }

  async function changeStatus(task, status) {
    try {
      await updateTask(task.id, {
        title: task.title,
        description: task.description || "",
        status,
        priority: task.priority,
        tag_ids: task.tags?.map((tag) => tag.id) || [],
      });
      await refreshTasks();
    } catch (err) {
      setError(err.message);
    }
  }

  return (
    <div>
      {error && (
        <p
          className="mb-4 rounded-lg bg-red-50 p-3 text-sm text-red-700"
          role="alert"
        >
          {error}
          <button className="ml-3 font-semibold" onClick={() => setError("")}>
            Fechar
          </button>
        </p>
      )}
      <section className="mb-6 flex flex-wrap items-end justify-between gap-4 rounded-2xl bg-white p-5 shadow-sm">
        <label className="grid min-w-56 gap-1.5 text-sm font-medium text-slate-700">
          Lista selecionada
          <span className="flex items-center gap-2">
            <Folder size={18} className="shrink-0 text-slate-500" />
            <select
              className="w-full rounded-lg border border-slate-300 px-3 py-2"
              value={selected?.id || ""}
              onChange={(event) => {
                const list = lists.find(
                  (item) => String(item.id) === event.target.value,
                );
                setTasks([]);
                setSelected(list || null);
              }}
            >
              {!lists.length && <option value="">Nenhuma lista</option>}
              {lists.map((list) => (
                <option key={list.id} value={list.id}>
                  {list.title}
                </option>
              ))}
            </select>
          </span>
        </label>
        <div className="flex flex-wrap gap-2">
          <Button onClick={() => setModal("task")} disabled={!selected}>
            <PlusCircle size={18} className="mr-2" />
            Nova Tarefa
          </Button>
        </div>
      </section>
      <section className="rounded-2xl bg-white p-6 shadow-sm">
        <h1 className="mb-1 text-2xl font-bold">
          {selected?.title || "Minhas tarefas"}
        </h1>
        {selected?.description && (
          <p className="mb-4 text-slate-500">{selected.description}</p>
        )}
        <TaskTable
          tasks={tasks}
          onStatusChange={changeStatus}
          onEdit={(task) => {
            setEditingTask(task);
            setModal("editTask");
          }}
          onDelete={removeTask}
        />
      </section>
      {modal === "task" && selected && (
        <Modal title="Criar tarefa" onClose={() => setModal(null)}>
          <TaskForm
            listId={selected.id}
            onSubmit={addTask}
            onCancel={() => setModal(null)}
          />
        </Modal>
      )}
      {modal === "editTask" && selected && editingTask && (
        <Modal
          title="Editar tarefa"
          onClose={() => {
            setModal(null);
            setEditingTask(null);
          }}
        >
          <TaskForm
            listId={selected.id}
            task={editingTask}
            onSubmit={saveTask}
            onCancel={() => {
              setModal(null);
              setEditingTask(null);
            }}
          />
        </Modal>
      )}
    </div>
  );
}
