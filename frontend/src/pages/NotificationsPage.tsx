import { Link } from 'react-router-dom';
import Layout from '../components/Layout';
import LoadingSpinner from '../components/LoadingSpinner';
import { useNotifications } from '../hooks/useNotifications';
import { timeAgo } from '../utils/timeAgo';

// 通知の種類に応じたラベルを返す（API の action_type 文字列を人が読める日本語に）
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

/**
 * 通知一覧（ProtectedRoute）。useNotifications の deleteOne で一覧キャッシュを部分更新。
 */
export default function NotificationsPage() {
  const { notifications, loading, deleteOne } = useNotifications();

  return (
    <Layout>
      <div className="max-w-2xl mx-auto">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">通知</h1>

        <div className="bg-white rounded-xl border border-gray-200 divide-y divide-gray-100">
          {loading ? (
            <LoadingSpinner />
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
                  type="button"
                  onClick={() => deleteOne(n.id)}
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
