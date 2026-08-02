export interface User {
  id: string;
  username: string;
  email: string;
  bio: string;
  avatar_url: string;
  is_admin: boolean;
  created_at: string;
  last_login?: string;
}

export interface UserPublic {
  id: string;
  username: string;
  bio: string;
  avatar_url: string;
  is_admin: boolean;
  created_at: string;
}

export interface Post {
  id: string;
  user_id: string;
  title: string;
  content: string;
  view_count: number;
  is_pinned: boolean;
  is_closed: boolean;
  created_at: string;
  updated_at: string;
  username: string;
  avatar_url: string;
  like_count: number;
  comment_count: number;
  liked: boolean;
  tags: Tag[];
  attachments: Attachment[];
  is_owner: boolean;
}

export interface Tag {
  id: string;
  name: string;
  slug: string;
  description: string;
}

export interface Attachment {
  id: string;
  post_id: string;
  file_name: string;
  file_path: string;
  file_size: number;
  mime_type: string;
}

export interface Comment {
  id: string;
  post_id: string;
  user_id: string;
  parent_id: string | null;
  content: string;
  is_deleted: boolean;
  created_at: string;
  updated_at: string;
  username: string;
  avatar_url: string;
  like_count: number;
  liked: boolean;
  is_owner: boolean;
  replies: Comment[];
}

export interface PostListResponse {
  posts: Post[];
  total_count: number;
  page: number;
  page_size: number;
}

export interface AuthTokens {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export interface ApiError {
  error: string;
  code?: string;
}
