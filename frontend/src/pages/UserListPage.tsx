import { useAtom } from 'jotai';
import { useState } from 'react';
import { deleteUser } from '../api/client';
import Layout from '../components/Layout';
import LoadingSpinner from '../components/LoadingSpinner';
import Pagination from '../components/Pagination';
import UserCard from '../components/UserCard';
import { useUserList } from '../hooks/useUserList';
import { currentUserAtom } from '../store/auth';

export default function UserListPage() {
  const [currentUser] = useAtom(currentUserAtom);
  const [page, setPage] = useState(1);
  // query = 実際に API に送る検索キーワード（送信ボタンで確定）
  const [query, setQuery] = useState('');
  // inputQuery = 入力欄の一時的な値（Enter/送信ボタンまで API には送らない）
  const [inputQuery, setInputQuery] = useState('');

  // ユーザー一覧を取得するカスタムフック
  const { users, setUsers, pagination, loading } = useUserList(page, query);

  // 検索フォーム送信
  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setPage(1);
    setQuery(inputQuery);
  };

  // ユーザー削除（管理者のみ）
  const handleDelete = async (id: number) => {
    if (!window.confirm('このユーザーを削除しますか？')) return;
    try {
      await deleteUser(id);
      setUsers((prev) => prev.filter((u) => u.id !== id));
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
