import { useSetAtom } from 'jotai';
import { useEffect, useState } from 'react';
import { useNavigate, useParams, useSearchParams } from 'react-router-dom';
import { checkPasswordResetToken, resetPassword } from '../api/client';
import { getErrorList } from '../api/errors';
import ErrorMessage from '../components/ErrorMessage';
import Layout from '../components/Layout';
import PasswordInput from '../components/PasswordInput';
import { useDocumentTitle } from '../hooks/useDocumentTitle';
import { authBootstrapEpochAtom, currentUserAtom } from '../store/auth';

/**
 * メール内リンクから遷移するパスワード再設定（トークン・email を検証してから PATCH）。
 * 成功時はログイン済みとして currentUserAtom を更新しプロフィールへ。
 */
export default function PasswordResetEditPage() {
  useDocumentTitle('パスワード再設定');
  const { token } = useParams<{ token: string }>();
  const [searchParams] = useSearchParams();
  const email = searchParams.get('email') || '';
  const setCurrentUser = useSetAtom(currentUserAtom);
  const bumpAuthEpoch = useSetAtom(authBootstrapEpochAtom);
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
      bumpAuthEpoch((n) => n + 1);
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
              type="button"
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

          {errors.length > 0 && <ErrorMessage messages={errors} />}

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label
                htmlFor="reset-password"
                className="block text-sm font-medium text-gray-700 mb-1"
              >
                新しいパスワード
              </label>
              <PasswordInput
                id="reset-password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                autoComplete="new-password"
                helperText="右のボタンで入力内容の表示／非表示を切り替えられます。"
              />
            </div>
            <div>
              <label
                htmlFor="reset-password-confirm"
                className="block text-sm font-medium text-gray-700 mb-1"
              >
                パスワード（確認）
              </label>
              <PasswordInput
                id="reset-password-confirm"
                value={passwordConfirmation}
                onChange={(e) => setPasswordConfirmation(e.target.value)}
                required
                autoComplete="new-password"
                helperText="右のボタンで表示／非表示を切り替えられます。"
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
