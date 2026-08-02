package models

import "time"

type Comment struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	ParentID  *string   `json:"parent_id"`
	Content   string    `json:"content"`
	IsDeleted bool      `json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommentDetail struct {
	Comment
	Username  string          `json:"username"`
	AvatarURL string          `json:"avatar_url"`
	LikeCount int             `json:"like_count"`
	Liked     bool            `json:"liked"`
	IsOwner   bool            `json:"is_owner"`
	Replies   []CommentDetail `json:"replies"`
}

type CreateCommentRequest struct {
	Content  string  `json:"content"`
	ParentID *string `json:"parent_id"`
}

type UpdateCommentRequest struct {
	Content string `json:"content"`
}
