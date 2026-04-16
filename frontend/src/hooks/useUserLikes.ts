/**
 * いいね一覧のデータ取得フック
 *
 * 役割: 指定ユーザーがいいねした投稿一覧を取得する
 *       GET /api/users/:id/likes に対応
 *
 * 使い方:
 *   const { posts, statSummary, pagination, loading, error } = useUserLikes(id, page);
 *
 * 戻り値:
 *   posts      - いいねした投稿の配列
 *   statSummary - サイドバー用の統計情報
 *   pagination - ページネーション情報
 *   loading    - 取得中は true
 *   error      - エラーメッセージ
 */
import { useEffect, useState } from 'react';
import { getUserLikes } from '../api/client';
import { getErrorMessage } from '../api/errors';
import type { Micropost, Pagination, UserStatSummary } from '../types';

export function useUserLikes(id: string | undefined, page: number) {
  const [posts, setPosts] = useState<Micropost[]>([]);
  const [statSummary, setStatSummary] = useState<UserStatSummary | null>(null);
  const [pagination, setPagination] = useState<Pagination | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;

    async function loadLikes() {
      try {
        setLoading(true);
        setError(null);
        const data = await getUserLikes(Number(id), page);
        setPosts(data.microposts);
        setPagination(data.pagination);
        // いいね一覧APIはブックマーク数を返さないため 0 で補完
        setStatSummary({
          user: data.user,
          micropost_count: data.micropost_count,
          following_count: data.following_count,
          followers_count: data.followers_count,
          liked_count: data.liked_count,
          bookmark_count: 0,
          is_current_user: false,
        });
      } catch (err) {
        setError(getErrorMessage(err, 'いいね一覧の取得に失敗しました'));
      } finally {
        setLoading(false);
      }
    }

    loadLikes();
  }, [id, page]);

  return { posts, setPosts, statSummary, pagination, loading, error };
}
