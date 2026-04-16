import axios from 'axios';

// axios レスポンスエラーから文字列メッセージを取り出すヘルパー
export function getErrorMessage(err: unknown, fallback = 'エラーが発生しました'): string {
  if (axios.isAxiosError(err)) {
    const data = err.response?.data;
    if (typeof data?.error === 'string') return data.error;
    if (Array.isArray(data?.errors)) return data.errors.join(', ');
  }
  return fallback;
}

// axios レスポンスエラーからエラー配列を取り出すヘルパー
export function getErrorList(err: unknown, fallback = 'エラーが発生しました'): string[] {
  if (axios.isAxiosError(err)) {
    const data = err.response?.data;
    if (Array.isArray(data?.errors)) return data.errors;
    if (typeof data?.error === 'string') return [data.error];
  }
  return [fallback];
}
