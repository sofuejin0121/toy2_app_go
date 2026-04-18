/**
 * ユーザープロフィール（1 人分の詳細 + 投稿一覧）を SWR で取得するフック。
 *
 * SWR の基本:
 * - 第 1 引数 `key` … キャッシュの識別子。文字列が変わると「別データ」として再フェッチされます。
 * - `null` を渡すと「まだ取得しない」（条件付きフェッチ）。id が無いときに使います。
 * - 第 2 引数 `fetcher` … key に対応するデータを返す Promise 関数。既存の API ラッパー（getUser）をそのまま使えます。
 * - 戻り値の `mutate` … キャッシュを手動で更新したり、サーバーに取りに行き直したりする関数です。
 */
import useSWR from 'swr';
import { getUser } from '../api/client';
import { getErrorMessage } from '../api/errors';
import type { UserProfile } from '../types';

export function useUserProfile(id: string | undefined, page = 1) {
  // ページ番号も key に含めると、ページを変えたときに自動で別リクエストになる
  const { data: profile, isLoading: loading, error, mutate } = useSWR<UserProfile>(
    id ? `user-${id}-page-${page}` : null,
    () => getUser(Number(id), page),
  );

  // SWR の error は unknown 系なので、画面に出しやすい文字列に変換する
  const errorMessage = error ? getErrorMessage(error, 'プロフィールの取得に失敗しました') : null;

  return { profile: profile ?? null, loading, error: errorMessage, mutate };
}
