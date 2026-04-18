import axios from 'axios';
import type { ApiError } from '../types';

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null;
}

/** axios の response.data を API エラー JSON として読み取る（型アサーションなし） */
function parseApiErrorPayload(data: unknown): ApiError | undefined {
  if (!isRecord(data)) return undefined;
  const out: ApiError = {};
  if (typeof data.error === 'string') out.error = data.error;
  if (Array.isArray(data.errors)) {
    const strings: string[] = [];
    for (const item of data.errors) {
      if (typeof item !== 'string') break;
      strings.push(item);
    }
    if (strings.length === data.errors.length) out.errors = strings;
  }
  return out;
}

// axios レスポンスエラーから文字列メッセージを取り出すヘルパー
// バックエンドは {"error":"..."} または {"errors":["...",...]} を返す
export function getErrorMessage(err: unknown, fallback = 'エラーが発生しました'): string {
  if (axios.isAxiosError(err)) {
    const data = parseApiErrorPayload(err.response?.data);
    if (typeof data?.error === 'string') return data.error;
    if (Array.isArray(data?.errors)) return data.errors.join(', ');
  }
  return fallback;
}

// axios レスポンスエラーからエラー配列を取り出すヘルパー（フォームバリデーションエラー表示用）
export function getErrorList(err: unknown, fallback = 'エラーが発生しました'): string[] {
  if (axios.isAxiosError(err)) {
    const data = parseApiErrorPayload(err.response?.data);
    if (Array.isArray(data?.errors)) return data.errors;
    if (typeof data?.error === 'string') return [data.error];
  }
  return [fallback];
}
