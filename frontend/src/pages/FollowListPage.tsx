import { useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import Layout from '../components/Layout';
import LoadingSpinner from '../components/LoadingSpinner';
import Pagination from '../components/Pagination';
import UserCard from '../components/UserCard';
import UserStatBar from '../components/UserStatBar';
import { useFollowList } from '../hooks/useFollowList';

interface Props {
  mode: 'following' | 'followers';
}

export default function FollowListPage({ mode }: Props) {
  const { id } = useParams<{ id: string }>();
  const [page, setPage] = useState(1);

  // フォロー/フォロワー一覧を取得するカスタムフック
  const { users, statSummary, pagination, loading, error } = useFollowList(id, mode, page);

  return (
    <Layout>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {statSummary && (
          <aside className="md:col-span-1">
            <div className="bg-white rounded-xl border border-gray-200 p-4 sticky top-20">
              <div className="text-center mb-4">
                <img
                  src={statSummary.user.avatar_url}
                  alt={statSummary.user.name}
                  className="w-16 h-16 rounded-full mx-auto mb-2"
                />
                <Link
                  to={`/users/${statSummary.user.id}`}
                  className="font-bold text-gray-900 hover:underline"
                >
                  {statSummary.user.name}
                </Link>
              </div>
              <UserStatBar profile={statSummary} />
            </div>
          </aside>
        )}

        <div className="md:col-span-2">
          <h2 className="text-xl font-bold text-gray-900 mb-4">
            {mode === 'following' ? 'フォロー中' : 'フォロワー'}
          </h2>
          <div className="bg-white rounded-xl border border-gray-200 p-4">
            {error && (
              <div className="p-3 bg-red-50 border border-red-200 text-red-700 text-sm rounded-lg">
                {error}
              </div>
            )}
            {loading ? (
              <LoadingSpinner />
            ) : users.length === 0 ? (
              <div className="text-center py-8 text-gray-400">まだいません</div>
            ) : (
              users.map((user) => <UserCard key={user.id} user={user} />)
            )}
          </div>
          {pagination && <Pagination pagination={pagination} onPageChange={setPage} />}
        </div>
      </div>
    </Layout>
  );
}
