/**
 * ページ単位の「まだデータを取っています」用のシンプルなスピナー。
 *
 * 使い方の典型:
 *   if (loading) return <Layout><LoadingSpinner /></Layout>;
 *
 * 認証の「確認中」（currentUser === undefined）は AuthLoader 経由で各所が同様に扱う。
 * このコンポーネントはあくまで UI 部品で、データフェッチのロジックは持たない。
 */
export default function LoadingSpinner() {
  return (
    <div className="min-h-[40vh] flex items-center justify-center">
      <div className="animate-spin w-8 h-8 border-2 border-blue-600 border-t-transparent rounded-full" />
    </div>
  );
}
