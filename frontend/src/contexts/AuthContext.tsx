// このファイルはアプリ起動時に一度だけ /api/me を叩いて
// currentUserAtom に値をセットするコンポーネントです。
// 認証状態の読み書きは jotai の useAtom/useSetAtom を直接使ってください。
//
// 使い方:
//   import { useAtom } from 'jotai';
//   import { currentUserAtom } from '../store/auth';
//
//   const [currentUser, setCurrentUser] = useAtom(currentUserAtom);

import { useSetAtom } from 'jotai';
import { type ReactNode, useEffect } from 'react';
import { getMe } from '../api/client';
import { currentUserAtom } from '../store/auth';

export function AuthLoader({ children }: { children: ReactNode }) {
  const setCurrentUser = useSetAtom(currentUserAtom);

  useEffect(() => {
    getMe()
      .then((user) => setCurrentUser(user))
      .catch(() => setCurrentUser(null));
  }, [setCurrentUser]);

  return <>{children}</>;
}
