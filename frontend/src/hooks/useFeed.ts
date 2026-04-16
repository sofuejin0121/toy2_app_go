/**
 * ホームフィードのデータ取得フック
 *
 * 役割: ログインユーザーのフィード（投稿一覧）とサイドバー用プロフィールを取得する
 *
 * 使い方:
 *   const { feed, profile, loading, error } = useFeed(currentUser, page);
 *
 * 引数:
 *   currentUser - ログイン中のユーザー（null=未ログイン, undefined=確認中）
 *   page        - 表示するページ番号（1始まり）
 *
 * 戻り値:
 *   feed    - フィードデータ（投稿一覧 + ページネーション情報）
 *   profile - サイドバー用プロフィール（GET /api/users/:id のレスポンス）
 *   loading - データ取得中は true
 *   error   - エラーメッセージ（正常時は null）
 */
import { useEffect, useState } from 'react';
import { getFeed, getUser } from '../api/client';
import { getErrorMessage } from '../api/errors';
import type { Micropost, Pagination, User, UserProfile } from '../types';

// getFeed のレスポンス型
interface FeedData {
  items: Micropost[];
  pagination: Pagination;
}

export function useFeed(currentUser: User | null | undefined, page: number) {
  const [feed, setFeed] = useState<FeedData | null>(null);
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    // undefined = /me 取得中。データが確定してから処理する
    if (currentUser === undefined) return;

    // 未ログインなら取得不要
    if (currentUser === null) {
      setLoading(false);
      return;
    }

    // クロージャ内で TypeScript の型ガードを有効にするためローカル変数に代入する
    const user = currentUser;

    // useEffect のコールバックは async にできないため、
    // 内部で async 関数を定義してすぐ呼び出す（Go の goroutine に近いイメージ）
    async function loadFeed() {
      try {
        setLoading(true);
        setError(null);
        // フィードとプロフィールを並列取得（Go の goroutine + WaitGroup に相当）
        const [feedData, profileData] = await Promise.all([
          getFeed(page),
          getUser(user.id),
        ]);
        setFeed(feedData);
        setProfile(profileData);
      } catch (err) {
        setError(getErrorMessage(err, 'フィードの取得に失敗しました'));
      } finally {
        setLoading(false);
      }
    }

    loadFeed();
  }, [currentUser, page]); // currentUser か page が変わったら再取得

  // 新しい投稿をフィードの先頭に追加する（投稿フォーム送信後に呼ぶ）
  function addPost(post: Micropost) {
    setFeed((prev) => (prev ? { ...prev, items: [post, ...prev.items] } : prev));
  }

  // 投稿をフィードから除去する（削除後に呼ぶ）
  function removePost(postId: number) {
    setFeed((prev) =>
      prev ? { ...prev, items: prev.items.filter((p) => p.id !== postId) } : prev,
    );
  }

  // 投稿の内容を更新する（いいね/ブックマーク後に呼ぶ）
  function updatePost(updated: Micropost) {
    setFeed((prev) =>
      prev ? { ...prev, items: prev.items.map((p) => (p.id === updated.id ? updated : p)) } : prev,
    );
  }

  return { feed, profile, loading, error, addPost, removePost, updatePost };
}
