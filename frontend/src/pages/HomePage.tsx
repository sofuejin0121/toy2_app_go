import { useAtom } from 'jotai';
import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { getFeed, getUser } from '../api/client';
import { getErrorMessage } from '../api/errors';
import Layout from '../components/Layout';
import MicropostCard from '../components/MicropostCard';
import MicropostForm from '../components/MicropostForm';
import Pagination from '../components/Pagination';
import UserStatBar from '../components/UserStatBar';
import { currentUserAtom } from '../store/auth';
import type { Micropost, Pagination as PaginationType, UserProfile } from '../types';

export default function HomePage() {
  const [currentUser] = useAtom(currentUserAtom);
  const [posts, setPosts] = useState<Micropost[]>([]);
  const [pagination, setPagination] = useState<PaginationType | null>(null);
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    // undefined = /me 取得中。null になってからウェルカム画面を表示する
    if (currentUser === undefined) return;
    if (currentUser === null) {
      setLoading(false);
      return;
    }
    setLoading(true);
    setError(null);
    Promise.all([getFeed(page), getUser(currentUser.id)])
      .then(([feed, prof]) => {
        setPosts(feed.items);
        setPagination(feed.pagination);
        setProfile(prof);
      })
      .catch((err: unknown) => setError(getErrorMessage(err, 'フィードの取得に失敗しました')))
      .finally(() => setLoading(false));
  }, [currentUser, page]);

  // /me 取得中は画面全体をスピナーで待機（ウェルカム画面の一瞬表示を防ぐ）
  if (currentUser === undefined) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin w-8 h-8 border-2 border-blue-600 border-t-transparent rounded-full" />
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
          <MicropostForm onCreated={(post) => setPosts((prev) => [post, ...prev])} />

          {error && (
            <div className="p-3 bg-red-50 border border-red-200 text-red-700 text-sm rounded-lg">
              {error}
            </div>
          )}
          {loading ? (
            <div className="text-center py-10 text-gray-400">読み込み中...</div>
          ) : posts.length === 0 ? (
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
              {posts.map((post) => (
                <MicropostCard
                  key={post.id}
                  post={post}
                  onDelete={(id) => setPosts((prev) => prev.filter((p) => p.id !== id))}
                  onUpdate={(updated) =>
                    setPosts((prev) => prev.map((p) => (p.id === updated.id ? updated : p)))
                  }
                />
              ))}
              {pagination && <Pagination pagination={pagination} onPageChange={setPage} />}
            </>
          )}
        </div>
      </div>
    </Layout>
  );
}
