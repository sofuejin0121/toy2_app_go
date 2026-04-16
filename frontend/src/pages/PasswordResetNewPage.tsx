import { useState } from 'react';
import { Link } from 'react-router-dom';
import { requestPasswordReset } from '../api/client';
import { getErrorMessage } from '../api/errors';
import Layout from '../components/Layout';

export default function PasswordResetNewPage() {
  const [email, setEmail] = useState('');
  const [loading, setLoading] = useState(false);
  const [sent, setSent] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    try {
      await requestPasswordReset(email);
      setSent(true);
    } catch (err: unknown) {
      setError(getErrorMessage(err));
    }
    setLoading(false);
  };

  if (sent) {
    return (
      <Layout>
        <div className="max-w-md mx-auto text-center py-12">
          <div className="bg-green-50 border border-green-200 rounded-xl p-8">
            <svg
              className="w-12 h-12 text-green-500 mx-auto mb-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
              />
            </svg>
            <h2 className="text-xl font-bold text-gray-900 mb-2">メールを送信しました</h2>
            <p className="text-gray-600 text-sm">
              パスワード再設定のメールを送信しました。メールをご確認ください。
            </p>
            <Link to="/login" className="mt-4 inline-block text-blue-600 hover:underline text-sm">
              ログインページへ
            </Link>
          </div>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="max-w-md mx-auto">
        <div className="bg-white rounded-xl border border-gray-200 p-8">
          <h1 className="text-2xl font-bold text-gray-900 mb-2">パスワードを忘れた方へ</h1>
          <p className="text-sm text-gray-500 mb-6">登録済みのメールアドレスを入力してください</p>

          {error && (
            <div className="mb-4 p-3 bg-red-50 border border-red-200 text-red-700 text-sm rounded-lg">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">メールアドレス</label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <button
              type="submit"
              disabled={loading}
              className="w-full bg-blue-600 text-white py-2 rounded-lg text-sm font-medium hover:bg-blue-700 disabled:opacity-50"
            >
              {loading ? '送信中...' : 'リセットメールを送信'}
            </button>
          </form>
          <p className="mt-4 text-center text-sm">
            <Link to="/login" className="text-blue-600 hover:underline">
              ログインページへ戻る
            </Link>
          </p>
        </div>
      </div>
    </Layout>
  );
}
