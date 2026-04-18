import { useAtom } from 'jotai';
import { useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { follow, unfollow } from '../api/client';
import ErrorMessage from '../components/ErrorMessage';
import Layout from '../components/Layout';
import LoadingSpinner from '../components/LoadingSpinner';
import MicropostCard from '../components/MicropostCard';
import Pagination from '../components/Pagination';
import UserStatBar from '../components/UserStatBar';
import { useUserProfile } from '../hooks/useUserProfile';
import { currentUserAtom } from '../store/auth';

/**
 * プロフィール表示ページ。
 * フォロー / アンフォロー後は、手で state をいじるのではなく useSWR の mutate() を呼び出し、
 * サーバーからプロフィールを取り直します（楽観的更新より実装が単純で、数値のズレも起きにくい）。
 */
export default function UserShowPage() {
  const { id } = useParams<{ id: string }>();
  const [currentUser] = useAtom(currentUserAtom);
  const [page, setPage] = useState(1);
  const [followLoading, setFollowLoading] = useState(false);
  const [alert, setAlert] = useState<{ type: 'success' | 'error'; message: string } | null>(null);

  const { profile, loading, error, mutate } = useUserProfile(id, page);

  if (loading)
    return (
      <Layout>
        <LoadingSpinner />
      </Layout>
    );

  if (error)
    return (
      <Layout alert={alert || undefined}>
        <ErrorMessage message={error} />
      </Layout>
    );

  if (!profile)
    return (
      <Layout alert={alert || undefined}>
        <div className="text-center py-10 text-gray-400">ユーザーが見つかりません</div>
      </Layout>
    );

  const { user } = profile;

  const handleFollow = async () => {
    if (followLoading) return;
    setFollowLoading(true);
    try {
      if (profile.is_following) {
        const rid = profile.relationship_id;
        if (rid == null) throw new Error('relationship_id missing');
        await unfollow(rid);
      } else {
        await follow(user.id);
      }
      // 引数なし mutate = この key のデータをサーバーから再取得（リフェッチ）
      await mutate();
    } catch (_e) {
      setAlert({ type: 'error', message: 'フォロー操作に失敗しました' });
    }
    setFollowLoading(false);
  };

  return (
    <Layout alert={alert || undefined}>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
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
                type="button"
                onClick={() => void handleFollow()}
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

        <div className="md:col-span-2 space-y-4">
          <h2 className="text-lg font-semibold text-gray-900">投稿</h2>
          {profile.microposts.length === 0 ? (
            <div className="text-center py-10 text-gray-400">まだ投稿がありません</div>
          ) : (
            <>
              {profile.microposts.map((post) => (
                <MicropostCard
                  key={post.id}
                  post={post}
                  onDelete={() => void mutate()}
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
