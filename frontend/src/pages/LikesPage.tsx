import { useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { getUserLikes } from '../api/client';
import Layout from '../components/Layout';
import MicropostCard from '../components/MicropostCard';
import Pagination from '../components/Pagination';
import UserStatBar from '../components/UserStatBar';
import type { Micropost, Pagination as PaginationType, UserProfile } from '../types';

export default function LikesPage() {
  const { id } = useParams<{ id: string }>();
  const [posts, setPosts] = useState<Micropost[]>([]);
  const [pagination, setPagination] = useState<PaginationType | null>(null);
  const [profileData, setProfileData] = useState<Partial<UserProfile> | null>(null);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id) return;
    setLoading(true);
    getUserLikes(Number(id), page)
      .then((data) => {
        setPosts(data.microposts);
        setPagination(data.pagination);
        setProfileData({
          user: data.user,
          micropost_count: data.micropost_count,
          following_count: data.following_count,
          followers_count: data.followers_count,
          liked_count: data.liked_count,
          bookmark_count: 0,
          is_current_user: false,
          is_following: false,
          microposts: [],
          pagination: data.pagination,
        });
      })
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [id, page]);

  return (
    <Layout>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {profileData?.user && (
          <aside className="md:col-span-1">
            <div className="bg-white rounded-xl border border-gray-200 p-4 sticky top-20">
              <div className="text-center mb-4">
                <img
                  src={profileData.user.avatar_url}
                  alt={profileData.user.name}
                  className="w-16 h-16 rounded-full mx-auto mb-2"
                />
                <Link
                  to={`/users/${profileData.user.id}`}
                  className="font-bold text-gray-900 hover:underline"
                >
                  {profileData.user.name}
                </Link>
              </div>
              {(profileData as UserProfile) && <UserStatBar profile={profileData as UserProfile} />}
            </div>
          </aside>
        )}

        <div className="md:col-span-2 space-y-4">
          <h2 className="text-xl font-bold text-gray-900">いいねした投稿</h2>
          {loading ? (
            <div className="text-center py-10 text-gray-400">読み込み中...</div>
          ) : posts.length === 0 ? (
            <div className="text-center py-10 text-gray-400">まだいいねした投稿がありません</div>
          ) : (
            <>
              {posts.map((post) => (
                <MicropostCard
                  key={post.id}
                  post={post}
                  onDelete={(pid) => setPosts((prev) => prev.filter((p) => p.id !== pid))}
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
