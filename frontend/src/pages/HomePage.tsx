import { useAtom } from 'jotai';
import { useState } from 'react';
import { Link } from 'react-router-dom';
import Layout from '../components/Layout';
import LoadingSpinner from '../components/LoadingSpinner';
import MicropostCard from '../components/MicropostCard';
import MicropostForm from '../components/MicropostForm';
import Pagination from '../components/Pagination';
import UserStatBar from '../components/UserStatBar';
import { useFeed } from '../hooks/useFeed';
import { currentUserAtom } from '../store/auth';

export default function HomePage() {
  // Jotai atom からログイン中のユーザーを取得
  // undefined = /me 確認中、null = 未ログイン、User = ログイン済み
  const [currentUser] = useAtom(currentUserAtom);
  const [page, setPage] = useState(1);

  // フィードとサイドバー用プロフィールを取得する（カスタムフック）
  const { feed, profile, loading, error, addPost, removePost, updatePost } = useFeed(
    currentUser,
    page,
  );

  // 認証確認中は全画面スピナー（ウェルカム画面の一瞬表示を防ぐ）
  if (currentUser === undefined) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner />
      </div>
    );
  }

  // 未ログインの場合はウェルカム画面を表示
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
        {/* サイドバー */}
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

        {/* メインフィード */}
        <div className="md:col-span-2 space-y-4">
          {/* 投稿フォーム：送信成功時に addPost でフィードの先頭に追加 */}
          <MicropostForm onCreated={addPost} />

          {error && (
            <div className="p-3 bg-red-50 border border-red-200 text-red-700 text-sm rounded-lg">
              {error}
            </div>
          )}
          {loading ? (
            <div className="text-center py-10 text-gray-400">読み込み中...</div>
          ) : !feed || feed.items.length === 0 ? (
            <div className="text-center py-10 text-gray-400">
              <p>まだ投稿がありません。</p>
              <p className="text-sm mt-1">
                <Link to="/users" className="text-blue-600 hover:underline">
                  他のユーザーをフォロー
                </Link>
                して投稿を見てみましょう。
              </p>
            </div>
          ) : (
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
                <Pagination pagination={feed.pagination} onPageChange={setPage} />
              )}
            </>
          )}
        </div>
      </div>
    </Layout>
  );
}
