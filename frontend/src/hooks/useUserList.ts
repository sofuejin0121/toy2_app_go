import useSWR from 'swr';
import { listUsers } from '../api/client';
import type { Pagination, User } from '../types';

export function useUserList(page: number, query: string) {
  const key = `users-p${page}-q${query}`;

  const { data, isLoading: loading, mutate } = useSWR(key, () => listUsers(page, query));

  const users: User[] = data?.users ?? [];
  const pagination: Pagination | null = data?.pagination ?? null;

  return { users, pagination, loading, mutate };
}
