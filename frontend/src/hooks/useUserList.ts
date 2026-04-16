/**
 * ユーザー一覧のデータ取得フック
 *
 * 役割: ページネーション付きのユーザー一覧を取得する
 *       GET /api/users?page=:page&q=:query に対応
 *
 * 使い方:
 *   const { users, pagination, loading } = useUserList(page, query);
 *
 * 引数:
 *   page  - ページ番号（1始まり）
 *   query - 検索キーワード（空文字列で全件取得）
 *
 * 戻り値:
 *   users      - ユーザーの配列
 *   pagination - ページネーション情報
 *   loading    - 取得中は true
 */
import { useEffect, useState } from 'react';
import { listUsers } from '../api/client';
import type { Pagination, User } from '../types';

export function useUserList(page: number, query: string) {
  const [users, setUsers] = useState<User[]>([]);
  const [pagination, setPagination] = useState<Pagination | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function loadUsers() {
      try {
        setLoading(true);
        const data = await listUsers(page, query);
        setUsers(data.users);
        setPagination(data.pagination);
      } finally {
        setLoading(false);
      }
    }

    loadUsers();
  }, [page, query]); // page か query が変わったら再取得

  return { users, setUsers, pagination, loading };
}
