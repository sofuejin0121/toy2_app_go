import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import PasswordInput from '../components/PasswordInput';
import { signUp } from '../api/client';
import { getErrorList } from '../api/errors';
import ErrorMessage from '../components/ErrorMessage';
import Layout from '../components/Layout';

type SignupFieldName = 'name' | 'email' | 'password' | 'password_confirmation';

const SIGNUP_FIELDS: { name: SignupFieldName; label: string; type: string }[] = [
  { name: 'name', label: '名前', type: 'text' },
  { name: 'email', label: 'メールアドレス', type: 'email' },
  { name: 'password', label: 'パスワード', type: 'password' },
  { name: 'password_confirmation', label: 'パスワード（確認）', type: 'password' },
];

/**
 * 新規登録（POST /users）。成功時は確認メール案内 UI、失敗時は getErrorList で複数エラー表示。
 * SIGNUP_FIELDS で name の型を固定し、フォーム値へのアクセスで as を不要にしている。
 */
export default function SignupPage() {
  const navigate = useNavigate();
  const [form, setForm] = useState({
    name: '',
    email: '',
    password: '',
    password_confirmation: '',
  });
  const [errors, setErrors] = useState<string[]>([]);
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setErrors([]);
    try {
      await signUp(form);
      setSuccess(true);
    } catch (err: unknown) {
      setErrors(getErrorList(err));
    }
    setLoading(false);
  };

  if (success) {
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
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <h2 className="text-xl font-bold text-gray-900 mb-2">登録完了！</h2>
            <p className="text-gray-600 text-sm mb-4">
              確認メールを送信しました。メール内のリンクをクリックしてアカウントを有効化してください。
            </p>
            <button
              onClick={() => navigate('/login')}
              className="text-blue-600 hover:underline text-sm"
            >
              ログインページへ
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
          <h1 className="text-2xl font-bold text-gray-900 mb-6">新規登録</h1>

          {errors.length > 0 && <ErrorMessage messages={errors} />}

          <form onSubmit={handleSubmit} className="space-y-4">
            {SIGNUP_FIELDS.map((field) => (
              <div key={field.name}>
                <label
                  htmlFor={`signup-${field.name}`}
                  className="block text-sm font-medium text-gray-700 mb-1"
                >
                  {field.label}
                </label>
                {field.type === 'password' ? (
                  <PasswordInput
                    id={`signup-${field.name}`}
                    name={field.name}
                    value={form[field.name]}
                    onChange={handleChange}
                    required
                    autoComplete="new-password"
                    helperText={
                      field.name === 'password'
                        ? '右のボタンで入力内容の表示／非表示を切り替えられます。'
                        : '右のボタンで表示／非表示を切り替えられます。'
                    }
                  />
                ) : (
                  <input
                    id={`signup-${field.name}`}
                    type={field.type}
                    name={field.name}
                    value={form[field.name]}
                    onChange={handleChange}
                    required
                    autoComplete={field.name === 'email' ? 'email' : undefined}
                    className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                )}
              </div>
            ))}
            <button
              type="submit"
              disabled={loading}
              className="w-full bg-blue-600 text-white py-2 rounded-lg text-sm font-medium hover:bg-blue-700 disabled:opacity-50"
            >
              {loading ? '登録中...' : 'アカウントを作成'}
            </button>
          </form>

          <p className="mt-4 text-center text-sm text-gray-500">
            すでにアカウントをお持ちの方は{' '}
            <Link to="/login" className="text-blue-600 hover:underline">
              ログイン
            </Link>
          </p>
        </div>
      </div>
    </Layout>
  );
}
