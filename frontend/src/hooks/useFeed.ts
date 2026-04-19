/**
 * ホームのフィード用フック（ログイン済みユーザー向け）。
 *
 * useSWR を 2 つに分ける理由:
 * - フィード本体 … key `feed-page-${page}` → getFeed（投稿一覧）
 * - サイドバー用 … key `feed-profile-${userId}` → getUser（UserStatBar に必要な統計）
 * 失敗・再取得を独立させ、ゲスト（currentUser === null）では key を null にしてフェッチしない。
 *
 * SWR オプション:
 * - 以前: main の SWRConfig の dedupingInterval: 2000 のみ。直前の 401 などがキャッシュに残ったまま
 *   別ページからすぐホームに戻ると、dedupe で再フェッチされず古いエラーが表示され得た。
 * - 今回: このフックだけ dedupingInterval: 0 にし、ホームに戻るたびに最新の HTTP を取りに行く。
 *
 * addPost / removePost / updatePost:
 * - mutateFeed に関数を渡し、updateIfDefined で「キャッシュ未取得なら触らない」更新をする（revalidate: false）。
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

  const swrDedupeOff = { dedupingInterval: 0 };

  const { data: feed, isLoading: loadingFeed, error: errFeed, mutate: mutateFeed } = useSWR<
    FeedData | undefined
  >(feedKey, () => getFeed(page), swrDedupeOff);

  const { data: profile, isLoading: loadingProfile, error: errProfile } = useSWR(
    profileKey,
    () => {
      // profileKey が付いている時点で currentUser は存在する想定だが、型のため明示的にチェック
      if (!currentUser) throw new Error('feed profile requires login');
      return getUser(currentUser.id);
    },
    swrDedupeOff,
  );

  // ゲスト（null）や未取得（undefined）のときは「フィード用ローディング」を出さない
  const loading =
    currentUser !== undefined && currentUser !== null && (loadingFeed || loadingProfile);

  const rawError = errFeed ?? errProfile;
  const error = rawError ? getErrorMessage(rawError, 'フィードの取得に失敗しました') : null;

  function addPost(post: Micropost) {
    mutateFeed((prev) => updateIfDefined(prev, (p) => ({ ...p, items: [post, ...p.items] })), {
      revalidate: false, //フィードのキャッシュを更新しない
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
