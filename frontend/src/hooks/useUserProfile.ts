/**
 * ユーザープロフィールのデータ取得フック
 *
 * 役割: 指定した userId のプロフィール・投稿一覧・統計情報を取得する
 *       GET /api/users/:id に対応
 *
 * 使い方:
 *   const { profile, loading, error } = useUserProfile(id, page);
 *
 * 引数:
 *   id   - URL パラメータから取得したユーザーID文字列（例: "42"）
 *   page - 投稿一覧のページ番号
 *
 * 戻り値:
 *   profile - ユーザープロフィール（投稿一覧・統計情報を含む）
 *   loading - 取得中は true
 *   error   - エラーメッセージ（正常時は null）
 *   setProfile - プロフィールを直接更新する関数（フォロー状態の楽観的更新に使用）
 */
import { useEffect, useState } from 'react';
import { getUser } from '../api/client';
import { getErrorMessage } from '../api/errors';
import type { UserProfile } from '../types';

export function useUserProfile(id: string | undefined, page = 1) {
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;

    async function loadProfile() {
      try {
        setLoading(true);
        setError(null);
        const data = await getUser(Number(id), page);
        setProfile(data);
      } catch (err) {
        setError(getErrorMessage(err, 'プロフィールの取得に失敗しました'));
      } finally {
        setLoading(false);
      }
    }

    loadProfile();
  }, [id, page]); // id か page が変わったら再取得

  return { profile, loading, error, setProfile };
}
