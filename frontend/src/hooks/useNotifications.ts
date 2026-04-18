/**
 * 通知一覧。マウント時に SWR が自動で GET します（key が固定文字列なので「この画面用のキャッシュ」が 1 つ）。
 *
 * 削除後に一覧を取り直す代わりに、mutate でキャッシュから該当 ID だけ除いています。
 * 体感が速く、サーバー負荷も減ります。取り直したい場合は mutate() で再検証も可能です。
 */
import useSWR from 'swr';
import { deleteNotification, listNotifications } from '../api/client';
import type { Notification } from '../types';

export function useNotifications() {
  const { data, isLoading: loading, mutate } = useSWR('notifications', listNotifications);

  const notifications: Notification[] = data?.notifications ?? [];

  async function deleteOne(notificationId: number) {
    await deleteNotification(notificationId);
    mutate(
      (prev) =>
        prev
          ? { notifications: prev.notifications.filter((n) => n.id !== notificationId) }
          : prev,
      { revalidate: false },
    );
  }

  return { notifications, loading, deleteOne };
}
