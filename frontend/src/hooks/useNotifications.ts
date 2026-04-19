/**
 * 通知一覧（GET /notifications）。
 *
 * - SWR key は固定文字列なので「通知画面」のキャッシュは常に 1 つ。
 * - deleteOne … API で削除後、mutate でキャッシュの配列から該当通知だけ除去（全件 GET し直さない）。
 */
import useSWR from 'swr';
import { deleteNotification, listNotifications } from '../api/client';
import type { Notification } from '../types';
import { updateIfDefined } from '../utils/updateIfDefined';

export function useNotifications() {
  const { data, isLoading: loading, mutate } = useSWR('notifications', listNotifications);

  // `?.` … 未取得は undefined / `??` … 通知リストは空配列にする
  const notifications: Notification[] = data?.notifications ?? [];

  async function deleteOne(notificationId: number) {
    await deleteNotification(notificationId);
    // 一覧を取り直さず、キャッシュ上の配列だけから該当 ID を除く（prev が無いときは updateIfDefined が何もしない）
    mutate(
      (prev) =>
        updateIfDefined(prev, (p) => ({
          ...p,
          notifications: p.notifications.filter((n) => n.id !== notificationId),
        })),
      { revalidate: false },
    );
  }

  return { notifications, loading, deleteOne };
}
