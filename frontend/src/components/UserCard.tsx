import { Link } from 'react-router-dom';
import type { User } from '../types';

interface Props {
  user: User;
  onDelete?: (id: number) => void;
  showAdmin?: boolean;
}

export default function UserCard({ user, onDelete, showAdmin }: Props) {
  return (
    <div className="flex items-center justify-between py-3 border-b border-gray-100 last:border-0">
      <div className="flex items-center gap-3">
        <Link to={`/users/${user.id}`}>
          <img
            src={user.avatar_url}
            alt={user.name}
            className="w-10 h-10 rounded-full object-cover"
          />
        </Link>
        <div>
          <Link
            to={`/users/${user.id}`}
            className="font-medium text-gray-900 hover:underline text-sm"
          >
            {user.name}
          </Link>
          {user.bio && <p className="text-xs text-gray-500 mt-0.5 max-w-xs truncate">{user.bio}</p>}
        </div>
      </div>
      {showAdmin && onDelete && (
        <button
          onClick={() => onDelete(user.id)}
          className="text-xs text-red-400 hover:text-red-600 border border-red-200 px-2 py-1 rounded"
        >
          削除
        </button>
      )}
    </div>
  );
}
