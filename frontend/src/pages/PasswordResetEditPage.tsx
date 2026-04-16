import { useSetAtom } from 'jotai';
import { useEffect, useState } from 'react';
import { useNavigate, useParams, useSearchParams } from 'react-router-dom';
import { checkPasswordResetToken, resetPassword } from '../api/client';
import { getErrorList } from '../api/errors';
import Layout from '../components/Layout';
import { currentUserAtom } from '../store/auth';

export default function PasswordResetEditPage() {
  const { token } = useParams<{ token: string }>();
  const [searchParams] = useSearchParams();
  const email = searchParams.get('email') || '';
  const setCurrentUser = useSetAtom(currentUserAtom);
  const navigate = useNavigate();
  const [valid, setValid] = useState<boolean | null>(null);
  const [password, setPassword] = useState('');
  const [passwordConfirmation, setPasswordConfirmation] = useState('');
  const [errors, setErrors] = useState<string[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!token || !email) {
      setValid(false);
      return;
    }
    checkPasswordResetToken(token, email)
      .then(() => setValid(true))
      .catch(() => setValid(false));
  }, [token, email]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token) return;
    setLoading(true);
    setErrors([]);
    try {
      const data = await resetPassword(token, {
        email,
        password,
        password_confirmation: passwordConfirmation,
      });
      setCurrentUser(data.user);
      navigate(`/users/${data.user.id}`);
    } catch (err: unknown) {
      setErrors(getErrorList(err));
    }
    setLoading(false);
  };

  if (valid === null) {
    return (
      <Layout>
        <div className="text-center py-10 text-gray-400">確認中...</div>
      </Layout>
    );
  }
  if (!valid) {
    return (
      <Layout>
        <div className="max-w-md mx-auto text-center py-12">
          <div className="bg-red-50 border border-red-200 rounded-xl p-8">
            <h2 className="text-xl font-bold text-gray-900 mb-2">無効なリンク</h2>
            <p className="text-gray-600 text-sm">このリンクは無効または期限切れです</p>
            <button
              onClick={() => navigate('/password_resets/new')}
              className="mt-4 text-blue-600 hover:underline text-sm"
            >
              もう一度送る
            </button>
          </div>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="max-w-md mx-auto">
        <div className="bg-white rounded-xl border border-gray-200 p-8">
          <h1 className="text-2xl font-bold text-gray-900 mb-6">パスワードを再設定</h1>

          {errors.length > 0 && (
            <div className="mb-4 p-3 bg-red-50 border border-red-200 text-red-700 text-sm rounded-lg">
              <ul className="list-disc list-inside space-y-1">
                {errors.map((e, i) => (
                  <li key={i}>{e}</li>
                ))}
              </ul>
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                新しいパスワード
              </label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                パスワード（確認）
              </label>
              <input
                type="password"
                value={passwordConfirmation}
                onChange={(e) => setPasswordConfirmation(e.target.value)}
                required
                className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <button
              type="submit"
              disabled={loading}
              className="w-full bg-blue-600 text-white py-2 rounded-lg text-sm font-medium hover:bg-blue-700 disabled:opacity-50"
            >
              {loading ? '更新中...' : 'パスワードを更新'}
            </button>
          </form>
        </div>
      </div>
    </Layout>
  );
}
