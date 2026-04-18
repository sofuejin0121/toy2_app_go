import { useAtom } from 'jotai';
import { useState } from 'react';
import { deleteUser } from '../api/client';
import Layout from '../components/Layout';
import LoadingSpinner from '../components/LoadingSpinner';
import Pagination from '../components/Pagination';
import UserCard from '../components/UserCard';
import { useUserList } from '../hooks/useUserList';
import { currentUserAtom } from '../store/auth';
import { updateIfDefined } from '../utils/updateIfDefined';

export default function UserListPage() {
  const [currentUser] = useAtom(currentUserAtom);
  const [page, setPage] = useState(1);
  const [query, setQuery] = useState('');
  const [inputQuery, setInputQuery] = useState('');

  const { users, pagination, loading, mutate } = useUserList(page, query);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setPage(1);
    setQuery(inputQuery);
  };

  const handleDelete = async (userId: number) => {
    if (!window.confirm('このユーザーを削除しますか？')) return;
    try {
      await deleteUser(userId);
      // SWR のキャッシュから削除済みユーザーを除く（未取得なら prev は undefined のまま）
      mutate(
        (prev) => updateIfDefined(prev, (p) => ({ ...p, users: p.users.filter((u) => u.id !== userId) })),
        { revalidate: false },
      );
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <Layout>
      <div className="max-w-2xl mx-auto">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">ユーザー一覧</h1>

        <form onSubmit={handleSearch} className="mb-6 flex gap-2">
          <input
            type="text"
            value={inputQuery}
            onChange={(e) => setInputQuery(e.target.value)}
            placeholder="ユーザーを検索..."
            className="flex-1 border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <button
            type="submit"
            className="bg-blue-600 text-white px-4 py-2 rounded-lg text-sm hover:bg-blue-700"
          >
            検索
          </button>
          {query && (
            <button
              type="button"
              onClick={() => {
                setQuery('');
                setInputQuery('');
                setPage(1);
              }}
              className="border border-gray-300 text-gray-600 px-3 py-2 rounded-lg text-sm hover:bg-gray-50"
            >
              クリア
            </button>
          )}
        </form>

        <div className="bg-white rounded-xl border border-gray-200 p-4">
          {loading ? (
            <LoadingSpinner />
          ) : users.length === 0 ? (
            <div className="text-center py-8 text-gray-400">ユーザーが見つかりません</div>
          ) : (
            users.map((user) => (
              <UserCard
                key={user.id}
                user={user}
                showAdmin={currentUser?.admin && currentUser.id !== user.id}
                onDelete={handleDelete}
              />
            ))
          )}
        </div>

        {pagination && <Pagination pagination={pagination} onPageChange={setPage} />}
      </div>
    </Layout>
  );
}
