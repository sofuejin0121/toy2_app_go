import useSWR from 'swr';
import { getFeed, getUser } from '../api/client';
import { getErrorMessage } from '../api/errors';
import type { Micropost, Pagination, User } from '../types';

interface FeedData {
  items: Micropost[];
  pagination: Pagination;
}

export function useFeed(currentUser: User | null | undefined, page: number) {
  const feedKey = currentUser ? `feed-page-${page}` : null;
  const profileKey = currentUser ? `feed-profile-${currentUser.id}` : null;

  const { data: feed, isLoading: loadingFeed, error: errFeed, mutate: mutateFeed } = useSWR<
    FeedData | undefined
  >(feedKey, () => getFeed(page));

  const { data: profile, isLoading: loadingProfile, error: errProfile } = useSWR(
    profileKey,
    () => {
      if (!currentUser) throw new Error('feed profile requires login');
      return getUser(currentUser.id);
    },
  );

  const loading =
    currentUser !== undefined && currentUser !== null && (loadingFeed || loadingProfile);

  const rawError = errFeed ?? errProfile;
  const error = rawError ? getErrorMessage(rawError, 'フィードの取得に失敗しました') : null;

  function addPost(post: Micropost) {
    mutateFeed((prev) => (prev ? { ...prev, items: [post, ...prev.items] } : prev), {
      revalidate: false,
    });
  }

  function removePost(postId: number) {
    mutateFeed(
      (prev) => (prev ? { ...prev, items: prev.items.filter((p) => p.id !== postId) } : prev),
      { revalidate: false },
    );
  }

  function updatePost(updated: Micropost) {
    mutateFeed(
      (prev) =>
        prev
          ? { ...prev, items: prev.items.map((p) => (p.id === updated.id ? updated : p)) }
          : prev,
      { revalidate: false },
    );
  }

  return {
    feed: feed ?? null,
    profile: profile ?? null,
    loading,
    error,
    addPost,
    removePost,
    updatePost,
    mutateFeed,
  };
}
