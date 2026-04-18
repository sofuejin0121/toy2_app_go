import useSWR from 'swr';
import { getMicropost } from '../api/client';
import type { Micropost } from '../types';

export function useMicropostThread(id: string | undefined) {
  const { data, isLoading: loading, mutate } = useSWR(
    id ? `micropost-thread-${id}` : null,
    () => getMicropost(Number(id)),
  );

  const post = data?.post ?? null;
  const replies = data?.replies ?? [];

  function setPost(p: Micropost) {
    mutate((d) => (d ? { ...d, post: p } : d), { revalidate: false });
  }

  function setReplies(updater: (prev: Micropost[]) => Micropost[]) {
    mutate((d) => (d ? { ...d, replies: updater(d.replies ?? []) } : d), { revalidate: false });
  }

  return { post, replies, loading, mutate, setPost, setReplies };
}
