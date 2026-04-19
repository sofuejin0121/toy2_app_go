/**
 * 1 件のマイクロポスト詳細＋返信一覧（GET /microposts/:id）。
 *
 * - SWR キャッシュは post + replies が一体となったオブジェクト。
 * - setPost / setReplies … MicropostPage からいいね更新や返信追加後にキャッシュだけ部分更新するときに使う。
 * - mutate に渡す更新関数では updateIfDefined を使い、未取得時は undefined のままにする。
 */
import useSWR from 'swr';
import { getMicropost } from '../api/client';
import type { Micropost } from '../types';
import { updateIfDefined } from '../utils/updateIfDefined';

export function useMicropostThread(id: string | undefined) {
  const { data, isLoading: loading, mutate } = useSWR(
    id ? `micropost-thread-${id}` : null,
    () => getMicropost(Number(id)),
  );

  const post = data?.post ?? null;
  // `?.` … data がまだ無いときは undefined を避ける / `??` … 左が null/undefined なら右（空配列）を使う
  const replies = data?.replies ?? [];

  function setPost(p: Micropost) {
    mutate((d) => updateIfDefined(d, (cur) => ({ ...cur, post: p })), { revalidate: false });
  }

  function setReplies(updater: (prev: Micropost[]) => Micropost[]) {
    mutate(
      (d) => updateIfDefined(d, (cur) => ({ ...cur, replies: updater(cur.replies ?? []) })),
      { revalidate: false },
    );
  }

  return { post, replies, loading, mutate, setPost, setReplies };
}
