/**
 * バックエンド API への HTTP 呼び出しをまとめたモジュール。
 *
 * - axios インスタンス `api` に baseURL・Cookie（withCredentials）・JSON ヘッダを設定。
 * - 各 export は「どのパスに何を送るか」と「レスポンスの型」を 1 関数に閉じ込める。
 * - 画面や SWR の fetcher はここを呼ぶだけにし、URL 文字列の重複を避ける。
 *
 * 環境:
 * - VITE_API_BASE_URL があれば `${それ}/api`、無ければ同一オリジンの `/api`（プロキシ前提の開発など）。
 */
import axios from 'axios';
import type {
  AdminStats,
  Micropost,
  Notification,
  Pagination,
  Settings,
  User,
  UserProfile,
} from '../types';

// Netlify デプロイ時は VITE_API_BASE_URL に Render の URL を設定する
// 例: https://your-app.onrender.com
// 未設定の場合は同一オリジン（Render のみ or ローカル開発）
const baseURL = import.meta.env.VITE_API_BASE_URL
  ? `${import.meta.env.VITE_API_BASE_URL}/api`
  : '/api';

// withCredentials: true … セッション Cookie を付与する（ログイン状態の維持に必須）
const api = axios.create({
  baseURL,
  withCredentials: true,
  headers: { 'Content-Type': 'application/json' },
});

// ---- 認証 ----
export const getMe = () => api.get<User>('/me').then((r) => r.data);

export const login = (email: string, password: string, remember = false) =>
  api.post<User>('/login', { email, password, remember }).then((r) => r.data);

export const logout = () => api.delete('/logout').then((r) => r.data);

// ---- ユーザー ----
export const signUp = (data: {
  name: string;
  email: string;
  password: string;
  password_confirmation: string;
}) => api.post<User>('/users', data).then((r) => r.data);

export const listUsers = (page = 1, q = '') =>
  api
    .get<{ users: User[]; pagination: Pagination }>('/users', { params: { page, q } })
    .then((r) => r.data);

export const getUser = (id: number, page = 1) =>
  api.get<UserProfile>(`/users/${id}`, { params: { page } }).then((r) => r.data);

export const updateUser = (
  id: number,
  data: {
    name: string;
    email: string;
    bio: string;
    password?: string;
    password_confirmation?: string;
  },
) => api.patch<User>(`/users/${id}`, data).then((r) => r.data);

export const deleteUser = (id: number) => api.delete(`/users/${id}`).then((r) => r.data);

export const getFollowing = (id: number, page = 1) =>
  api
    .get<{
      user: User;
      users: User[];
      following_count: number;
      followers_count: number;
      micropost_count: number;
      liked_count: number;
      bookmark_count: number;
      is_current_user: boolean;
      pagination: Pagination;
    }>(`/users/${id}/following`, { params: { page } })
    .then((r) => r.data);

export const getFollowers = (id: number, page = 1) =>
  api
    .get<{
      user: User;
      users: User[];
      following_count: number;
      followers_count: number;
      micropost_count: number;
      liked_count: number;
      bookmark_count: number;
      is_current_user: boolean;
      pagination: Pagination;
    }>(`/users/${id}/followers`, { params: { page } })
    .then((r) => r.data);

export const getUserLikes = (id: number, page = 1) =>
  api
    .get<{
      user: User;
      microposts: Micropost[];
      liked_count: number;
      bookmark_count: number;
      is_current_user: boolean;
      following_count: number;
      followers_count: number;
      micropost_count: number;
      pagination: Pagination;
    }>(`/users/${id}/likes`, { params: { page } })
    .then((r) => r.data);

export const getUserBookmarks = (id: number, page = 1) =>
  api
    .get<{
      user: User;
      microposts: Micropost[];
      bookmark_count: number;
      following_count: number;
      followers_count: number;
      micropost_count: number;
      liked_count: number;
      pagination: Pagination;
    }>(`/users/${id}/bookmarks`, { params: { page } })
    .then((r) => r.data);

// ---- フィード ----
export const getFeed = (page = 1) =>
  api
    .get<{ items: Micropost[]; pagination: Pagination }>('/feed', { params: { page } })
    .then((r) => r.data);

// ---- マイクロポスト ----
export const getMicropost = (id: number) =>
  api
    .get<{ post: Micropost; replies: Micropost[]; reply_count: number }>(`/microposts/${id}`)
    .then((r) => r.data);

export const createMicropost = (formData: FormData) =>
  api
    .post<Micropost>('/microposts', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    .then((r) => r.data);

export const deleteMicropost = (id: number) => api.delete(`/microposts/${id}`).then((r) => r.data);

// ---- フォロー ----
export const follow = (followedId: number) =>
  api
    .post<{ message: string; relationship_id: number }>('/relationships', {
      followed_id: followedId,
    })
    .then((r) => r.data);

export const unfollow = (relationshipId: number) =>
  api.delete(`/relationships/${relationshipId}`).then((r) => r.data);

// ---- いいね ----
export const like = (micropostId: number) =>
  api
    .post<{ liked: boolean; count: number }>('/likes', { micropost_id: micropostId })
    .then((r) => r.data);

export const unlike = (micropostId: number) =>
  api.delete<{ liked: boolean; count: number }>(`/likes/${micropostId}`).then((r) => r.data);

// ---- ブックマーク ----
export const bookmark = (micropostId: number) =>
  api
    .post<{ bookmarked: boolean }>('/bookmarks', { micropost_id: micropostId })
    .then((r) => r.data);

export const unbookmark = (micropostId: number) =>
  api.delete<{ bookmarked: boolean }>(`/bookmarks/${micropostId}`).then((r) => r.data);

// ---- 通知 ----
export const listNotifications = () =>
  api.get<{ notifications: Notification[] }>('/notifications').then((r) => r.data);

export const getUnreadCount = () =>
  api.get<{ count: number }>('/notifications/unread_count').then((r) => r.data);

export const deleteNotification = (id: number) =>
  api.delete(`/notifications/${id}`).then((r) => r.data);

// ---- 管理者 ----
export const getAdminStats = () => api.get<AdminStats>('/admin').then((r) => r.data);

// ---- 設定 ----
export const getSettings = () => api.get<Settings>('/settings').then((r) => r.data);

export const updateSettings = (data: Settings) =>
  api.patch<Settings>('/settings', data).then((r) => r.data);

// ---- アカウント有効化 ----
export const activateAccount = (token: string, email: string) =>
  api
    .get<{ message: string; user: User }>(`/account_activations/${token}/edit`, {
      params: { email },
    })
    .then((r) => r.data);

// ---- パスワードリセット ----
export const requestPasswordReset = (email: string) =>
  api.post('/password_resets', { email }).then((r) => r.data);

export const checkPasswordResetToken = (token: string, email: string) =>
  api
    .get<{ email: string; token: string }>(`/password_resets/${token}/edit`, { params: { email } })
    .then((r) => r.data);

export const resetPassword = (
  token: string,
  data: {
    email: string;
    password: string;
    password_confirmation: string;
  },
) =>
  api.patch<{ message: string; user: User }>(`/password_resets/${token}`, data).then((r) => r.data);

export default api;
