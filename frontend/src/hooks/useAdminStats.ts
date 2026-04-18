import useSWR from 'swr';
import { getAdminStats } from '../api/client';
import type { AdminStats } from '../types';

export function useAdminStats() {
  const { data: stats, isLoading: loading } = useSWR<AdminStats>('admin-stats', getAdminStats);

  return { stats: stats ?? null, loading };
}
