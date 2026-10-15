import { useState } from "react";
import { LogIn, UserPlus } from "lucide-react";
import { Link, Navigate, useLocation } from "react-router-dom";
import Button from "../components/ui/Button";
import Input from "../components/ui/Input";
import { useAuth } from "../context/AuthContext";
import { login } from "../services/authService";

export default function LoginView() {
  const { user, signIn } = useAuth();
  const location = useLocation();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);
  async function submit(event) {
    event.preventDefault();
    setBusy(true);
    setError("");
    try {
      await signIn(login, { email, password });
    } catch (err) {
      setError(err.message);
    } finally {
      setBusy(false);
    }
  }
  if (user) {
    return <Navigate to={location.state?.from?.pathname || "/main"} replace />;
  }
  return (
    <AuthLayout title="Bem-vindo de volta">
      <form onSubmit={submit} className="grid gap-4">
        <Input
          label="E-mail"
          type="email"
          autoComplete="email"
          value={email}
          onChange={(event) => setEmail(event.target.value)}
          required
        />
        <Input
          label="Senha"
          type="password"
          autoComplete="current-password"
          value={password}
          onChange={(event) => setPassword(event.target.value)}
          required
        />
        {error && (
          <p className="text-sm text-red-600" role="alert">
            {error}
          </p>
        )}
        <Button disabled={busy}><LogIn size={18} className="mr-2" />{busy ? "Entrando…" : "Entrar"}</Button>
      </form>
    </AuthLayout>
  );
}

function AuthLayout({ title, children }) {
  return (
    <div className="grid min-h-full place-items-center p-4">
      <section className="w-full max-w-md rounded-2xl bg-white p-8 shadow-sm">
        <h1 className="mb-2 text-3xl font-bold text-slate-900">TaskFlow</h1>
        <h2 className="mb-6 text-lg text-slate-600">{title}</h2>
        {children}
        <p className="mt-5 text-center text-sm text-slate-600">
          Ainda não tem conta?{" "}
          <Link className="font-semibold text-indigo-600" to="/register">
            <UserPlus size={18} className="mr-1 inline" /> Cadastre-se
          </Link>
        </p>
      </section>
    </div>
  );
}
