import useSWR from 'swr';
import { getUserBookmarks } from '../api/client';
import { getErrorMessage } from '../api/errors';
import type { Micropost, Pagination, UserStatSummary } from '../types';

export function useUserBookmarks(id: string | undefined, page: number) {
  const key = id ? `bookmarks-${id}-p${page}` : null;

  const { data, isLoading: loading, error, mutate } = useSWR(
    key,
    () => getUserBookmarks(Number(id), page),
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
        is_current_user: true,
      }
    : null;

  const errorMessage = error ? getErrorMessage(error, 'ブックマークの取得に失敗しました') : null;

  return { posts, statSummary, pagination, loading, error: errorMessage, mutate };
}
