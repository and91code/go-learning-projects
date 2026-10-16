import { Navigate, Route, Routes } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext'
import AppLayout from './components/layout/AppLayout'
import RequireAuth from './components/routing/RequireAuth'
import LoginView from './views/LoginView'
import RegisterView from './views/RegisterView'
import MainView from './views/MainView'
import ProfileView from './views/ProfileView'
import ListsView from './views/ListsView'
import TagsView from './views/TagsView'

export default function App() {
  return (
    <AuthProvider>
      <Routes>
        <Route element={<AppLayout />}>
          <Route path="/login" element={<LoginView />} />
          <Route path="/register" element={<RegisterView />} />
          <Route element={<RequireAuth />}>
            <Route path="/main" element={<MainView />} />
            <Route path="/my-profile" element={<ProfileView />} />
            <Route path="/lists" element={<ListsView />} />
            <Route path="/tags" element={<TagsView />} />
          </Route>
        </Route>
        <Route path="*" element={<Navigate to="/main" replace />} />
      </Routes>
    </AuthProvider>
  )
}
