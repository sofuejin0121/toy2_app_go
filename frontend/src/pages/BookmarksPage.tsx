import { useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import ErrorMessage from '../components/ErrorMessage';
import Layout from '../components/Layout';
import LoadingSpinner from '../components/LoadingSpinner';
import MicropostCard from '../components/MicropostCard';
import Pagination from '../components/Pagination';
import UserStatBar from '../components/UserStatBar';
import { useUserBookmarks } from '../hooks/useUserBookmarks';
import { updateIfDefined } from '../utils/updateIfDefined';

export default function BookmarksPage() {
  const { id } = useParams<{ id: string }>();
  const [page, setPage] = useState(1);

  const { posts, statSummary, pagination, loading, error, mutate } = useUserBookmarks(id, page);

  return (
    <Layout>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {statSummary?.user && (
          <aside className="md:col-span-1">
            <div className="bg-white rounded-xl border border-gray-200 p-4 sticky top-20">
              <div className="text-center mb-4">
                <img
                  src={statSummary.user.avatar_url}
                  alt={statSummary.user.name}
                  className="w-16 h-16 rounded-full mx-auto mb-2"
                />
                <Link
                  to={`/users/${statSummary.user.id}`}
                  className="font-bold text-gray-900 hover:underline"
                >
                  {statSummary.user.name}
                </Link>
              </div>
              <UserStatBar profile={statSummary} />
            </div>
          </aside>
        )}

        <div className="md:col-span-2 space-y-4">
          <h2 className="text-xl font-bold text-gray-900">ブックマーク</h2>
          {error && <ErrorMessage message={error} />}
          {loading ? (
            <LoadingSpinner />
          ) : posts.length === 0 ? (
            <div className="text-center py-10 text-gray-400">まだブックマークがありません</div>
          ) : (
            <>
              {posts.map((post) => (
                <MicropostCard
                  key={post.id}
                  post={post}
                  onDelete={(pid) =>
                    mutate(
                      (prev) =>
                        updateIfDefined(prev, (p) => ({
                          ...p,
                          microposts: p.microposts.filter((x) => x.id !== pid),
                        })),
                      { revalidate: false },
                    )
                  }
                />
              ))}
              {pagination && <Pagination pagination={pagination} onPageChange={setPage} />}
            </>
          )}
        </div>
      </div>
    </Layout>
  );
}
