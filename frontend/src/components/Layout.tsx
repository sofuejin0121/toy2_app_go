import { useAtom } from 'jotai';
import { type ReactNode, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { logout } from '../api/client';
import { currentUserAtom } from '../store/auth';

interface AlertState {
  type: 'success' | 'error' | 'info' | 'warning';
  message: string;
}

interface LayoutProps {
  children: ReactNode;
  alert?: AlertState | null;
}

export default function Layout({ children, alert }: LayoutProps) {
  const [currentUser, setCurrentUser] = useAtom(currentUserAtom);
  const navigate = useNavigate();
  const [menuOpen, setMenuOpen] = useState(false);

  const handleLogout = async () => {
    await logout();
    setCurrentUser(null);
    navigate('/login');
  };

  const alertColors = {
    success: 'bg-green-50 border-green-400 text-green-800',
    error: 'bg-red-50 border-red-400 text-red-800',
    info: 'bg-blue-50 border-blue-400 text-blue-800',
    warning: 'bg-yellow-50 border-yellow-400 text-yellow-800',
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white border-b border-gray-200 sticky top-0 z-50">
        <div className="max-w-5xl mx-auto px-4 flex items-center justify-between h-14">
          <Link to="/" className="text-xl font-bold text-blue-600 hover:text-blue-700">
            MicroApp
          </Link>

          {/* Desktop Nav */}
          <div className="hidden md:flex items-center gap-4">
            {currentUser ? (
              <>
                <Link to="/" className="text-gray-600 hover:text-gray-900 text-sm font-medium">
                  Home
                </Link>
                <Link
                  to={`/users/${currentUser.id}`}
                  className="flex items-center gap-2 text-sm text-gray-700 hover:text-gray-900"
                >
                  <img
                    src={currentUser.avatar_url}
                    alt={currentUser.name}
                    className="w-7 h-7 rounded-full"
                  />
                  <span className="font-medium">{currentUser.name}</span>
                </Link>
                <Link to="/notifications" className="text-gray-600 hover:text-gray-900 text-sm">
                  通知
                </Link>
                <Link to="/settings" className="text-gray-600 hover:text-gray-900 text-sm">
                  設定
                </Link>
                {currentUser.admin && (
                  <Link
                    to="/admin"
                    className="text-orange-600 hover:text-orange-700 text-sm font-medium"
                  >
                    Admin
                  </Link>
                )}
                <button onClick={handleLogout} className="text-sm text-gray-500 hover:text-red-600">
                  ログアウト
                </button>
              </>
            ) : (
              <>
                <Link to="/login" className="text-sm text-gray-600 hover:text-gray-900">
                  ログイン
                </Link>
                <Link
                  to="/signup"
                  className="text-sm bg-blue-600 text-white px-4 py-1.5 rounded-full hover:bg-blue-700"
                >
                  新規登録
                </Link>
              </>
            )}
          </div>

          {/* Mobile Hamburger */}
          <button className="md:hidden text-gray-600" onClick={() => setMenuOpen(!menuOpen)}>
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              {menuOpen ? (
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M6 18L18 6M6 6l12 12"
                />
              ) : (
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M4 6h16M4 12h16M4 18h16"
                />
              )}
            </svg>
          </button>
        </div>

        {/* Mobile Menu */}
        {menuOpen && (
          <div className="md:hidden border-t border-gray-100 px-4 py-3 space-y-2 bg-white">
            {currentUser ? (
              <>
                <Link
                  to="/"
                  onClick={() => setMenuOpen(false)}
                  className="block text-gray-700 py-1"
                >
                  Home
                </Link>
                <Link
                  to={`/users/${currentUser.id}`}
                  onClick={() => setMenuOpen(false)}
                  className="block text-gray-700 py-1"
                >
                  プロフィール
                </Link>
                <Link
                  to="/notifications"
                  onClick={() => setMenuOpen(false)}
                  className="block text-gray-700 py-1"
                >
                  通知
                </Link>
                <Link
                  to="/settings"
                  onClick={() => setMenuOpen(false)}
                  className="block text-gray-700 py-1"
                >
                  設定
                </Link>
                {currentUser.admin && (
                  <Link
                    to="/admin"
                    onClick={() => setMenuOpen(false)}
                    className="block text-orange-600 py-1"
                  >
                    Admin
                  </Link>
                )}
                <button onClick={handleLogout} className="block text-red-500 py-1 w-full text-left">
                  ログアウト
                </button>
              </>
            ) : (
              <>
                <Link
                  to="/login"
                  onClick={() => setMenuOpen(false)}
                  className="block text-gray-700 py-1"
                >
                  ログイン
                </Link>
                <Link
                  to="/signup"
                  onClick={() => setMenuOpen(false)}
                  className="block text-blue-600 py-1"
                >
                  新規登録
                </Link>
              </>
            )}
          </div>
        )}
      </nav>

      <main className="max-w-5xl mx-auto px-4 py-6">
        {alert && (
          <div className={`mb-4 px-4 py-3 rounded-lg border ${alertColors[alert.type]}`}>
            {alert.message}
          </div>
        )}
        {children}
      </main>
    </div>
  );
}
