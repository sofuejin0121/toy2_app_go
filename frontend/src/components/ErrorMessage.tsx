/**
 * エラー文言の表示用プレゼンテーションコンポーネント。
 *
 * Props の使い分け:
 * - message … 1 件だけ出したいとき
 * - messages … 複数行（バリデーションの配列）。message より優先される
 * - variant … 'box'（枠付き・ページ上部向け） / 'inline'（フォーム直下の小さめテキスト）
 *
 * データ取得はしない。親が getErrorMessage / getErrorList の結果を渡す。
 */
interface Props {
  message?: string;
  messages?: string[];
  className?: string;
  /** フォーム内など、枠なしの1行表示 */
  variant?: 'box' | 'inline';
}

export default function ErrorMessage({
  message,
  messages,
  className = '',
  variant = 'box',
}: Props) {
  const lines = messages?.length ? messages : message ? [message] : [];
  if (lines.length === 0) return null;

  if (variant === 'inline') {
    return (
      <p className={`text-red-500 text-xs mt-1 ${className}`}>
        {lines.length === 1 ? lines[0] : lines.join(' / ')}
      </p>
    );
  }

  if (lines.length === 1) {
    return (
      <div
        className={`mb-4 p-3 bg-red-50 border border-red-200 text-red-700 text-sm rounded-lg ${className}`}
      >
        {lines[0]}
      </div>
    );
  }

  return (
    <div
      className={`mb-4 p-3 bg-red-50 border border-red-200 text-red-700 text-sm rounded-lg ${className}`}
    >
      <ul className="list-disc list-inside space-y-1">
        {lines.map((line) => (
          <li key={line}>{line}</li>
        ))}
      </ul>
    </div>
  );
}
