/**
 * 管理者ダッシュボード用の統計。
 * AdminRoute で「管理者だけがこのページに来る」ことが保証されているので、
 * フック側ではログインチェックや navigate を行いません。シンプルに固定 key で GET します。
 */
import useSWR from 'swr';
import { getAdminStats } from '../api/client';
import type { AdminStats } from '../types';

export function useAdminStats() {
  const { data: stats, isLoading: loading } = useSWR<AdminStats>('admin-stats', getAdminStats);

  return { stats: stats ?? null, loading };
}
