import useSWR from 'swr';
import { getUserLikes } from '../api/client';
import { getErrorMessage } from '../api/errors';
import type { Micropost, Pagination, UserStatSummary } from '../types';

export function useUserLikes(id: string | undefined, page: number) {
  const key = id ? `likes-${id}-p${page}` : null;

  const { data, isLoading: loading, error, mutate } = useSWR(
    key,
    () => getUserLikes(Number(id), page),
  );

  const posts: Micropost[] = data?.microposts ?? [];
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

  const errorMessage = error ? getErrorMessage(error, 'いいね一覧の取得に失敗しました') : null;

  return { posts, statSummary, pagination, loading, error: errorMessage, mutate };
}
