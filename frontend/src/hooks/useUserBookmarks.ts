/**
 * ブックマーク一覧のデータ取得フック
 *
 * 役割: ログインユーザー自身のブックマーク一覧を取得する
 *       GET /api/users/:id/bookmarks に対応
 *       ※ 他人のブックマークは見られないため、本人確認も行う
 *
 * 使い方:
 *   const { posts, statSummary, pagination, loading, error } =
 *     useUserBookmarks(id, currentUser, page);
 *
 * 引数:
 *   id          - URL パラメータのユーザー ID
 *   currentUser - Jotai atom から取得したログイン中のユーザー
 *   navigate    - ページ遷移関数（本人以外のアクセス時にリダイレクト）
 *   page        - ページ番号
 */
import { useEffect, useState } from 'react';
import type { NavigateFunction } from 'react-router-dom';
import { getUserBookmarks } from '../api/client';
import { getErrorMessage } from '../api/errors';
import type { Micropost, Pagination, User, UserStatSummary } from '../types';

export function useUserBookmarks(
  id: string | undefined,
  currentUser: User | null | undefined,
  navigate: NavigateFunction,
  page: number,
) {
  const [posts, setPosts] = useState<Micropost[]>([]);
  const [statSummary, setStatSummary] = useState<UserStatSummary | null>(null);
  const [pagination, setPagination] = useState<Pagination | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    // undefined = 認証確認中のためスキップ
    if (currentUser === undefined) return;
    // 未ログイン、または自分以外のページにアクセスした場合はホームへリダイレクト
    if (!currentUser || String(currentUser.id) !== id) {
      navigate('/');
      return;
    }

    async function loadBookmarks() {
      try {
        setLoading(true);
        setError(null);
        const data = await getUserBookmarks(Number(id), page);
        setPosts(data.microposts);
        setPagination(data.pagination);
        setStatSummary({
          user: data.user,
          micropost_count: data.micropost_count,
          following_count: data.following_count,
          followers_count: data.followers_count,
          liked_count: data.liked_count,
          bookmark_count: data.bookmark_count,
          is_current_user: true,
        });
      } catch (err) {
        setError(getErrorMessage(err, 'ブックマークの取得に失敗しました'));
      } finally {
        setLoading(false);
      }
    }

    loadBookmarks();
  }, [id, currentUser, page, navigate]);

  return { posts, setPosts, statSummary, pagination, loading, error };
}
