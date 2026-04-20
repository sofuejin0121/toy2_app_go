/**
 * 新規投稿フォーム（ホーム・返信スレッドで再利用）。
 *
 * - FormData で本文・任意の画像・返信先 in_reply_to_id を POST /microposts へ送る。
 * - 成功時は onCreated で親に Micropost を渡し、親がリスト先頭に追加する等する。
 * - 未ログインなら null を返して何も出さない（ホームのゲスト表示では非表示）。
 */
import { useAtom } from 'jotai';
import { useEffect, useRef, useState } from 'react';
import { createMicropost } from '../api/client';
import { getErrorMessage } from '../api/errors';
import ErrorMessage from './ErrorMessage';
import { currentUserAtom } from '../store/auth';
import type { Micropost } from '../types';

interface Props {
  inReplyToId?: number;
  onCreated?: (post: Micropost) => void;
  placeholder?: string;
}

export default function MicropostForm({
  inReplyToId,
  onCreated,
  placeholder = '今何してる？',
}: Props) {
  const [currentUser] = useAtom(currentUserAtom);
  const [content, setContent] = useState('');
  const [imageFile, setImageFile] = useState<File | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [imagePreviewUrl, setImagePreviewUrl] = useState<string | null>(null);
  // type="file" は React state で中身をリセットしづらいので、送信成功後に DOM の value を空にして同じファイルを選べるようにする。
  const fileRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (!imageFile) {
      setImagePreviewUrl(null);
      return;
    }
    const url = URL.createObjectURL(imageFile);
    setImagePreviewUrl(url);
    return () => URL.revokeObjectURL(url);
  }, [imageFile]);

  if (!currentUser) return null;

  const clearImage = () => {
    setImageFile(null);
    if (fileRef.current) fileRef.current.value = '';
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!content.trim() || loading) return;
    setLoading(true);
    setError('');
    try {
      const fd = new FormData();
      fd.append('content', content.trim());
      if (inReplyToId) fd.append('in_reply_to_id', String(inReplyToId));
      if (imageFile) fd.append('image', imageFile);
      const post = await createMicropost(fd);
      setContent('');
      setImageFile(null);
      if (fileRef.current) fileRef.current.value = '';
      onCreated?.(post);
    } catch (err: unknown) {
      setError(getErrorMessage(err));
    }
    setLoading(false);
  };

  return (
    <form onSubmit={handleSubmit} className="bg-white border border-gray-200 rounded-xl p-4">
      <div className="flex gap-3">
        <img src={currentUser.avatar_url} alt="" className="w-10 h-10 rounded-full shrink-0" />
        <div className="flex-1">
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            placeholder={placeholder}
            maxLength={140}
            rows={3}
            className="w-full resize-none border-0 outline-none text-gray-900 placeholder-gray-400 text-sm"
          />
          {imagePreviewUrl && imageFile && (
            <div className="relative mt-3 inline-block max-w-full">
              <img
                src={imagePreviewUrl}
                alt={imageFile.name}
                className="max-h-48 max-w-full rounded-lg border border-gray-200 object-contain"
              />
              <button
                type="button"
                onClick={clearImage}
                className="absolute top-1 right-1 rounded-full bg-black/50 p-1 text-white hover:bg-black/70"
                aria-label="画像を外す"
              >
                <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
            </div>
          )}
          {error && <ErrorMessage message={error} variant="inline" />}
          <div className="flex items-center justify-between mt-2 pt-2 border-t border-gray-100">
            <div className="flex items-center gap-2">
              <label className="cursor-pointer text-blue-500 hover:text-blue-600">
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                  />
                </svg>
                <input
                  ref={fileRef}
                  type="file"
                  accept="image/jpeg,image/png,image/gif"
                  className="hidden"
                  onChange={(e) => setImageFile(e.target.files?.[0] || null)}
                />
              </label>
              {imageFile && (
                <span className="text-xs text-gray-500 truncate max-w-[10rem]" title={imageFile.name}>
                  {imageFile.name}
                </span>
              )}
            </div>
            <div className="flex items-center gap-3">
              <span
                className={`text-xs ${content.length > 120 ? 'text-orange-500' : 'text-gray-400'}`}
              >
                {140 - content.length}
              </span>
              <button
                type="submit"
                disabled={!content.trim() || loading}
                className="bg-blue-600 text-white text-sm px-4 py-1.5 rounded-full hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {loading ? '送信中...' : '投稿'}
              </button>
            </div>
          </div>
        </div>
      </div>
    </form>
  );
}
