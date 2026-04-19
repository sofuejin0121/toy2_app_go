import { useSetAtom } from "jotai";
import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { login } from "../api/client";
import { getErrorMessage } from "../api/errors";
import ErrorMessage from "../components/ErrorMessage";
import Layout from "../components/Layout";
import PasswordInput from "../components/PasswordInput";
import { authBootstrapEpochAtom, currentUserAtom } from "../store/auth";

/**
 * ログイン（POST /login）。成功時に currentUserAtom を更新し、プロフィールへ navigate。
 */
export default function LoginPage() {
  const setCurrentUser = useSetAtom(currentUserAtom);
  const bumpAuthEpoch = useSetAtom(authBootstrapEpochAtom);
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [remember, setRemember] = useState(false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");
    try {
      const user = await login(email, password, remember);
      // ログイン前に飛んでいた getMe の 401 が遅れて届いても currentUser を消さないよう epoch を先に進める
      bumpAuthEpoch((n) => n + 1);
      setCurrentUser(user);
      navigate(`/`);
    } catch (err: unknown) {
      setError(getErrorMessage(err, "ログインに失敗しました"));
    }
    setLoading(false);
  };

  return (
    <Layout>
      <div className="max-w-md mx-auto">
        <div className="bg-white rounded-xl border border-gray-200 p-8">
          <h1 className="text-2xl font-bold text-gray-900 mb-6">ログイン</h1>

          {error && <ErrorMessage message={error} />}

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label
                htmlFor="login-email"
                className="block text-sm font-medium text-gray-700 mb-1"
              >
                メールアドレス
              </label>
              <input
                id="login-email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                autoComplete="email"
                className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label
                htmlFor="login-password"
                className="block text-sm font-medium text-gray-700 mb-1"
              >
                パスワード
              </label>
              <PasswordInput
                id="login-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                autoComplete="current-password"
                helperText="右のボタンで入力したパスワードの表示／非表示を切り替えられます。"
              />
            </div>
            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                id="remember"
                checked={remember}
                onChange={(e) => setRemember(e.target.checked)}
                className="rounded"
              />
              <label htmlFor="remember" className="text-sm text-gray-600">
                ログイン状態を保持
              </label>
            </div>
            <button
              type="submit"
              disabled={loading}
              className="w-full bg-blue-600 text-white py-2 rounded-lg text-sm font-medium hover:bg-blue-700 disabled:opacity-50"
            >
              {loading ? "ログイン中..." : "ログイン"}
            </button>
          </form>

          <div className="mt-4 text-center text-sm text-gray-500 space-y-1">
            <p>
              <Link
                to="/password_resets/new"
                className="text-blue-600 hover:underline"
              >
                パスワードをお忘れですか？
              </Link>
            </p>
            <p>
              アカウントをお持ちでない方は{" "}
              <Link to="/signup" className="text-blue-600 hover:underline">
                新規登録
              </Link>
            </p>
          </div>
        </div>
      </div>
    </Layout>
  );
}
