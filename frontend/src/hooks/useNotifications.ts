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
