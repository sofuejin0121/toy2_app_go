/**
 * ユーザープロフィール（1 人分の詳細 + マイクロポスト一覧 + ページネーション）を SWR で取得。
 *
 * - key: `user-${id}-page-${page}` … id かページが変わると別リクエスト。
 * - key が null のときは id 未確定なのでフェッチしない。
 * - fetcher は client.getUser（GET /users/:id?page=）。
 * - mutate() … 引数なしでリフェッチ（フォロー後に UserShowPage が呼ぶ）、など。
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
