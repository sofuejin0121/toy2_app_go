import { useAtom } from 'jotai';
import { Navigate, Route, Routes } from 'react-router-dom';
import AccountActivationPage from './pages/AccountActivationPage';
import AdminPage from './pages/AdminPage';
import BookmarksPage from './pages/BookmarksPage';
import FollowListPage from './pages/FollowListPage';
import HomePage from './pages/HomePage';
import LikesPage from './pages/LikesPage';
import LoginPage from './pages/LoginPage';
import MicropostPage from './pages/MicropostPage';
import NotificationsPage from './pages/NotificationsPage';
import PasswordResetEditPage from './pages/PasswordResetEditPage';
import PasswordResetNewPage from './pages/PasswordResetNewPage';
import SettingsPage from './pages/SettingsPage';
import SignupPage from './pages/SignupPage';
import UserEditPage from './pages/UserEditPage';
import UserListPage from './pages/UserListPage';
import UserShowPage from './pages/UserShowPage';
import { currentUserAtom } from './store/auth';

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  // undefined = /me 取得中, null = 未ログイン, User = ログイン済み
  const [currentUser] = useAtom(currentUserAtom);
  if (currentUser === undefined)
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin w-8 h-8 border-2 border-blue-600 border-t-transparent rounded-full" />
      </div>
    );
  if (!currentUser) return <Navigate to="/login" replace />;
  return <>{children}</>;
}

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/login" element={<LoginPage />} />
      <Route path="/signup" element={<SignupPage />} />
      <Route
        path="/users"
        element={
          <ProtectedRoute>
            <UserListPage />
          </ProtectedRoute>
        }
      />
      <Route path="/users/:id" element={<UserShowPage />} />
      <Route
        path="/users/:id/edit"
        element={
          <ProtectedRoute>
            <UserEditPage />
          </ProtectedRoute>
        }
      />
      <Route
        path="/users/:id/following"
        element={
          <ProtectedRoute>
            <FollowListPage mode="following" />
          </ProtectedRoute>
        }
      />
      <Route
        path="/users/:id/followers"
        element={
          <ProtectedRoute>
            <FollowListPage mode="followers" />
          </ProtectedRoute>
        }
      />
      <Route
        path="/users/:id/likes"
        element={
          <ProtectedRoute>
            <LikesPage />
          </ProtectedRoute>
        }
      />
      <Route
        path="/users/:id/bookmarks"
        element={
          <ProtectedRoute>
            <BookmarksPage />
          </ProtectedRoute>
        }
      />
      <Route path="/microposts/:id" element={<MicropostPage />} />
      <Route
        path="/notifications"
        element={
          <ProtectedRoute>
            <NotificationsPage />
          </ProtectedRoute>
        }
      />
      <Route
        path="/admin"
        element={
          <ProtectedRoute>
            <AdminPage />
          </ProtectedRoute>
        }
      />
      <Route
        path="/settings"
        element={
          <ProtectedRoute>
            <SettingsPage />
          </ProtectedRoute>
        }
      />
      <Route path="/account_activations/:token/edit" element={<AccountActivationPage />} />
      <Route path="/password_resets/new" element={<PasswordResetNewPage />} />
      <Route path="/password_resets/:token/edit" element={<PasswordResetEditPage />} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
