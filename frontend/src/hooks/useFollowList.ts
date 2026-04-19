/**
 * フォロー中 / フォロワー一覧（GET .../following または .../followers）。
 *
 * - `mode` と URL の `id` と `page` が SWR の key に入るため、切り替えのたびに別キャッシュとして扱われる。
 * - id が無いときは key を null にしてフェッチしない（型安全のためのガード）。
 * - 返却する statSummary は UserStatBar 用に、API のトップレベル user と各 count をまとめたオブジェクト。
 */
import useSWR from 'swr';
import { getFollowers, getFollowing } from '../api/client';
import { getErrorMessage } from '../api/errors';
import type { Pagination, User, UserStatSummary } from '../types';

export function useFollowList(
  id: string | undefined,
  mode: 'following' | 'followers',
  page: number,
) {
  const key = id ? `followlist-${mode}-${id}-p${page}` : null;

  const { data, isLoading: loading, error } = useSWR(
    key,
    async () => {
      const fetch = mode === 'following' ? getFollowing : getFollowers;
      return fetch(Number(id), page);
    },
  );

  // `?.` … 未取得は undefined / `??` … ユーザー配列は空として扱う
  const users: User[] = data?.users ?? [];
  const pagination: Pagination | null = data?.pagination ?? null;

  const statSummary: UserStatSummary | null = data
    ? {
        user: data.user,
        micropost_count: data.micropost_count,
        following_count: data.following_count,
        followers_count: data.followers_count,
        liked_count: data.liked_count,
        bookmark_count: data.bookmark_count,
        is_current_user: data.is_current_user,
      }
    : null;

  const errorMessage = error ? getErrorMessage(error, 'ユーザー一覧の取得に失敗しました') : null;

  return { users, statSummary, pagination, loading, error: errorMessage };
}
