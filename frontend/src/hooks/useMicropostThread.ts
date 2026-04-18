/**
 * 1 件のマイクロポストと、その返信一覧を取得するフック。
 *
 * API（getMicropost）は 1 回のレスポンスで post + replies を返しますが、
 * SWR のキャッシュは「そのレスポンス丸ごと」1 キーに紐づきます。
 *
 * `mutate` に関数を渡すと「前のキャッシュを受け取り、新しいキャッシュを返す」更新ができます。
 * いいね後の post 更新や、返信追加で replies を触るときに revalidate: false で即座に UI を合わせます。
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
