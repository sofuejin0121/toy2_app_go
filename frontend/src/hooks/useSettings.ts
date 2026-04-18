/**
 * 通知設定の取得（SWR）と保存（通常の async 関数）。
 *
 * - 読み込み: useSWR('settings', getSettings) … マウント時に自動フェッチ、タブ復帰時の再検証は SWR のデフォルト動作に任せる。
 * - チェックボックス: patchSettings でキャッシュ上の settings を部分的に書き換える（React の setState より単純な引数だけ）。
 * - 保存成功後: mutate(updated, { revalidate: false }) でサーバーと同じ内容をキャッシュに直接書き込み、余計な GET を避ける。
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
