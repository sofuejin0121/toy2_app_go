/**
 * ユーザー一覧（検索クエリ query と page が key に入る）。
 * 検索ボタンで query が変わると SWR が別キャッシュとして扱うため、古い結果と混ざりません。
 *
 * mutate は管理者がユーザーを削除したあと、一覧からその行だけ除く用途などで使います。
 */
import useSWR from 'swr';
import { listUsers } from '../api/client';
import type { Pagination, User } from '../types';

export function useUserList(page: number, query: string) {
  const key = `users-p${page}-q${query}`;

  const { data, isLoading: loading, mutate } = useSWR(key, () => listUsers(page, query));

  // `?.` … まだ取得前は data が undefined / `??` … そのときは空配列として扱う
  const users: User[] = data?.users ?? [];
  const pagination: Pagination | null = data?.pagination ?? null;

  return { users, pagination, loading, mutate };
}
