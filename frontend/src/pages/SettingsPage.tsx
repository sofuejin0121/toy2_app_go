import Layout from '../components/Layout';
import LoadingSpinner from '../components/LoadingSpinner';
import SettingsToggleRow from '../components/SettingsToggleRow';
import { useSettings } from '../hooks/useSettings';

export default function SettingsPage() {
  // 通知設定の読み込みと保存を管理するカスタムフック
  // settings - 現在の設定値
  // patchSettings - チェックボックスの変更をキャッシュ上に部分的に反映する関数
  // save - 保存ボタン押下時に API へ送信する関数
  const { settings, patchSettings, loading, saving, alert, save } = useSettings();

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
              <SettingsToggleRow
                title="フォロー通知メール"
                description="フォローされたときにメールを受け取る"
                checked={settings.email_on_follow}
                onCheckedChange={(checked) => patchSettings({ email_on_follow: checked })}
              />

              <SettingsToggleRow
                title="いいね通知メール"
                description="いいねされたときにメールを受け取る"
                checked={settings.email_on_like}
                onCheckedChange={(checked) => patchSettings({ email_on_like: checked })}
              />

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
