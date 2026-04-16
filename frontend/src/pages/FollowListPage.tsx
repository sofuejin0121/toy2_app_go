import { useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { getFollowers, getFollowing } from '../api/client';
import { getErrorMessage } from '../api/errors';
import Layout from '../components/Layout';
import Pagination from '../components/Pagination';
import UserCard from '../components/UserCard';
import UserStatBar from '../components/UserStatBar';
import type { Pagination as PaginationType, User, UserStatSummary } from '../types';

interface Props {
  mode: 'following' | 'followers';
}

export default function FollowListPage({ mode }: Props) {
  const { id } = useParams<{ id: string }>();
  const [users, setUsers] = useState<User[]>([]);
  // フォロー一覧APIは liked_count / bookmark_count を返さないため UserStatSummary を部分的に保持
  const [profileData, setProfileData] = useState<Omit<
    UserStatSummary,
    'liked_count' | 'bookmark_count' | 'is_current_user'
  > | null>(null);
  const [pagination, setPagination] = useState<PaginationType | null>(null);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;
    setLoading(true);
    const fn = mode === 'following' ? getFollowing : getFollowers;
    fn(Number(id), page)
      .then((data) => {
        setUsers(data.users);
        setPagination(data.pagination);
        setProfileData({
          user: data.user,
          following_count: data.following_count,
          followers_count: data.followers_count,
          micropost_count: data.micropost_count,
        });
      })
      .catch((err: unknown) => setError(getErrorMessage(err, 'ユーザー一覧の取得に失敗しました')))
      .finally(() => setLoading(false));
  }, [id, mode, page]);

  // UserStatBar 用に不足フィールドを補完する（フォロー一覧APIでは liked_count 等が返らない）
  const statSummary: UserStatSummary | null = profileData
    ? { ...profileData, liked_count: 0, bookmark_count: 0, is_current_user: false }
    : null;

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
              <div className="text-center py-8 text-gray-400">読み込み中...</div>
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
