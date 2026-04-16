import { atom } from 'jotai';
import type { User } from '../types';

// ログイン中のユーザーを保持するatom
// null = 未ログイン, undefined = まだ確認中（ローディング）
export const currentUserAtom = atom<User | null | undefined>(undefined);
