/**
 * ユーザー一覧＋検索（GET /users?page=&q=）。
 *
 * - SWR key に page と query の両方を含める → 検索語やページが変わるたびに自動で再フェッチ。
 * - mutate は削除後に一覧キャッシュから 1 行除くなど、UserListPage から `updateIfDefined` と組み合わせて使う。
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
