import { useAtom } from 'jotai';
import { Link, useParams } from 'react-router-dom';
import Layout from '../components/Layout';
import LoadingSpinner from '../components/LoadingSpinner';
import MicropostCard from '../components/MicropostCard';
import MicropostForm from '../components/MicropostForm';
import { useDocumentTitle } from '../hooks/useDocumentTitle';
import { useMicropostThread } from '../hooks/useMicropostThread';
import { micropostTabTitle } from '../lib/documentTitle';
import { currentUserAtom } from '../store/auth';

/**
 * 1 投稿＋返信スレッド（/microposts/:id）。useMicropostThread で取得、返信は MicropostForm + setReplies。
 */
export default function MicropostPage() {
  const { id } = useParams<{ id: string }>();
  const [currentUser] = useAtom(currentUserAtom);

  // 投稿本文とリプライ一覧を取得するカスタムフック
  const { post, setPost, replies, setReplies, loading } = useMicropostThread(id);

  const docTitle = loading || !post ? '投稿' : micropostTabTitle(post.content);
  useDocumentTitle(docTitle);

  if (loading)
    return (
      <Layout>
        <LoadingSpinner />
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

        <MicropostCard post={post} onDelete={() => window.history.back()} onUpdate={setPost} />

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
