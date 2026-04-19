/**
 * 認証の「起動時ブートストラップ」用コンポーネント。
 *
 * 何をするか:
 * - main.tsx で <AuthLoader><App /></AuthLoader> のように App より外側に置く。
 * - GET /api/me（getMe）で結果を currentUserAtom に書き込む。
 * - authBootstrapEpochAtom が進むたびに getMe をやり直す（ログイン直後の再確認など）。
 *
 * currentUserAtom の値の意味（store/auth.ts とセットで覚える）:
 * - undefined … まだ getMe の結果が出ていない（アプリ全体では「確認中」）
 * - null … 未ログイン（セッションなし or 期限切れ）
 * - User … ログイン済み。Layout や各ページがこのオブジェクトを参照する。
 *
 * 以前: マウント時 1 回だけ getMe。ログイン前に飛ばした getMe が遅れて 401 になると
 * `.catch(() => setCurrentUser(null))` でログイン済み状態を消していた。
 * 今回: (1) effect のクリーンアップで「取り消し」フラグ (2) epoch で「古い getMe」を無視し、
 * ログイン後に進めた epoch の新しい getMe だけが atom を更新する。
 *
 * ページ側での使い方:
 * - 認証状態の読み書きはこのファイルを import する必要はない。
 * - `import { useAtom } from 'jotai'` と `import { currentUserAtom } from '../store/auth'` でよい。
 */
import { useAtomValue, useSetAtom } from 'jotai';
import { type ReactNode, useEffect, useRef } from 'react';
import { getMe } from '../api/client';
import { authBootstrapEpochAtom, currentUserAtom } from '../store/auth';

export function AuthLoader({ children }: { children: ReactNode }) {
  const setCurrentUser = useSetAtom(currentUserAtom);
  const epoch = useAtomValue(authBootstrapEpochAtom);
  const epochRef = useRef(epoch);
  epochRef.current = epoch;

  useEffect(() => {
    let cancelled = false;
    const startedEpoch = epoch;

    getMe()
      .then((user) => {
        if (cancelled || epochRef.current !== startedEpoch) return;
        setCurrentUser(user);
      })
      .catch(() => {
        if (cancelled || epochRef.current !== startedEpoch) return;
        setCurrentUser(null);
      });

    return () => {
      cancelled = true;
    };
  }, [setCurrentUser, epoch]);

  // 子（App 以下のルート）は常に描画する。未ログインでも公開ルート（/login 等）は表示される
  return <>{children}</>;
}
