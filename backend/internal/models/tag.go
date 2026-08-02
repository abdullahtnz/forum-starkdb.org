package models

import (
	"time"
)

type Tag struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Attachment struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	FileName  string    `json:"file_name"`
	FilePath  string    `json:"file_path"`
	FileSize  int64     `json:"file_size"`
	MimeType  string    `json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
}

type PostLike struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	PostID    string    `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentLike struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	CommentID string    `json:"comment_id"`
	CreatedAt time.Time `json:"created_at"`
}
