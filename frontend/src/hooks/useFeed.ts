/**
 * ホームのフィード用フック（ログイン済みユーザー向け）。
 *
 * useSWR を 2 つに分ける理由:
 * - フィード本体 … key `feed-page-${page}` → getFeed（投稿一覧）
 * - サイドバー用 … key `feed-profile-${userId}` → getUser（UserStatBar に必要な統計）
 * 失敗・再取得を独立させ、ゲスト（currentUser === null）では key を null にしてフェッチしない。
 *
 * SWR オプション:
 * - dedupingInterval: 0 … ホームに戻るたびに最新の HTTP を取りに行く。
 * - リトライ … 一時的な 5xx 向け。401/403/404 はリトライしない。
 * - 表示用エラー … 再取得中（isValidating）は画面に出さず、短い遅延後にのみ出す（一瞬の赤文字を防ぐ）。
 *
 * addPost / removePost / updatePost:
 * - mutateFeed に関数を渡し、updateIfDefined で「キャッシュ未取得なら触らない」更新をする（revalidate: false）。
 */
import axios from 'axios';
import { useEffect, useMemo, useState } from 'react';
import useSWR from 'swr';
import { getFeed, getUser } from '../api/client';
import { getErrorMessage } from '../api/errors';
import type { Micropost, Pagination, User } from '../types';
import { updateIfDefined } from '../utils/updateIfDefined';

interface FeedData {
  items: Micropost[];
  pagination: Pagination;
}

function shouldRetrySwrError(err: unknown): boolean {
  if (axios.isAxiosError(err)) {
    const s = err.response?.status;
    if (s === 401 || s === 403 || s === 404) return false;
  }
  return true;
}

const swrHomeOpts = {
  dedupingInterval: 0,
  errorRetryCount: 6,
  /** 指数バックオフ相当に近づけるため徐々に長く（SWR の型は number のみ） */
  errorRetryInterval: 2500,
  shouldRetryOnError: shouldRetrySwrError,
};

export function useFeed(currentUser: User | null | undefined, page: number) {
  // 未ログインのときは key を null にして、SWR にフェッチさせない
  const feedKey = currentUser ? `feed-page-${page}` : null;
  const profileKey = currentUser ? `feed-profile-${currentUser.id}` : null;

  const {
    data: feed,
    isLoading: loadingFeed,
    isValidating: validatingFeed,
    error: errFeed,
    mutate: mutateFeed,
  } = useSWR<FeedData | undefined>(feedKey, () => getFeed(page), swrHomeOpts);

  const {
    data: profile,
    isLoading: loadingProfile,
    isValidating: validatingProfile,
    error: errProfile,
  } = useSWR(
    profileKey,
    () => {
      // profileKey が付いている時点で currentUser は存在する想定だが、型のため明示的にチェック
      if (!currentUser) throw new Error('feed profile requires login');
      return getUser(currentUser.id);
    },
    swrHomeOpts,
  );

  const validating = validatingFeed || validatingProfile;

  // ゲスト（null）や未取得（undefined）のときは「フィード用ローディング」を出さない。
  // データ未取得のまま再検証中（リトライ含む）も読み込み扱いにし、エラー文言のチラつきを抑える。
  const loading =
    currentUser !== undefined &&
    currentUser !== null &&
    (loadingFeed ||
      loadingProfile ||
      ((!feed || !profile) && validating));

  const rawError = errFeed ?? errProfile;
  const tentativeError = useMemo(() => {
    if (!rawError || validating) return null;
    return getErrorMessage(rawError, 'フィードの取得に失敗しました');
  }, [rawError, validating]);

  const [error, setError] = useState<string | null>(null);
  useEffect(() => {
    if (!tentativeError) {
      setError(null);
      return;
    }
    const id = window.setTimeout(() => setError(tentativeError), 400);
    return () => window.clearTimeout(id);
  }, [tentativeError]);

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
