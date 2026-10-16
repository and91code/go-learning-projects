import { Folder, Pencil, PlusCircle, Trash2 } from "lucide-react";
import { useEffect, useState } from "react";
import { useSearchParams } from "react-router-dom";
import TagForm from "../components/forms/TagForm";
import Button from "../components/ui/Button";
import Modal from "../components/ui/Modal";
import { useAuth } from "../context/AuthContext";
import { getLists } from "../services/listService";
import {
  createTag,
  deleteTag,
  getTagById,
  getTagsByList,
  updateTag,
} from "../services/tagService";

export default function TagsView() {
  const { resourceVersions } = useAuth();
  const [searchParams] = useSearchParams();
  const [lists, setLists] = useState([]);
  const [listId, setListId] = useState("");
  const [tags, setTags] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [modal, setModal] = useState(null);

  useEffect(() => {
    let ignore = false;
    async function loadLists() {
      try {
        const result = await getLists();
        if (ignore) return;
        const items = Array.isArray(result) ? result : result?.data || [];
        setLists(items);
        const tagId = searchParams.get("tag_id");
        if (tagId) {
          const tag = await getTagById(tagId);
          if (ignore) return;
          setListId(String(tag.list_id));
        } else {
          setListId(items[0] ? String(items[0].id) : "");
        }
      } catch (err) {
        if (!ignore) setError(err.message);
      }
    }
    loadLists();
    return () => {
      ignore = true;
    };
  }, [resourceVersions.list, searchParams]);

  useEffect(() => {
    if (!listId) return;
    let ignore = false;
    getTagsByList(listId)
      .then((result) => {
        if (!ignore)
          setTags(Array.isArray(result) ? result : result?.data || []);
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
  }, [listId, resourceVersions.tag]);

  useEffect(() => {
    const tagId = searchParams.get("tag_id");
    if (!tagId || !tags.some((tag) => String(tag.id) === tagId)) return;
    document.getElementById(`tag-${tagId}`)?.scrollIntoView({ behavior: "smooth", block: "center" });
  }, [tags, searchParams]);

  function changeList(value) {
    setListId(value);
    setTags([]);
    setLoading(true);
    setError("");
  }

  async function submit(data) {
    if (modal === "create") {
      const result = await createTag(data);
      if (result?.status !== "accepted")
        setTags((current) => [...current, result?.data || result]);
    } else {
      const result = await updateTag(modal.id, data);
      const updated = result?.data || result;
      if (result?.status !== "accepted")
        setTags((current) => current.map((tag) => tag.id === modal.id ? { ...tag, ...updated } : tag));
    }
    setModal(null);
  }

  async function remove(tag) {
    if (
      !window.confirm(
        `Excluir a tag "${tag.name}"? Ela será removida das tarefas que a utilizam.`,
      )
    )
      return;
    try {
      const result = await deleteTag(tag.id);
      if (result?.status !== "accepted")
        setTags((current) => current.filter((item) => item.id !== tag.id));
    } catch (err) {
      setError(err.message);
    }
  }

  return (
    <div className="grid gap-6">
      <section className="flex flex-wrap items-end justify-between gap-4 rounded-2xl bg-white p-6 shadow-sm">
        <div className="flex flex-wrap items-end gap-4">
          <label className="grid min-w-56 gap-1.5 text-sm font-medium text-slate-700">
            Lista
            <span className="flex items-center gap-2">
              <Folder size={18} className="shrink-0 text-slate-500" />
              <select
                className="w-full rounded-lg border border-slate-300 px-3 py-2"
                value={listId}
                onChange={(event) => changeList(event.target.value)}
                disabled={!lists.length}
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
        </div>
        <Button onClick={() => setModal("create")} disabled={!listId}>
          <PlusCircle size={18} className="mr-2" />
          Adicionar
        </Button>
      </section>

      <section className="rounded-2xl bg-white p-6 shadow-sm">
        <h1 className="mb-4 text-2xl font-bold">Tags</h1>
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
        {!listId ? (
          <p className="py-8 text-center text-slate-500">
            Crie uma lista para cadastrar tags.
          </p>
        ) : loading ? (
          <p className="py-8 text-center text-slate-500">Carregando tags...</p>
        ) : tags.length === 0 ? (
          <p className="py-8 text-center text-slate-500">
            Nenhuma tag cadastrada nesta lista.
          </p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left">
              <thead>
                <tr className="border-b text-sm text-slate-500">
                  <th className="py-3 font-medium">ID</th>
                  <th className="py-3 font-medium">Nome</th>
                  <th className="py-3 font-medium">Cor</th>
                  <th className="py-3 text-right font-medium">Ações</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {tags.map((tag) => (
                  <tr key={tag.id} id={`tag-${tag.id}`}>
                    <td className="py-3 text-sm text-slate-600">{tag.id}</td>
                    <td className="py-3">
                      <span className="inline-flex items-center gap-2 font-medium text-slate-800">
                        <span
                          className="h-3 w-3 shrink-0 rounded-full"
                          style={{ backgroundColor: tag.color }}
                          aria-hidden="true"
                        />
                        {tag.name}
                      </span>
                    </td>
                    <td className="py-3 font-mono text-sm text-slate-600">
                      {tag.color}
                    </td>
                    <td className="py-3">
                      <div className="flex justify-end gap-3">
                        <button
                          className="inline-flex items-center gap-1 text-sm font-medium text-blue-600 hover:text-blue-800"
                          onClick={() => setModal(tag)}
                          aria-label={`Editar ${tag.name}`}
                        >
                          <Pencil size={16} />
                          Editar
                        </button>
                        <button
                          className="inline-flex items-center gap-1 text-sm font-medium text-red-600 hover:text-red-800"
                          onClick={() => remove(tag)}
                          aria-label={`Excluir ${tag.name}`}
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

      {modal && listId && (
        <Modal
          title={modal === "create" ? "Nova tag" : "Editar tag"}
          onClose={() => setModal(null)}
        >
          <TagForm
            listId={listId}
            tag={modal === "create" ? null : modal}
            onSubmit={submit}
            onCancel={() => setModal(null)}
          />
        </Modal>
      )}
    </div>
  );
}
