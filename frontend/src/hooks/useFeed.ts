/**
 * ホームのフィード用フック。
 *
 * 以前は 1 つの useEffect で Promise.all していましたが、SWR では「データの塊ごと」に useSWR を分けます。
 * - フィード本体 … `feed-page-${page}`
 * - サイドバー用プロフィール … `feed-profile-${userId}`
 * それぞれ独立してキャッシュ・再取得されるので、片方だけ失敗したときの扱いもしやすくなります。
 *
 * mutate と revalidate: false:
 * - 投稿直後など「サーバーに取りに行かず、今のキャッシュだけ手で直す」ときに使います。
 * - 再フェッチしたいときは `mutate()` だけ呼ぶ（オプション省略）など、別の書き方もできます。
 */
import useSWR from 'swr';
import { getFeed, getUser } from '../api/client';
import { getErrorMessage } from '../api/errors';
import type { Micropost, Pagination, User } from '../types';
import { updateIfDefined } from '../utils/updateIfDefined';

interface FeedData {
  items: Micropost[];
  pagination: Pagination;
}

export function useFeed(currentUser: User | null | undefined, page: number) {
  // 未ログインのときは key を null にして、SWR にフェッチさせない
  const feedKey = currentUser ? `feed-page-${page}` : null;
  const profileKey = currentUser ? `feed-profile-${currentUser.id}` : null;

  const { data: feed, isLoading: loadingFeed, error: errFeed, mutate: mutateFeed } = useSWR<
    FeedData | undefined
  >(feedKey, () => getFeed(page));

  const { data: profile, isLoading: loadingProfile, error: errProfile } = useSWR(
    profileKey,
    () => {
      // profileKey が付いている時点で currentUser は存在する想定だが、型のため明示的にチェック
      if (!currentUser) throw new Error('feed profile requires login');
      return getUser(currentUser.id);
    },
  );

  // ゲスト（null）や未取得（undefined）のときは「フィード用ローディング」を出さない
  const loading =
    currentUser !== undefined && currentUser !== null && (loadingFeed || loadingProfile);

  const rawError = errFeed ?? errProfile;
  const error = rawError ? getErrorMessage(rawError, 'フィードの取得に失敗しました') : null;

  function addPost(post: Micropost) {
    mutateFeed((prev) => updateIfDefined(prev, (p) => ({ ...p, items: [post, ...p.items] })), {
      revalidate: false,
    });
  }

  function removePost(postId: number) {
    mutateFeed(
      (prev) => updateIfDefined(prev, (p) => ({ ...p, items: p.items.filter((x) => x.id !== postId) })),
      { revalidate: false },
    );
  }

  function updatePost(updated: Micropost) {
    mutateFeed(
      (prev) =>
        updateIfDefined(prev, (p) => ({
          ...p,
          items: p.items.map((x) => (x.id === updated.id ? updated : x)),
        })),
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
