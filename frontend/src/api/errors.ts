import axios from 'axios';
import type { ApiError } from '../types';

// axios レスポンスエラーから文字列メッセージを取り出すヘルパー
// バックエンドは {"error":"..."} または {"errors":["...",...]} を返す
export function getErrorMessage(err: unknown, fallback = 'エラーが発生しました'): string {
  if (axios.isAxiosError(err)) {
    const data = err.response?.data as ApiError | undefined;
    if (typeof data?.error === 'string') return data.error;
    if (Array.isArray(data?.errors)) return data.errors.join(', ');
  }
  return fallback;
}

// axios レスポンスエラーからエラー配列を取り出すヘルパー（フォームバリデーションエラー表示用）
export function getErrorList(err: unknown, fallback = 'エラーが発生しました'): string[] {
  if (axios.isAxiosError(err)) {
    const data = err.response?.data as ApiError | undefined;
    if (Array.isArray(data?.errors)) return data.errors;
    if (typeof data?.error === 'string') return [data.error];
  }
  return [fallback];
}
