/**
 * 管理者統計情報の取得フック
 *
 * 役割: 管理者ページ用の統計情報を取得する
 *       管理者以外のアクセスはホームへリダイレクトする
 *       GET /api/admin/stats に対応
 *
 * 使い方:
 *   const { stats, loading } = useAdminStats(currentUser, navigate);
 *
 * 引数:
 *   currentUser - ログイン中のユーザー（admin フラグを確認）
 *   navigate    - ページ遷移関数（非管理者のリダイレクトに使用）
 */
import { useEffect, useState } from 'react';
import type { NavigateFunction } from 'react-router-dom';
import { getAdminStats } from '../api/client';
import type { AdminStats, User } from '../types';

export function useAdminStats(currentUser: User | null | undefined, navigate: NavigateFunction) {
  const [stats, setStats] = useState<AdminStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // undefined = 認証確認中のためスキップ
    if (currentUser === undefined) return;
    // 管理者以外はホームへリダイレクト
    if (!currentUser?.admin) {
      navigate('/');
      return;
    }

    async function loadStats() {
      try {
        setLoading(true);
        const data = await getAdminStats();
        setStats(data);
      } finally {
        setLoading(false);
      }
    }

    loadStats();
  }, [currentUser, navigate]);

  return { stats, loading };
}
