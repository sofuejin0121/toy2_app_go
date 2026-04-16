import { useAtom } from 'jotai';
import { useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { follow, getUser, unfollow } from '../api/client';
import Layout from '../components/Layout';
import MicropostCard from '../components/MicropostCard';
import Pagination from '../components/Pagination';
import UserStatBar from '../components/UserStatBar';
import { currentUserAtom } from '../store/auth';
import type { Micropost, UserProfile } from '../types';

export default function UserShowPage() {
  const { id } = useParams<{ id: string }>();
  const [currentUser] = useAtom(currentUserAtom);
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [posts, setPosts] = useState<Micropost[]>([]);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [followLoading, setFollowLoading] = useState(false);
  const [alert, setAlert] = useState<{ type: 'success' | 'error'; message: string } | null>(null);

  useEffect(() => {
    if (!id) return;
    setLoading(true);
    getUser(Number(id), page)
      .then((data) => {
        setProfile(data);
        setPosts(data.microposts);
      })
      .catch(() => setAlert({ type: 'error', message: 'ユーザーが見つかりません' }))
      .finally(() => setLoading(false));
  }, [id, page]);

  const handleFollow = async () => {
    if (!profile || followLoading) return;
    setFollowLoading(true);
    try {
      if (profile.is_following) {
        await unfollow(profile.relationship_id!);
        setProfile((prev) =>
          prev
            ? {
                ...prev,
                is_following: false,
                followers_count: prev.followers_count - 1,
                relationship_id: undefined,
              }
            : prev,
        );
      } else {
        const res = await follow(profile.user.id);
        setProfile((prev) =>
          prev
            ? {
                ...prev,
                is_following: true,
                followers_count: prev.followers_count + 1,
                relationship_id: res.relationship_id,
              }
            : prev,
        );
      }
    } catch (e) {
      console.error(e);
    }
    setFollowLoading(false);
  };

  if (loading)
    return (
      <Layout>
        <div className="text-center py-10 text-gray-400">読み込み中...</div>
      </Layout>
    );
  if (!profile)
    return (
      <Layout alert={alert || undefined}>
        <div className="text-center py-10 text-gray-400">ユーザーが見つかりません</div>
      </Layout>
    );

  const { user } = profile;

  return (
    <Layout alert={alert || undefined}>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* サイドバー */}
        <aside className="md:col-span-1">
          <div className="bg-white rounded-xl border border-gray-200 p-4 sticky top-20">
            <div className="text-center mb-4">
              <img
                src={user.avatar_url}
                alt={user.name}
                className="w-20 h-20 rounded-full mx-auto mb-2"
              />
              <h1 className="text-xl font-bold text-gray-900">{user.name}</h1>
              {user.bio && <p className="text-sm text-gray-500 mt-1">{user.bio}</p>}
            </div>

            <UserStatBar profile={profile} />

            {currentUser && !profile.is_current_user && (
              <button
                onClick={handleFollow}
                disabled={followLoading}
                className={`mt-4 w-full py-2 rounded-full text-sm font-medium transition-colors ${
                  profile.is_following
                    ? 'border border-gray-300 text-gray-700 hover:border-red-300 hover:text-red-500'
                    : 'bg-blue-600 text-white hover:bg-blue-700'
                } disabled:opacity-50`}
              >
                {followLoading ? '...' : profile.is_following ? 'フォロー中' : 'フォロー'}
              </button>
            )}

            {profile.is_current_user && (
              <Link
                to={`/users/${user.id}/edit`}
                className="mt-4 block w-full text-center py-2 border border-gray-300 rounded-full text-sm text-gray-700 hover:bg-gray-50"
              >
                プロフィールを編集
              </Link>
            )}
          </div>
        </aside>

        {/* 投稿一覧 */}
        <div className="md:col-span-2 space-y-4">
          <h2 className="text-lg font-semibold text-gray-900">投稿</h2>
          {posts.length === 0 ? (
            <div className="text-center py-10 text-gray-400">まだ投稿がありません</div>
          ) : (
            <>
              {posts.map((post) => (
                <MicropostCard
                  key={post.id}
                  post={post}
                  onDelete={(postId) => setPosts((prev) => prev.filter((p) => p.id !== postId))}
                />
              ))}
              {profile.pagination && (
                <Pagination pagination={profile.pagination} onPageChange={setPage} />
              )}
            </>
          )}
        </div>
      </div>
    </Layout>
  );
}
