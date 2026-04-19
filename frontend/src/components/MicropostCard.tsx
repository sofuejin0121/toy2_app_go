/**
 * 1 件のマイクロポストをカード表示する。
 *
 * このコンポーネント内で行う処理:
 * - いいね / 解除 → like / unlike API → ローカル state と親への onUpdate で件数を同期
 * - ブックマーク / 解除 → bookmark / unbookmark
 * - 削除（本人 or 管理者）→ deleteMicropost → onDelete(id) で親のリストから外す等
 *
 * 未ログイン時はいいね・BM ボタンを無効化し、削除も出さない（currentUser で判定）。
 */
import { useAtom } from 'jotai';
import { useState } from 'react';
import { Link } from 'react-router-dom';
import { bookmark, deleteMicropost, like, unbookmark, unlike } from '../api/client';
import { currentUserAtom } from '../store/auth';
import type { Micropost } from '../types';
import { timeAgo } from '../utils/timeAgo';

interface Props {
  post: Micropost;
  onDelete?: (id: number) => void;
  onUpdate?: (post: Micropost) => void;
}

export default function MicropostCard({ post, onDelete, onUpdate }: Props) {
  const [currentUser] = useAtom(currentUserAtom);
  const [likeCount, setLikeCount] = useState(post.like_count);
  const [isLiked, setIsLiked] = useState(post.is_liked);
  const [isBookmarked, setIsBookmarked] = useState(post.is_bookmarked);
  const [loading, setLoading] = useState(false);
  const [isOpenImage, setIsOpenImage] = useState(false);



  const handleOpenImage = () => {
    setIsOpenImage(true);
  }
  const handleCloseImage = () => {
    setIsOpenImage(false);
  }
  const handleLike = async () => {
    if (!currentUser || loading) return;
    setLoading(true);
    try {
      if (isLiked) {
        const res = await unlike(post.id);
        setIsLiked(false);
        setLikeCount(res.count);
        // ホームのフィード等: 親が持つ配列上の post も最新の count / is_liked に揃える
        onUpdate?.({ ...post, is_liked: false, like_count: res.count });
      } else {
        const res = await like(post.id);
        setIsLiked(true);
        setLikeCount(res.count);
        onUpdate?.({ ...post, is_liked: true, like_count: res.count });
      }
    } catch (e) {
      console.error(e);
    }
    setLoading(false);
  };

  const handleBookmark = async () => {
    if (!currentUser || loading) return;
    setLoading(true);
    try {
      if (isBookmarked) {
        await unbookmark(post.id);
        setIsBookmarked(false);
      } else {
        await bookmark(post.id);
        setIsBookmarked(true);
      }
    } catch (e) {
      console.error(e);
    }
    setLoading(false);
  };

  const handleDelete = async () => {
    if (!window.confirm('削除しますか？')) return;
    try {
      await deleteMicropost(post.id);
      onDelete?.(post.id);
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <article className="bg-white rounded-xl border border-gray-200 p-4 hover:border-gray-300 transition-colors">
      {/* リプライ元表示 */}
      {post.parent && (
        <div className="mb-2 pl-3 border-l-2 border-gray-200">
          <p className="text-xs text-gray-500">
            <Link to={`/users/${post.parent.user.id}`} className="font-medium hover:underline">
              {post.parent.user.name}
            </Link>{' '}
            への返信
          </p>
          <p className="text-xs text-gray-400 truncate">{post.parent.content}</p>
        </div>
      )}

      <div className="flex gap-3">
        <Link to={`/users/${post.user.id}`} className="shrink-0">
          <img
            src={post.user.avatar_url}
            alt={post.user.name}
            className="w-10 h-10 rounded-full object-cover"
          />
        </Link>

        <div className="flex-1 min-w-0">
          <div className="flex items-baseline gap-2 flex-wrap">
            <Link
              to={`/users/${post.user.id}`}
              className="font-semibold text-gray-900 hover:underline text-sm"
            >
              {post.user.name}
            </Link>
            <span className="text-xs text-gray-400">{timeAgo(post.created_at)}</span>
          </div>

          <Link to={`/microposts/${post.id}`}>
            <p className="mt-1 text-gray-800 text-sm whitespace-pre-wrap break-words">
              {post.content}
            </p>
          </Link>

          {post.image_path && (
            <img
              src={post.image_path}
              alt="投稿画像"
              onClick={handleOpenImage}
              className="mt-2 rounded-lg max-h-64 w-auto object-cover cursor-pointer"
            />
          )}

          {isOpenImage && (
            <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
              <img src={post.image_path} alt="投稿画像" className="max-w-full max-h-full object-contain" />
              <button onClick={handleCloseImage} className="absolute top-4 right-4 text-white text-2xl">×</button>
            </div>
          )}

          {/* アクションボタン */}
          <div className="mt-2 flex items-center gap-4">
            {/* いいね */}
            <button
              type="button"
              onClick={handleLike}
              disabled={!currentUser}
              className={`flex items-center gap-1 text-sm transition-colors ${
                isLiked ? 'text-red-500' : 'text-gray-400 hover:text-red-400'
              } disabled:opacity-50`}
            >
              <svg
                className="w-4 h-4"
                fill={isLiked ? 'currentColor' : 'none'}
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
                />
              </svg>
              <span>{likeCount}</span>
            </button>

            {/* リプライ */}
            <Link
              to={`/microposts/${post.id}`}
              className="flex items-center gap-1 text-sm text-gray-400 hover:text-blue-400 transition-colors"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                />
              </svg>
            </Link>

            {/* ブックマーク */}
            {currentUser && (
              <button
                type="button"
                onClick={handleBookmark}
                className={`flex items-center gap-1 text-sm transition-colors ${
                  isBookmarked ? 'text-yellow-500' : 'text-gray-400 hover:text-yellow-400'
                }`}
              >
                <svg
                  className="w-4 h-4"
                  fill={isBookmarked ? 'currentColor' : 'none'}
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"
                  />
                </svg>
              </button>
            )}

            {/* 削除ボタン（自分の投稿のみ） */}
            {currentUser && (currentUser.id === post.user_id || currentUser.admin) && (
              <button
                type="button"
                onClick={handleDelete}
                className="ml-auto text-xs text-gray-300 hover:text-red-400 transition-colors"
              >
                削除
              </button>
            )}
          </div>
        </div>
      </div>
    </article>
  );
}
