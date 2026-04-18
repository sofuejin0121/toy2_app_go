import useSWR from 'swr';
import { getUser } from '../api/client';
import { getErrorMessage } from '../api/errors';
import type { UserProfile } from '../types';

export function useUserProfile(id: string | undefined, page = 1) {
  const { data: profile, isLoading: loading, error, mutate } = useSWR<UserProfile>(
    id ? `user-${id}-page-${page}` : null,
    () => getUser(Number(id), page),
  );

  const errorMessage = error ? getErrorMessage(error, 'プロフィールの取得に失敗しました') : null;

  return { profile: profile ?? null, loading, error: errorMessage, mutate };
}
