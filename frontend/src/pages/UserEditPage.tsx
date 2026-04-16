import { useAtom } from 'jotai';
import { useEffect, useRef, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { getUser, updateUser } from '../api/client';
import { getErrorList } from '../api/errors';
import Layout from '../components/Layout';
import LoadingSpinner from '../components/LoadingSpinner';
import { currentUserAtom } from '../store/auth';

export default function UserEditPage() {
  const { id } = useParams<{ id: string }>();
  const [currentUser, setCurrentUser] = useAtom(currentUserAtom);
  const navigate = useNavigate();
  const [form, setForm] = useState({
    name: '',
    email: '',
    bio: '',
    password: '',
    password_confirmation: '',
  });
  const [errors, setErrors] = useState<string[]>([]);
  const [submitting, setSubmitting] = useState(false);
  const [pageLoading, setPageLoading] = useState(true);

  // 自分以外のユーザーの編集ページへのアクセスをリダイレクト
  // render 中に navigate を呼ぶと React の警告が出るため useEffect に寄せる
  // redirectedRef はリダイレクトを1回だけ実行するためのフラグ
  const redirectedRef = useRef(false);
  useEffect(() => {
    if (
      !redirectedRef.current &&
      currentUser !== undefined &&
      (!currentUser || currentUser.id !== Number(id))
    ) {
      redirectedRef.current = true;
      navigate('/');
    }
  }, [currentUser, id, navigate]);

  // 編集対象ユーザーの現在値をフォームに反映する
  useEffect(() => {
    if (!id) return;

    async function loadUser() {
      try {
        const data = await getUser(Number(id));
        setForm((prev) => ({
          ...prev,
          name: data.user.name,
          email: data.user.email,
          bio: data.user.bio ?? '',
        }));
      } finally {
        setPageLoading(false);
      }
    }

    loadUser();
  }, [id]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  // フォーム送信：プロフィールを更新して詳細ページへ遷移
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setErrors([]);
    try {
      const updated = await updateUser(Number(id), form);
      // Jotai atom を更新してヘッダーの表示名などを即時反映する
      setCurrentUser({ ...currentUser, ...updated });
      navigate(`/users/${id}`);
    } catch (err: unknown) {
      setErrors(getErrorList(err));
    }
    setSubmitting(false);
  };

  if (pageLoading)
    return (
      <Layout>
        <LoadingSpinner />
      </Layout>
    );

  return (
    <Layout>
      <div className="max-w-lg mx-auto">
        <div className="bg-white rounded-xl border border-gray-200 p-8">
          <h1 className="text-2xl font-bold text-gray-900 mb-6">プロフィール編集</h1>

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
              <label className="block text-sm font-medium text-gray-700 mb-1">名前</label>
              <input
                type="text"
                name="name"
                value={form.name}
                onChange={handleChange}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">メールアドレス</label>
              <input
                type="email"
                name="email"
                value={form.email}
                onChange={handleChange}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">自己紹介</label>
              <textarea
                name="bio"
                value={form.bio}
                onChange={handleChange}
                rows={3}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div className="border-t border-gray-100 pt-4">
              <p className="text-xs text-gray-500 mb-3">
                パスワードを変更する場合のみ入力してください
              </p>
              <div className="space-y-3">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    新しいパスワード
                  </label>
                  <input
                    type="password"
                    name="password"
                    value={form.password}
                    onChange={handleChange}
                    className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    新しいパスワード（確認）
                  </label>
                  <input
                    type="password"
                    name="password_confirmation"
                    value={form.password_confirmation}
                    onChange={handleChange}
                    className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
              </div>
            </div>
            <button
              type="submit"
              disabled={submitting}
              className="w-full bg-blue-600 text-white py-2 rounded-lg text-sm font-medium hover:bg-blue-700 disabled:opacity-50"
            >
              {submitting ? '保存中...' : '変更を保存'}
            </button>
          </form>
        </div>
      </div>
    </Layout>
  );
}
