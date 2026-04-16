import Layout from '../components/Layout';
import LoadingSpinner from '../components/LoadingSpinner';
import { useSettings } from '../hooks/useSettings';

export default function SettingsPage() {
  // 通知設定の読み込みと保存を管理するカスタムフック
  // settings - 現在の設定値
  // setSettings - チェックボックスの変更をローカルに反映する関数
  // save - 保存ボタン押下時に API へ送信する関数
  const { settings, setSettings, loading, saving, alert, save } = useSettings();

  if (loading)
    return (
      <Layout>
        <LoadingSpinner />
      </Layout>
    );

  return (
    <Layout alert={alert || undefined}>
      <div className="max-w-lg mx-auto">
        <div className="bg-white rounded-xl border border-gray-200 p-8">
          <h1 className="text-2xl font-bold text-gray-900 mb-6">通知設定</h1>

          {settings && (
            <form
              onSubmit={(e) => {
                e.preventDefault();
                save(settings);
              }}
              className="space-y-4"
            >
              <div className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                <div>
                  <p className="font-medium text-gray-900 text-sm">フォロー通知メール</p>
                  <p className="text-xs text-gray-500 mt-0.5">
                    フォローされたときにメールを受け取る
                  </p>
                </div>
                <label className="relative inline-flex items-center cursor-pointer">
                  <input
                    type="checkbox"
                    checked={settings.email_on_follow}
                    onChange={(e) =>
                      setSettings((s) => (s ? { ...s, email_on_follow: e.target.checked } : s))
                    }
                    className="sr-only peer"
                  />
                  <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:bg-blue-600 after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all" />
                </label>
              </div>

              <div className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                <div>
                  <p className="font-medium text-gray-900 text-sm">いいね通知メール</p>
                  <p className="text-xs text-gray-500 mt-0.5">いいねされたときにメールを受け取る</p>
                </div>
                <label className="relative inline-flex items-center cursor-pointer">
                  <input
                    type="checkbox"
                    checked={settings.email_on_like}
                    onChange={(e) =>
                      setSettings((s) => (s ? { ...s, email_on_like: e.target.checked } : s))
                    }
                    className="sr-only peer"
                  />
                  <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:bg-blue-600 after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all" />
                </label>
              </div>

              <button
                type="submit"
                disabled={saving}
                className="w-full bg-blue-600 text-white py-2 rounded-lg text-sm font-medium hover:bg-blue-700 disabled:opacity-50"
              >
                {saving ? '保存中...' : '設定を保存'}
              </button>
            </form>
          )}
        </div>
      </div>
    </Layout>
  );
}
