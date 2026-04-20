import { Link } from 'react-router-dom';
import Layout from '../components/Layout';
import { useDocumentTitle } from '../hooks/useDocumentTitle';

/**
 * 未定義ルート（App.tsx の path="*"）。存在しない URL へ来たことを伝え、ホームへ誘導する。
 */
export default function NotFoundPage() {
  useDocumentTitle('ページが見つかりません');

  return (
    <Layout>
      <div className="max-w-lg mx-auto text-center py-16 px-4">
        <p className="text-6xl font-bold text-gray-200 mb-2">404</p>
        <h1 className="text-xl font-semibold text-gray-900 mb-2">ページが見つかりません</h1>
        <p className="text-sm text-gray-500 mb-8">
          お探しのページは存在しないか、URL が間違っている可能性があります。
        </p>
        <Link
          to="/"
          className="inline-block bg-blue-600 text-white text-sm px-5 py-2 rounded-full hover:bg-blue-700"
        >
          ホームへ戻る
        </Link>
      </div>
    </Layout>
  );
}
