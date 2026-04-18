import type { Dispatch, SetStateAction } from 'react';
import { useState } from 'react';
import useSWR from 'swr';
import { getSettings, updateSettings } from '../api/client';
import { getErrorMessage } from '../api/errors';
import type { Settings } from '../types';

interface Alert {
  type: 'success' | 'error';
  message: string;
}

export function useSettings() {
  const { data: settings, isLoading: loading, mutate } = useSWR('settings', getSettings);

  const setSettings: Dispatch<SetStateAction<Settings | null>> = (action) => {
    mutate(
      (prev) => {
        const cur = prev ?? null;
        const next =
          typeof action === 'function'
            ? (action as (p: Settings | null) => Settings | null)(cur)
            : action;
        if (next === null) return prev;
        return next;
      },
      { revalidate: false },
    );
  };

  const [saving, setSaving] = useState(false);
  const [alert, setAlert] = useState<Alert | null>(null);

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

  return { settings, setSettings, loading, saving, alert, save };
}
