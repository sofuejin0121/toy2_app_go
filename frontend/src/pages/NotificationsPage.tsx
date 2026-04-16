import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { deleteNotification, listNotifications } from '../api/client';
import Layout from '../components/Layout';
import type { Notification } from '../types';

function timeAgo(dateStr: string): string {
  const diff = (Date.now() - new Date(dateStr).getTime()) / 1000;
  if (diff < 60) return `${Math.floor(diff)}秒前`;
  if (diff < 3600) return `${Math.floor(diff / 60)}分前`;
  if (diff < 86400) return `${Math.floor(diff / 3600)}時間前`;
  return `${Math.floor(diff / 86400)}日前`;
}

function actionLabel(type: string): string {
  switch (type) {
    case 'like':
      return 'があなたの投稿にいいねしました';
    case 'follow':
      return 'があなたをフォローしました';
    default:
      return `が${type}しました`;
  }
}

export default function NotificationsPage() {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    listNotifications()
      .then((data) => setNotifications(data.notifications || []))
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  const handleDelete = async (id: number) => {
    try {
      await deleteNotification(id);
      setNotifications((prev) => prev.filter((n) => n.id !== id));
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <Layout>
      <div className="max-w-2xl mx-auto">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">通知</h1>

        <div className="bg-white rounded-xl border border-gray-200 divide-y divide-gray-100">
          {loading ? (
            <div className="text-center py-10 text-gray-400">読み込み中...</div>
          ) : notifications.length === 0 ? (
            <div className="text-center py-10 text-gray-400">通知はありません</div>
          ) : (
            notifications.map((n) => (
              <div
                key={n.id}
                className={`flex items-start gap-3 p-4 ${!n.read ? 'bg-blue-50' : ''}`}
              >
                <Link to={`/users/${n.actor.id}`} className="shrink-0">
                  <img
                    src={n.actor.avatar_url}
                    alt={n.actor.name}
                    className="w-9 h-9 rounded-full"
                  />
                </Link>
                <div className="flex-1 min-w-0">
                  <p className="text-sm text-gray-800">
                    <Link to={`/users/${n.actor.id}`} className="font-semibold hover:underline">
                      {n.actor.name}
                    </Link>
                    {actionLabel(n.action_type)}
                  </p>
                  {n.target_content && (
                    <Link
                      to={n.target_id ? `/microposts/${n.target_id}` : '#'}
                      className="text-xs text-gray-500 mt-0.5 line-clamp-1 hover:underline block"
                    >
                      「{n.target_content}」
                    </Link>
                  )}
                  <p className="text-xs text-gray-400 mt-1">{timeAgo(n.created_at)}</p>
                </div>
                <button
                  onClick={() => handleDelete(n.id)}
                  className="text-xs text-gray-300 hover:text-red-400 shrink-0"
                >
                  ✕
                </button>
              </div>
            ))
          )}
        </div>
      </div>
    </Layout>
  );
}
