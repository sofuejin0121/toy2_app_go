import { useSetAtom } from 'jotai';
import { useEffect, useState } from 'react';
import { useNavigate, useParams, useSearchParams } from 'react-router-dom';
import { activateAccount } from '../api/client';
import { getErrorMessage } from '../api/errors';
import ErrorMessage from '../components/ErrorMessage';
import Layout from '../components/Layout';
import { currentUserAtom } from '../store/auth';

export default function AccountActivationPage() {
  const { token } = useParams<{ token: string }>();
  const [searchParams] = useSearchParams();
  const email = searchParams.get('email') || '';
  const setCurrentUser = useSetAtom(currentUserAtom);
  const navigate = useNavigate();
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');
  const [message, setMessage] = useState('');

  useEffect(() => {
    if (!token || !email) {
      setStatus('error');
      setMessage('無効なリンクです');
      return;
    }
    activateAccount(token, email)
      .then((data) => {
        setCurrentUser(data.user);
        setStatus('success');
        setMessage(data.message);
        setTimeout(() => navigate(`/users/${data.user.id}`), 2000);
      })
      .catch((err: unknown) => {
        setStatus('error');
        setMessage(getErrorMessage(err, '有効化に失敗しました'));
      });
  }, [token, email, setCurrentUser, navigate]);

  return (
    <Layout>
      <div className="max-w-md mx-auto text-center py-16">
        {status === 'loading' && (
          <div className="text-gray-500">
            <div className="animate-spin w-8 h-8 border-2 border-blue-600 border-t-transparent rounded-full mx-auto mb-4" />
            アカウントを有効化中...
          </div>
        )}
        {status === 'success' && (
          <div className="bg-green-50 border border-green-200 rounded-xl p-8">
            <svg
              className="w-12 h-12 text-green-500 mx-auto mb-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <h2 className="text-xl font-bold text-gray-900 mb-2">有効化完了！</h2>
            <p className="text-gray-600 text-sm">{message}</p>
            <p className="text-xs text-gray-400 mt-2">プロフィールページに移動します...</p>
          </div>
        )}
        {status === 'error' && (
          <div className="bg-red-50 border border-red-200 rounded-xl p-8">
            <h2 className="text-xl font-bold text-gray-900 mb-2">有効化失敗</h2>
            <ErrorMessage message={message} variant="inline" className="text-red-700 text-sm" />
            <button
              type="button"
              onClick={() => navigate('/')}
              className="mt-4 text-blue-600 hover:underline text-sm"
            >
              ホームへ
            </button>
          </div>
        )}
      </div>
    </Layout>
  );
}
