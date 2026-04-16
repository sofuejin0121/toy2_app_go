/**
 * 通知一覧のデータ取得フック
 *
 * 役割: ログインユーザーの通知一覧を取得し、削除操作も提供する
 *       GET /api/notifications に対応
 *
 * 使い方:
 *   const { notifications, loading, deleteOne } = useNotifications();
 *
 * 戻り値:
 *   notifications - 通知の配列
 *   loading       - 取得中は true
 *   deleteOne     - 通知を1件削除する関数（引数: 通知ID）
 */
import { useEffect, useState } from 'react';
import { deleteNotification, listNotifications } from '../api/client';
import type { Notification } from '../types';

export function useNotifications() {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function loadNotifications() {
      try {
        setLoading(true);
        const data = await listNotifications();
        setNotifications(data.notifications ?? []);
      } finally {
        setLoading(false);
      }
    }

    loadNotifications();
  }, []); // マウント時に1回だけ実行

  // 通知を1件削除してローカルの表示からも除去する
  async function deleteOne(notificationId: number) {
    await deleteNotification(notificationId);
    // API 成功後にローカルの state から除去（再フェッチ不要）
    setNotifications((prev) => prev.filter((n) => n.id !== notificationId));
  }

  return { notifications, loading, deleteOne };
}
