/**
 * 投稿スレッドのデータ取得フック
 *
 * 役割: 指定した投稿とそのリプライ一覧を取得する
 *       GET /api/microposts/:id に対応
 *
 * 使い方:
 *   const { post, replies, loading } = useMicropostThread(id);
 *
 * 戻り値:
 *   post    - 対象の投稿（取得前は null）
 *   replies - リプライの配列
 *   loading - 取得中は true
 */
import { useEffect, useState } from 'react';
import { getMicropost } from '../api/client';
import type { Micropost } from '../types';

export function useMicropostThread(id: string | undefined) {
  const [post, setPost] = useState<Micropost | null>(null);
  const [replies, setReplies] = useState<Micropost[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id) return;

    async function loadThread() {
      try {
        setLoading(true);
        const data = await getMicropost(Number(id));
        setPost(data.post);
        setReplies(data.replies ?? []);
      } finally {
        setLoading(false);
      }
    }

    loadThread();
  }, [id]);

  return { post, setPost, replies, setReplies, loading };
}
