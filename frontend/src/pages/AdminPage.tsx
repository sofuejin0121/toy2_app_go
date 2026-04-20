import Layout from '../components/Layout';
import LoadingSpinner from '../components/LoadingSpinner';
import { useDocumentTitle } from '../hooks/useDocumentTitle';
import { useAdminStats } from '../hooks/useAdminStats';

/**
 * 管理者ダッシュボード（AdminRoute）。useAdminStats で集計を表示し、棒グラフは maxCount で正規化。
 */
export default function AdminPage() {
  const { stats, loading } = useAdminStats();

  useDocumentTitle('管理');

  if (loading)
    return (
      <Layout>
        <LoadingSpinner />
      </Layout>
    );
  if (!stats) return null;

  const maxCount = Math.max(...stats.daily_signups.map((d) => d.count), 1);

  return (
    <Layout>
      <div className="max-w-3xl mx-auto">
        <h1 className="text-2xl font-bold text-gray-900 mb-6">管理者ダッシュボード</h1>

        <div className="grid grid-cols-3 gap-4 mb-8">
          {[
            { label: '総ユーザー数', value: stats.total_users, color: 'text-blue-600' },
            { label: '総投稿数', value: stats.total_posts, color: 'text-green-600' },
            { label: '本日の新規登録', value: stats.today_signups, color: 'text-purple-600' },
          ].map((item) => (
            <div
              key={item.label}
              className="bg-white rounded-xl border border-gray-200 p-4 text-center"
            >
              <div className={`text-3xl font-bold ${item.color}`}>{item.value}</div>
              <div className="text-sm text-gray-500 mt-1">{item.label}</div>
            </div>
          ))}
        </div>

        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">過去7日間の新規登録</h2>
          {stats.daily_signups.length === 0 ? (
            <p className="text-gray-400 text-sm">データがありません</p>
          ) : (
            <div className="space-y-2">
              {stats.daily_signups.map((d) => (
                <div key={d.date} className="flex items-center gap-3">
                  <span className="text-xs text-gray-500 w-24 shrink-0">{d.date}</span>
                  <div className="flex-1 bg-gray-100 rounded-full h-5 overflow-hidden">
                    <div
                      className="h-full bg-blue-500 rounded-full transition-all"
                      style={{ width: `${(d.count / maxCount) * 100}%` }}
                    />
                  </div>
                  <span className="text-sm font-medium text-gray-700 w-6 text-right">
                    {d.count}
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </Layout>
  );
}
