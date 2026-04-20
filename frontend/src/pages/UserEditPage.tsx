import { useAtom } from 'jotai';
import { useEffect, useRef, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import useSWR from 'swr';
import { getUser, updateUser } from '../api/client';
import { getErrorList } from '../api/errors';
import ErrorMessage from '../components/ErrorMessage';
import Layout from '../components/Layout';
import LoadingSpinner from '../components/LoadingSpinner';
import PasswordInput from '../components/PasswordInput';
import { useDocumentTitle } from '../hooks/useDocumentTitle';
import { currentUserAtom } from '../store/auth';

/**
 * プロフィール編集（OwnerRoute）。認可は App.tsx の OwnerRoute が担当。
 *
 * 初期表示用に useSWR で getUser を呼びます（編集に必要なのは user 部分だけだが、既存 API を流用）。
 * SWR がバックグラウンドで再検証しても、入力途中のフォームを上書きしないよう ref で「この id には初回だけ反映」と制御しています。
 */
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
  // SWR が再検証して data が更新されても、編集中の入力を上書きしないためのフラグ。
  // 「どのユーザー id に対してサーバー値をフォームへ流し込んだか」を保持し、同一 id では一度だけ setForm する。
  const filledForId = useRef<string | undefined>(undefined);

  // 第 1 ページのプロフィールを取ればフォーム初期値に十分（microposts は未使用）
  const { data, isLoading: pageLoading } = useSWR(
    id ? `user-edit-${id}` : null,
    () => getUser(Number(id), 1),
  );

  useDocumentTitle(
    data?.user ? `${data.user.name} · プロフィール編集` : 'プロフィール編集',
  );

  useEffect(() => {
    if (id && filledForId.current !== undefined && filledForId.current !== id) {
      filledForId.current = undefined;
    }
  }, [id]);

  useEffect(() => {
    if (!data || !id) return;
    if (filledForId.current === id) return;
    filledForId.current = id;
    setForm((prev) => ({
      ...prev,
      name: data.user.name,
      email: data.user.email,
      bio: data.user.bio ?? '',
    }));
  }, [id, data]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setForm((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!id) return;
    setSubmitting(true);
    setErrors([]);
    try {
      const updated = await updateUser(Number(id), form);
      if (currentUser) setCurrentUser({ ...currentUser, ...updated });
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

          {errors.length > 0 && <ErrorMessage messages={errors} />}

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
                  <label
                    htmlFor="edit-password"
                    className="block text-sm font-medium text-gray-700 mb-1"
                  >
                    新しいパスワード
                  </label>
                  <PasswordInput
                    id="edit-password"
                    name="password"
                    value={form.password}
                    onChange={handleChange}
                    autoComplete="new-password"
                    helperText="右のボタンで入力内容の表示／非表示を切り替えられます。"
                  />
                </div>
                <div>
                  <label
                    htmlFor="edit-password-confirm"
                    className="block text-sm font-medium text-gray-700 mb-1"
                  >
                    新しいパスワード（確認）
                  </label>
                  <PasswordInput
                    id="edit-password-confirm"
                    name="password_confirmation"
                    value={form.password_confirmation}
                    onChange={handleChange}
                    autoComplete="new-password"
                    helperText="右のボタンで表示／非表示を切り替えられます。"
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
