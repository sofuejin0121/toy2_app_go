/**
 * 通知設定のデータ取得・更新フック
 *
 * 役割: 通知設定の読み込みと保存を担当する
 *       GET /api/settings / PATCH /api/settings に対応
 *
 * 使い方:
 *   const { settings, loading, saving, alert, save } = useSettings();
 *
 * 戻り値:
 *   settings - 現在の設定値（取得前は null）
 *   loading  - 初期読み込み中は true
 *   saving   - 保存処理中は true
 *   alert    - 保存結果のメッセージ（成功/エラー）
 *   save     - 設定を保存する関数（引数: 更新後の settings）
 */
import { useEffect, useState } from 'react';
import { getSettings, updateSettings } from '../api/client';
import { getErrorMessage } from '../api/errors';
import type { Settings } from '../types';

interface Alert {
  type: 'success' | 'error';
  message: string;
}

export function useSettings() {
  const [settings, setSettings] = useState<Settings | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [alert, setAlert] = useState<Alert | null>(null);

  useEffect(() => {
    async function loadSettings() {
      try {
        setLoading(true);
        const data = await getSettings();
        setSettings(data);
      } finally {
        setLoading(false);
      }
    }

    loadSettings();
  }, []); // マウント時に1回だけ実行

  // 設定を保存する
  async function save(updated: Settings) {
    try {
      setSaving(true);
      await updateSettings(updated);
      setSettings(updated);
      setAlert({ type: 'success', message: '設定を保存しました' });
    } catch (err) {
      setAlert({ type: 'error', message: getErrorMessage(err, '設定の保存に失敗しました') });
    } finally {
      setSaving(false);
    }
  }

  return { settings, setSettings, loading, saving, alert, save };
}
