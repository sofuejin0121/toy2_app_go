import { useAtom } from 'jotai';
import { Navigate, Route, Routes, useParams } from 'react-router-dom';
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
import LoadingSpinner from './components/LoadingSpinner';
import { currentUserAtom } from './store/auth';

/**
 * ログイン済みユーザーのみ通すラッパー。
 * - currentUser === undefined … まだ「自分が誰か」取得中 → スピナー
 * - currentUser === null … 未ログイン → /login へ
 * - User オブジェクト … 子のページを表示
 */
function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const [currentUser] = useAtom(currentUserAtom);
  if (currentUser === undefined)
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner />
      </div>
    );
  if (!currentUser) return <Navigate to="/login" replace />;
  return <>{children}</>;
}

/**
 * 「URL の :id とログイン中ユーザーが同じ人」だけ通すラッパー。
 * プロフィール編集や自分のブックマークのように、他人が URL を直接叩いても見せたくない画面で使います。
 * 認可チェックを各ページや各フックに散らさず、ルート定義に集約するのが目的です。
 */
function OwnerRoute({ children }: { children: React.ReactNode }) {
  const [currentUser] = useAtom(currentUserAtom);
  const { id } = useParams<{ id: string }>();
  if (currentUser === undefined)
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner />
      </div>
    );
  if (!currentUser) return <Navigate to="/login" replace />;
  if (String(currentUser.id) !== id) return <Navigate to="/" replace />;
  return <>{children}</>;
}

/**
 * 管理者（admin フラグが true）だけ通すラッパー。
 * AdminPage 内でまた navigate する必要がなくなります。
 */
function AdminRoute({ children }: { children: React.ReactNode }) {
  const [currentUser] = useAtom(currentUserAtom);
  if (currentUser === undefined)
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner />
      </div>
    );
  if (!currentUser) return <Navigate to="/login" replace />;
  if (!currentUser.admin) return <Navigate to="/" replace />;
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
          <OwnerRoute>
            <UserEditPage />
          </OwnerRoute>
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
          <OwnerRoute>
            <BookmarksPage />
          </OwnerRoute>
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
          <AdminRoute>
            <AdminPage />
          </AdminRoute>
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
