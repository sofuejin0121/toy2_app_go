export interface User {
  id: number;
  name: string;
  email: string;
  bio: string;
  admin: boolean;
  activated: boolean;
  avatar_url: string;
  created_at: string;
}

export interface Micropost {
  id: number;
  content: string;
  user_id: number;
  image_path?: string;
  in_reply_to_id?: number;
  like_count: number;
  is_liked: boolean;
  is_bookmarked: boolean;
  user: User;
  parent?: {
    id: number;
    content: string;
    user: User;
  };
  created_at: string;
}

export interface Pagination {
  current_page: number;
  total_pages: number;
  total_items: number;
  per_page: number;
  has_prev: boolean;
  has_next: boolean;
}

// GET /api/users/:id のレスポンス型（UserStatSummary を拡張）
export interface UserProfile extends UserStatSummary {
  is_following: boolean;
  relationship_id?: number;
  microposts: Micropost[];
  pagination: Pagination;
}

export interface Notification {
  id: number;
  action_type: string;
  read: boolean;
  actor: User;
  target_id?: number;
  target_content?: string;
  created_at: string;
}

export interface AdminStats {
  total_users: number;
  total_posts: number;
  today_signups: number;
  daily_signups: { date: string; count: number }[];
}

export interface Settings {
  email_on_follow: boolean;
  email_on_like: boolean;
}

export interface ApiError {
  error?: string;
  errors?: string[];
}

// UserStatBar が必要とする最小限の統計情報
// GET /api/users/:id/likes など UserProfile 全体が返らない API でも使えるようにする
export interface UserStatSummary {
  user: User;
  micropost_count: number;
  following_count: number;
  followers_count: number;
  liked_count: number;
  // is_current_user=true のときのみブックマーク数リンクを表示する
  is_current_user: boolean;
  bookmark_count: number;
}
