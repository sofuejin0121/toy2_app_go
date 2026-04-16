import { Link } from 'react-router-dom';
import type { UserStatSummary } from '../types';

interface Props {
  profile: UserStatSummary;
}

export default function UserStatBar({ profile }: Props) {
  const {
    user,
    micropost_count,
    following_count,
    followers_count,
    liked_count,
    bookmark_count,
    is_current_user,
  } = profile;

  return (
    <div className="flex gap-4 text-sm flex-wrap">
      <Link to={`/users/${user.id}`} className="text-center hover:text-blue-600">
        <div className="font-bold text-gray-900">{micropost_count}</div>
        <div className="text-gray-500 text-xs">microposts</div>
      </Link>
      <Link to={`/users/${user.id}/following`} className="text-center hover:text-blue-600">
        <div className="font-bold text-gray-900">{following_count}</div>
        <div className="text-gray-500 text-xs">following</div>
      </Link>
      <Link to={`/users/${user.id}/followers`} className="text-center hover:text-blue-600">
        <div className="font-bold text-gray-900">{followers_count}</div>
        <div className="text-gray-500 text-xs">followers</div>
      </Link>
      <Link to={`/users/${user.id}/likes`} className="text-center hover:text-blue-600">
        <div className="font-bold text-gray-900">{liked_count}</div>
        <div className="text-gray-500 text-xs">likes</div>
      </Link>
      {is_current_user && (
        <Link to={`/users/${user.id}/bookmarks`} className="text-center hover:text-blue-600">
          <div className="font-bold text-gray-900">{bookmark_count}</div>
          <div className="text-gray-500 text-xs">bookmarks</div>
        </Link>
      )}
    </div>
  );
}
