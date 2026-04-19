import { atom } from 'jotai';
import type { User } from '../types';

// ログイン中のユーザーを保持するatom
// null = 未ログイン, undefined = まだ確認中（ローディング）
export const currentUserAtom = atom<User | null | undefined>(undefined);

/**
 * セッション境界を表すカウンタ。ログイン／ログアウト／有効化などで進める。
 * AuthLoader の getMe が「古いリクエスト」になったら結果を atom に反映しないために使う。
 */
export const authBootstrapEpochAtom = atom(0);
