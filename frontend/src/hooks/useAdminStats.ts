/**
 * 管理者ダッシュボード用の統計（GET /admin）。
 *
 * - SWR の key は固定文字列 `admin-stats`（画面が 1 枚なのでキャッシュも 1 つで足りる）。
 * - App.tsx の AdminRoute が「管理者以外はここに来ない」ことを保証するため、フック内では権限チェックしない。
 * - 戻り値の stats は未取得時 null（AdminPage で loading 後に参照する）。
 */
import useSWR from 'swr';
import { getAdminStats } from '../api/client';
import type { AdminStats } from '../types';

export function useAdminStats() {
  const { data: stats, isLoading: loading } = useSWR<AdminStats>('admin-stats', getAdminStats);

  return { stats: stats ?? null, loading }; //SWR の data が未取得のとき undefined なので、ページ側では null＝データなしとして扱いやすくするため
}
