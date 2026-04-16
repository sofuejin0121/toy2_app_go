import { useAtom } from 'jotai';
import { useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { getMicropost } from '../api/client';
import Layout from '../components/Layout';
import MicropostCard from '../components/MicropostCard';
import MicropostForm from '../components/MicropostForm';
import { currentUserAtom } from '../store/auth';
import type { Micropost } from '../types';

export default function MicropostPage() {
  const { id } = useParams<{ id: string }>();
  const [currentUser] = useAtom(currentUserAtom);
  const [post, setPost] = useState<Micropost | null>(null);
  const [replies, setReplies] = useState<Micropost[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id) return;
    setLoading(true);
    getMicropost(Number(id))
      .then((data) => {
        setPost(data.post);
        setReplies(data.replies);
      })
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [id]);

  if (loading)
    return (
      <Layout>
        <div className="text-center py-10 text-gray-400">読み込み中...</div>
      </Layout>
    );
  if (!post)
    return (
      <Layout>
        <div className="text-center py-10 text-gray-400">投稿が見つかりません</div>
      </Layout>
    );

  return (
    <Layout>
      <div className="max-w-2xl mx-auto space-y-4">
        <div className="flex items-center gap-2 mb-2">
          <Link to="/" className="text-blue-600 text-sm hover:underline">
            ← ホーム
          </Link>
        </div>

        <MicropostCard
          post={post}
          onDelete={() => window.history.back()}
          onUpdate={(updated) => setPost(updated)}
        />

        {currentUser && (
          <MicropostForm
            inReplyToId={post.id}
            placeholder={`${post.user.name} に返信...`}
            onCreated={(newReply) => setReplies((prev) => [newReply, ...prev])}
          />
        )}

        {replies.length > 0 && (
          <div className="space-y-3">
            <h3 className="text-sm font-medium text-gray-500">返信 {replies.length}件</h3>
            {replies.map((reply) => (
              <MicropostCard
                key={reply.id}
                post={reply}
                onDelete={(pid) => setReplies((prev) => prev.filter((r) => r.id !== pid))}
              />
            ))}
          </div>
        )}
      </div>
    </Layout>
  );
}
