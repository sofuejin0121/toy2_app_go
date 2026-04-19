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
import type { AlertState } from '../types';

/**
 * ユーザー詳細ページ（/users/:id）。
 *
 * データの出どころ:
 * - URL の `:id` とページ番号 `page` を `useUserProfile` に渡す。
 * - フック内部の useSWR が `user-${id}-page-${page}` キーで GET /api/users/:id を呼び、
 *   ユーザー情報・フォロー状態・マイクロポスト一覧・ページネーションをまとめて取得する。
 *
 * フォロー / アンフォロー:
 * - 成功後は手で `profile` の一部だけ書き換えず、`mutate()` でサーバーから取り直す（リフェッチ）。
 *   楽観的更新よりコードが単純で、フォロワー数などの表示ズレも起きにくい。
 *
 * 表示上の分岐:
 * - `currentUserAtom` … ログイン中なら User、未ログインなら null、未取得なら undefined（main の AuthLoader 経由）。
 * - フォローボタンは「ログイン済み かつ 他人のプロフィール」のときだけ出す。
 */
export default function UserShowPage() {
  // App.tsx の <Route path="/users/:id" /> から来る。常に string（未定義ならフック側で key を null にしてフェッチしない）
  const { id } = useParams<{ id: string }>();
  const [currentUser] = useAtom(currentUserAtom);

  // 投稿一覧のページネーション用。変わると useUserProfile の SWR キーも変わり、別ページを取りに行く
  const [page, setPage] = useState(1);
  // フォロー API 中の二重クリック防止とボタン表示用
  const [followLoading, setFollowLoading] = useState(false);
  // Layout 上部に出す一時メッセージ（エラー時など）
  const [alert, setAlert] = useState<AlertState | null>(null);

  const { profile, loading, error, mutate } = useUserProfile(id, page);

  // 初回・ページ切替の取得中。プロフィール本体はまだ描画しない
  if (loading)
    return (
      <Layout>
        <LoadingSpinner />
      </Layout>
    );

  // SWR の失敗やネットワークエラー。フック側で文字列に整形済み
  if (error)
    return (
      <Layout alert={alert || undefined}>
        <ErrorMessage message={error} />
      </Layout>
    );

  // 404 相当や id 不正などで API がデータを返さない場合
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
        // アンフォローは relationship の id が必要（バックエンドの設計）
        const rid = profile.relationship_id;
        if (rid == null) throw new Error('relationship_id missing');
        await unfollow(rid);
      } else {
        await follow(user.id);
      }
      // 引数なし mutate = この SWR キーのデータをサーバーから再取得（フォロー状態・数値・一覧を一括で正しい値に）
      await mutate();
    } catch (_e) {
      setAlert({ type: 'error', message: 'フォロー操作に失敗しました' });
    }
    setFollowLoading(false);
  };

  // async ハンドラを onClick に直接渡すと未処理 Promise になりやすいので分離し、失敗は握りつぶさず catch のみ
  const onFollowClick = () => {
    handleFollow().catch(() => {});
  };

  // 投稿削除後に一覧をサーバーと一致させる（MicropostCard は Promise を返さないのでここで mutate）
  const onMicropostDeleted = () => {
    mutate().catch(() => {});
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
                onClick={onFollowClick}
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
                  onDelete={onMicropostDeleted}
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
