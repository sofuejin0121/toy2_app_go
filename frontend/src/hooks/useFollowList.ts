/**
 * フォロー/フォロワー一覧のデータ取得フック
 *
 * 役割: 指定ユーザーのフォロー中、またはフォロワーの一覧を取得する
 *       GET /api/users/:id/following または GET /api/users/:id/followers に対応
 *
 * 使い方:
 *   const { users, statSummary, pagination, loading, error } =
 *     useFollowList(id, 'following', page);
 *
 * 引数:
 *   id   - 対象ユーザーの ID
 *   mode - 'following'（フォロー中）または 'followers'（フォロワー）
 *   page - ページ番号
 *
 * 戻り値:
 *   users      - ユーザーの配列
 *   statSummary - サイドバー用の統計情報（UserStatBar に渡す）
 *   pagination - ページネーション情報
 *   loading    - 取得中は true
 *   error      - エラーメッセージ
 */
import { useEffect, useState } from 'react';
import { getFollowers, getFollowing } from '../api/client';
import { getErrorMessage } from '../api/errors';
import type { Pagination, User, UserStatSummary } from '../types';

export function useFollowList(
  id: string | undefined,
  mode: 'following' | 'followers',
  page: number,
) {
  const [users, setUsers] = useState<User[]>([]);
  // フォロー一覧APIは liked_count / bookmark_count を返さないため、補完した型で保持する
  const [statSummary, setStatSummary] = useState<UserStatSummary | null>(null);
  const [pagination, setPagination] = useState<Pagination | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;

    async function loadList() {
      try {
        setLoading(true);
        setError(null);
        // mode によって呼ぶ API 関数を切り替える
        const fetch = mode === 'following' ? getFollowing : getFollowers;
        const data = await fetch(Number(id), page);
        setUsers(data.users);
        setPagination(data.pagination);
        // フォロー一覧API が返さないフィールドは 0 / false で補完
        setStatSummary({
          user: data.user,
          micropost_count: data.micropost_count,
          following_count: data.following_count,
          followers_count: data.followers_count,
          liked_count: 0,
          bookmark_count: 0,
          is_current_user: false,
        });
      } catch (err) {
        setError(getErrorMessage(err, 'ユーザー一覧の取得に失敗しました'));
      } finally {
        setLoading(false);
      }
    }

    loadList();
  }, [id, mode, page]);

  return { users, statSummary, pagination, loading, error };
}
