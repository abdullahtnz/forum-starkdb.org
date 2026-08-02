package models

import "time"

type Post struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ViewCount int       `json:"view_count"`
	IsPinned  bool      `json:"is_pinned"`
	IsClosed  bool      `json:"is_closed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostDetail struct {
	Post
	Username     string       `json:"username"`
	AvatarURL    string       `json:"avatar_url"`
	LikeCount    int          `json:"like_count"`
	CommentCount int          `json:"comment_count"`
	Liked        bool         `json:"liked"`
	Tags         []Tag        `json:"tags"`
	Attachments  []Attachment `json:"attachments"`
	IsOwner      bool         `json:"is_owner"`
}

type PostListResponse struct {
	Posts      []PostDetail `json:"posts"`
	TotalCount int          `json:"total_count"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
}

type CreatePostRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	TagIDs  []string `json:"tag_ids"`
}

type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
