import { useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { getFollowers, getFollowing } from '../api/client';
import Layout from '../components/Layout';
import Pagination from '../components/Pagination';
import UserCard from '../components/UserCard';
import UserStatBar from '../components/UserStatBar';
import type { Pagination as PaginationType, User, UserProfile } from '../types';

interface Props {
  mode: 'following' | 'followers';
}

export default function FollowListPage({ mode }: Props) {
  const { id } = useParams<{ id: string }>();
  const [users, setUsers] = useState<User[]>([]);
  const [profileData, setProfileData] = useState<{
    user: UserProfile['user'];
    following_count: number;
    followers_count: number;
    micropost_count: number;
  } | null>(null);
  const [pagination, setPagination] = useState<PaginationType | null>(null);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);

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
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [id, mode, page]);

  const fakeProfile = profileData
    ? {
        user: profileData.user,
        is_current_user: false,
        is_following: false,
        micropost_count: profileData.micropost_count,
        following_count: profileData.following_count,
        followers_count: profileData.followers_count,
        liked_count: 0,
        bookmark_count: 0,
        microposts: [],
        pagination: pagination || {
          current_page: 1,
          total_pages: 1,
          total_items: 0,
          per_page: 30,
          has_prev: false,
          has_next: false,
        },
      }
    : null;

  return (
    <Layout>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {fakeProfile && (
          <aside className="md:col-span-1">
            <div className="bg-white rounded-xl border border-gray-200 p-4 sticky top-20">
              <div className="text-center mb-4">
                <img
                  src={fakeProfile.user.avatar_url}
                  alt={fakeProfile.user.name}
                  className="w-16 h-16 rounded-full mx-auto mb-2"
                />
                <Link
                  to={`/users/${fakeProfile.user.id}`}
                  className="font-bold text-gray-900 hover:underline"
                >
                  {fakeProfile.user.name}
                </Link>
              </div>
              <UserStatBar profile={fakeProfile} />
            </div>
          </aside>
        )}

        <div className="md:col-span-2">
          <h2 className="text-xl font-bold text-gray-900 mb-4">
            {mode === 'following' ? 'フォロー中' : 'フォロワー'}
          </h2>
          <div className="bg-white rounded-xl border border-gray-200 p-4">
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
