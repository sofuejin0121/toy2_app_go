import { useAtom } from 'jotai';
import { useState } from 'react';
import { Link } from 'react-router-dom';
import ErrorMessage from '../components/ErrorMessage';
import Layout from '../components/Layout';
import LoadingSpinner from '../components/LoadingSpinner';
import MicropostCard from '../components/MicropostCard';
import MicropostForm from '../components/MicropostForm';
import Pagination from '../components/Pagination';
import UserStatBar from '../components/UserStatBar';
import { useFeed } from '../hooks/useFeed';
import { currentUserAtom } from '../store/auth';
import type { Micropost, Pagination as PaginationMeta } from '../types';

/** HomePage 内サブコンポーネント用: フィード・ページネーション・削除/更新コールバックを受け取る */
interface HomeFeedSectionProps {
  loading: boolean;
  feed: { items: Micropost[]; pagination: PaginationMeta } | null;
  removePost: (postId: number) => void;
  updatePost: (post: Micropost) => void;
  onPageChange: (page: number) => void;
}

/**
 * ログイン後メインカラムのフィード部分。
 * ローディング / 空 / 一覧を三項のネストなしで early return 分岐する。
 */
function HomeFeedSection({
  loading,
  feed,
  removePost,
  updatePost,
  onPageChange,
}: HomeFeedSectionProps) {
  if (loading) {
    return <div className="text-center py-10 text-gray-400">読み込み中...</div>;
  }

  if (!feed || feed.items.length === 0) {
    return (
      <div className="text-center py-10 text-gray-400">
        <p>まだ投稿がありません。</p>
        <p className="text-sm mt-1">
          <Link to="/users" className="text-blue-600 hover:underline">
            他のユーザーをフォロー
          </Link>
          して投稿を見てみましょう。
        </p>
      </div>
    );
  }

  return (
    <>
      {feed.items.map((post) => (
        <MicropostCard
          key={post.id}
          post={post}
          onDelete={removePost}
          onUpdate={updatePost}
        />
      ))}
      {feed.pagination && (
        <Pagination pagination={feed.pagination} onPageChange={onPageChange} />
      )}
    </>
  );
}

/**
 * トップページ（/）。currentUser の undefined / null / User で表示を完全に分岐し、ログイン後は useFeed。
 */
export default function HomePage() {
  const [currentUser] = useAtom(currentUserAtom);
  const [page, setPage] = useState(1);

  const { feed, profile, loading, error, addPost, removePost, updatePost } = useFeed(
    currentUser,
    page,
  );

  if (currentUser === undefined) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner />
      </div>
    );
  }

  if (currentUser === null) {
    return (
      <Layout>
        <div className="text-center py-20">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">Chirp へようこそ</h1>
          <p className="text-gray-500 mb-8 text-lg">短いメッセージをシェアしよう</p>
          <div className="flex justify-center gap-4">
            <Link
              to="/signup"
              className="bg-blue-600 text-white px-8 py-3 rounded-full text-lg hover:bg-blue-700"
            >
              今すぐ登録
            </Link>
            <Link
              to="/login"
              className="border border-blue-600 text-blue-600 px-8 py-3 rounded-full text-lg hover:bg-blue-50"
            >
              ログイン
            </Link>
          </div>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <aside className="md:col-span-1">
          {profile && (
            <div className="bg-white rounded-xl border border-gray-200 p-4 sticky top-20">
              <div className="flex items-center gap-3 mb-3">
                <img
                  src={currentUser.avatar_url}
                  alt={currentUser.name}
                  className="w-12 h-12 rounded-full"
                />
                <div>
                  <Link
                    to={`/users/${currentUser.id}`}
                    className="font-semibold text-gray-900 hover:underline"
                  >
                    {currentUser.name}
                  </Link>
                  <p className="text-xs text-gray-500">{currentUser.email}</p>
                </div>
              </div>
              <UserStatBar profile={profile} />
              <div className="mt-3">
                <Link to="/users" className="text-sm text-blue-600 hover:underline">
                  ユーザー一覧を見る →
                </Link>
              </div>
            </div>
          )}
        </aside>

        <div className="md:col-span-2 space-y-4">
          <MicropostForm onCreated={addPost} />

          {error && <ErrorMessage message={error} />}
          <HomeFeedSection
            loading={loading}
            feed={feed}
            removePost={removePost}
            updatePost={updatePost}
            onPageChange={setPage}
          />
        </div>
      </div>
    </Layout>
  );
}
