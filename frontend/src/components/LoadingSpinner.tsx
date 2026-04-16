// ページ全体のローディング表示
// 使い方: if (loading) return <LoadingSpinner />;
export default function LoadingSpinner() {
  return (
    <div className="min-h-[40vh] flex items-center justify-center">
      <div className="animate-spin w-8 h-8 border-2 border-blue-600 border-t-transparent rounded-full" />
    </div>
  );
}
