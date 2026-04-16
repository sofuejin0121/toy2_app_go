import { useEffect, useState } from 'react';
import { getSettings, updateSettings } from '../api/client';
import Layout from '../components/Layout';
import type { Settings } from '../types';

export default function SettingsPage() {
  const [settings, setSettings] = useState<Settings | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [alert, setAlert] = useState<{ type: 'success' | 'error'; message: string } | null>(null);

  useEffect(() => {
    getSettings()
      .then((data) => setSettings(data))
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!settings) return;
    setSaving(true);
    try {
      const updated = await updateSettings(settings);
      setSettings(updated);
      setAlert({ type: 'success', message: '設定を保存しました' });
    } catch (_e) {
      setAlert({ type: 'error', message: '保存に失敗しました' });
    }
    setSaving(false);
  };

  if (loading)
    return (
      <Layout>
        <div className="text-center py-10 text-gray-400">読み込み中...</div>
      </Layout>
    );

  return (
    <Layout alert={alert || undefined}>
      <div className="max-w-lg mx-auto">
        <div className="bg-white rounded-xl border border-gray-200 p-8">
          <h1 className="text-2xl font-bold text-gray-900 mb-6">通知設定</h1>

          {settings && (
            <form onSubmit={handleSubmit} className="space-y-4">
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
