import { Pencil, PlusCircle, Trash2 } from "lucide-react";
import { useEffect, useState } from "react";
import { useSearchParams } from "react-router-dom";
import ListForm from "../components/forms/ListForm";
import Button from "../components/ui/Button";
import Modal from "../components/ui/Modal";
import { useAuth } from "../context/AuthContext";
import {
  createList,
  deleteList,
  getLists,
  updateList,
} from "../services/listService";

export default function ListsView() {
  const { resourceVersions } = useAuth();
  const [searchParams] = useSearchParams();
  const [lists, setLists] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [modal, setModal] = useState(null);

  useEffect(() => {
    let ignore = false;
    getLists()
      .then((result) => {
        if (!ignore)
          setLists(Array.isArray(result) ? result : result?.data || []);
      })
      .catch((err) => {
        if (!ignore) setError(err.message);
      })
      .finally(() => {
        if (!ignore) setLoading(false);
      });
    return () => {
      ignore = true;
    };
  }, [resourceVersions.list]);

  useEffect(() => {
    const listId = searchParams.get("list_id");
    if (!listId || !lists.some((list) => String(list.id) === listId)) return;
    document.getElementById(`list-${listId}`)?.scrollIntoView({ behavior: "smooth", block: "center" });
  }, [lists, searchParams]);

  async function submit(data) {
    if (modal === "create") {
      const result = await createList(data);
      if (result?.status !== "accepted")
        setLists((current) => [...current, result?.data || result?.list || result]);
    } else {
      const result = await updateList(modal.id, data);
      const updated = result?.data || result?.list || result;
      if (result?.status !== "accepted")
        setLists((current) => current.map((item) => item.id === modal.id ? { ...item, ...updated } : item));
    }
    setModal(null);
  }

  async function remove(list) {
    if (
      !window.confirm(
        `Excluir a lista "${list.title}"? Suas tarefas e tags também serão excluídas.`,
      )
    )
      return;
    try {
      const result = await deleteList(list.id);
      if (result?.status !== "accepted")
        setLists((current) => current.filter((item) => item.id !== list.id));
    } catch (err) {
      setError(err.message);
    }
  }

  return (
    <div className="grid gap-6">
      <section className="flex flex-wrap items-center justify-between gap-4 rounded-2xl bg-white p-6 shadow-sm">
        <div>
          <h1 className="text-2xl font-bold">Listas</h1>
          <p className="text-slate-500">
            Gerencie as listas que organizam suas tarefas.
          </p>
        </div>
        <Button onClick={() => setModal("create")}>
          <PlusCircle size={18} className="mr-2" />
          Adicionar
        </Button>
      </section>

      <section className="rounded-2xl bg-white p-6 shadow-sm">
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
        {loading ? (
          <p className="py-8 text-center text-slate-500">
            Carregando listas...
          </p>
        ) : lists.length === 0 ? (
          <p className="py-8 text-center text-slate-500">
            Nenhuma lista cadastrada.
          </p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left">
              <thead>
                <tr className="border-b text-sm text-slate-500">
                  <th className="py-3 font-medium">ID</th>
                  <th className="py-3 font-medium">Título</th>
                  <th className="py-3 font-medium">Descrição</th>
                  <th className="py-3 font-medium">Criada em</th>
                  <th className="py-3 text-right font-medium">Ações</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {lists.map((list) => (
                  <tr key={list.id} id={`list-${list.id}`}>
                    <td className="py-3 text-sm text-slate-600">{list.id}</td>
                    <td className="py-3 font-medium text-slate-800">
                      {list.title}
                    </td>
                    <td className="py-3 text-sm text-slate-600">
                      {list.description || "—"}
                    </td>
                    <td className="py-3 text-sm text-slate-600">
                      {list.created_at || "—"}
                    </td>
                    <td className="py-3">
                      <div className="flex justify-end gap-3">
                        <button
                          className="inline-flex items-center gap-1 text-sm font-medium text-blue-600 hover:text-blue-800"
                          onClick={() => setModal(list)}
                          aria-label={`Editar ${list.title}`}
                        >
                          <Pencil size={16} />
                          Editar
                        </button>
                        <button
                          className="inline-flex items-center gap-1 text-sm font-medium text-red-600 hover:text-red-800"
                          onClick={() => remove(list)}
                          aria-label={`Excluir ${list.title}`}
                        >
                          <Trash2 size={16} />
                          Excluir
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      {modal && (
        <Modal
          title={modal === "create" ? "Nova lista" : "Editar lista"}
          onClose={() => setModal(null)}
        >
          <ListForm
            list={modal === "create" ? null : modal}
            onSubmit={submit}
            onCancel={() => setModal(null)}
          />
        </Modal>
      )}
    </div>
  );
}
