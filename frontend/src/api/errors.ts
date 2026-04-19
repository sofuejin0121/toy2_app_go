/**
 * axios が投げるエラーから、画面向けのメッセージ文字列を取り出す層。
 *
 * バックエンドが返す JSON の形（例）:
 * - { "error": "単一のメッセージ" }
 * - { "errors": ["バリデーション1", "バリデーション2"] }
 *
 * `as` で無理に型付けせず、parseApiErrorPayload で unknown を安全に読む。
 * フォーム全体のエラーは getErrorList、トーストや 1 行なら getErrorMessage を使う。
 */
import axios from 'axios';
import type { ApiError } from '../types';

/** オブジェクトかどうかの型ガード（response.data が想定外の形でも落ちにくくする） */
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

/** 1 本の文字列として表示したいとき（ログイン失敗・トースト等） */
export function getErrorMessage(err: unknown, fallback = 'エラーが発生しました'): string {
  if (axios.isAxiosError(err)) {
    const data = parseApiErrorPayload(err.response?.data);
    if (typeof data?.error === 'string') return data.error;
    if (Array.isArray(data?.errors)) return data.errors.join(', ');
  }
  return fallback;
}

/** フォームの箇条書きエラー表示用（SignupPage 等）。配列で返す */
export function getErrorList(err: unknown, fallback = 'エラーが発生しました'): string[] {
  if (axios.isAxiosError(err)) {
    const data = parseApiErrorPayload(err.response?.data);
    if (Array.isArray(data?.errors)) return data.errors;
    if (typeof data?.error === 'string') return [data.error];
  }
  return [fallback];
}
