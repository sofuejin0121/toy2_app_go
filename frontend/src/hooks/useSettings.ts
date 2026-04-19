/**
 * メール通知設定（GET/PATCH /settings）。
 *
 * - 読み込み: useSWR('settings', getSettings)。チェック変更は patchSettings でキャッシュのみ即時反映。
 * - save: PATCH 成功後に mutate(updated, { revalidate: false }) でサーバー応答とキャッシュを一致させ、余計な GET を避ける。
 * - alert: Layout に渡す成功・失敗メッセージ用（SettingsPage が表示）。
 */
import { useState } from 'react';
import useSWR from 'swr';
import { getSettings, updateSettings } from '../api/client';
import { getErrorMessage } from '../api/errors';
import type { AlertState, Settings } from '../types';

export function useSettings() {
  const { data: settings, isLoading: loading, mutate } = useSWR('settings', getSettings);

  /** チェックボックス変更など、キャッシュ上の設定を一部だけ差し替える */
  function patchSettings(patch: Partial<Settings>) {
    mutate(
      (prev) => {
        if (prev == null) return prev;
        return { ...prev, ...patch };
      },
      { revalidate: false },
    );
  }

  const [saving, setSaving] = useState(false);
  const [alert, setAlert] = useState<AlertState | null>(null);

  async function save(updated: Settings) {
    try {
      setSaving(true);
      await updateSettings(updated);
      mutate(updated, { revalidate: false });
      setAlert({ type: 'success', message: '設定を保存しました' });
    } catch (err) {
      setAlert({ type: 'error', message: getErrorMessage(err, '設定の保存に失敗しました') });
    } finally {
      setSaving(false);
    }
  }

  return { settings, patchSettings, loading, saving, alert, save };
}
