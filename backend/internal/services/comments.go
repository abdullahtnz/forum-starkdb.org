package services

import (
	"context"
	"fmt"

	"github.com/abdullahtnz/forum-starkdb/internal/models"
	"github.com/abdullahtnz/forum-starkdb/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentService struct {
	db *pgxpool.Pool
}

func NewCommentService(db *pgxpool.Pool) *CommentService {
	return &CommentService{db: db}
}

func (s *CommentService) Create(ctx context.Context, postID, userID string, req models.CreateCommentRequest) (*models.CommentDetail, error) {
	content := utils.FilterBadWords(utils.SanitizeHTML(req.Content), nil)
	if len(content) == 0 {
		return nil, fmt.Errorf("comment cannot be empty")
	}
	var closed bool
	_ = s.db.QueryRow(ctx, `SELECT is_closed FROM posts WHERE id = $1`, postID).Scan(&closed)
	if closed {
		return nil, fmt.Errorf("this post is closed for comments")
	}
	var commentID string
	err := s.db.QueryRow(ctx,
		`INSERT INTO comments (post_id, user_id, parent_id, content) VALUES ($1, $2, $3, $4) RETURNING id`,
		postID, userID, req.ParentID, content,
	).Scan(&commentID)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}
	return s.GetByID(ctx, commentID, userID)
}

func (s *CommentService) GetByPost(ctx context.Context, postID, viewerID string) ([]models.CommentDetail, error) {
	rows, err := s.db.Query(ctx,
		`SELECT c.id, c.post_id, c.user_id, c.parent_id, c.content, c.is_deleted, c.created_at, c.updated_at, u.username, u.avatar_url
		 FROM comments c JOIN users u ON c.user_id = u.id
		 WHERE c.post_id = $1 AND c.parent_id IS NULL AND c.is_deleted = FALSE
		 ORDER BY c.created_at ASC`, postID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	var comments []models.CommentDetail
	for rows.Next() {
		var d models.CommentDetail
		var av *string
		if err := rows.Scan(&d.ID, &d.PostID, &d.UserID, &d.ParentID, &d.Content, &d.IsDeleted, &d.CreatedAt, &d.UpdatedAt, &d.Username, &av); err != nil {
			continue
		}
		if av != nil {
			d.AvatarURL = *av
		}
		_ = s.db.QueryRow(ctx, `SELECT COUNT(*) FROM comment_likes WHERE comment_id = $1`, d.ID).Scan(&d.LikeCount)
		if viewerID != "" {
			var exists bool
			_ = s.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM comment_likes WHERE comment_id = $1 AND user_id = $2)`, d.ID, viewerID).Scan(&exists)
			d.Liked = exists
		}
		d.IsOwner = viewerID != "" && d.UserID == viewerID
		d.Replies = s.getReplies(ctx, d.ID, viewerID)
		comments = append(comments, d)
	}
	if comments == nil {
		comments = []models.CommentDetail{}
	}
	return comments, nil
}

func (s *CommentService) getReplies(ctx context.Context, parentID, viewerID string) []models.CommentDetail {
	rows, err := s.db.Query(ctx,
		`SELECT c.id, c.post_id, c.user_id, c.parent_id, c.content, c.is_deleted, c.created_at, c.updated_at, u.username, u.avatar_url
		 FROM comments c JOIN users u ON c.user_id = u.id
		 WHERE c.parent_id = $1 AND c.is_deleted = FALSE
		 ORDER BY c.created_at ASC`, parentID,
	)
	if err != nil {
		return []models.CommentDetail{}
	}
	defer rows.Close()
	var replies []models.CommentDetail
	for rows.Next() {
		var d models.CommentDetail
		var av *string
		if err := rows.Scan(&d.ID, &d.PostID, &d.UserID, &d.ParentID, &d.Content, &d.IsDeleted, &d.CreatedAt, &d.UpdatedAt, &d.Username, &av); err != nil {
			continue
		}
		if av != nil {
			d.AvatarURL = *av
		}
		_ = s.db.QueryRow(ctx, `SELECT COUNT(*) FROM comment_likes WHERE comment_id = $1`, d.ID).Scan(&d.LikeCount)
		if viewerID != "" {
			var exists bool
			_ = s.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM comment_likes WHERE comment_id = $1 AND user_id = $2)`, d.ID, viewerID).Scan(&exists)
			d.Liked = exists
		}
		d.IsOwner = viewerID != "" && d.UserID == viewerID
		d.Replies = s.getReplies(ctx, d.ID, viewerID)
		replies = append(replies, d)
	}
	if replies == nil {
		replies = []models.CommentDetail{}
	}
	return replies
}

func (s *CommentService) GetByID(ctx context.Context, commentID, viewerID string) (*models.CommentDetail, error) {
	var d models.CommentDetail
	var av *string
	err := s.db.QueryRow(ctx,
		`SELECT c.id, c.post_id, c.user_id, c.parent_id, c.content, c.is_deleted, c.created_at, c.updated_at, u.username, u.avatar_url
		 FROM comments c JOIN users u ON c.user_id = u.id WHERE c.id = $1`, commentID,
	).Scan(&d.ID, &d.PostID, &d.UserID, &d.ParentID, &d.Content, &d.IsDeleted, &d.CreatedAt, &d.UpdatedAt, &d.Username, &av)
	if err != nil {
		return nil, fmt.Errorf("comment not found: %w", err)
	}
	if av != nil {
		d.AvatarURL = *av
	}
	_ = s.db.QueryRow(ctx, `SELECT COUNT(*) FROM comment_likes WHERE comment_id = $1`, d.ID).Scan(&d.LikeCount)
	if viewerID != "" {
		var exists bool
		_ = s.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM comment_likes WHERE comment_id = $1 AND user_id = $2)`, d.ID, viewerID).Scan(&exists)
		d.Liked = exists
	}
	d.IsOwner = viewerID != "" && d.UserID == viewerID
	return &d, nil
}

func (s *CommentService) Update(ctx context.Context, commentID, userID string, req models.UpdateCommentRequest) (*models.CommentDetail, error) {
	content := utils.FilterBadWords(utils.SanitizeHTML(req.Content), nil)
	_, err := s.db.Exec(ctx,
		`UPDATE comments SET content = $1, updated_at = NOW() WHERE id = $2 AND user_id = $3`,
		content, commentID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}
	return s.GetByID(ctx, commentID, userID)
}

func (s *CommentService) Delete(ctx context.Context, commentID, userID string, isAdmin bool) error {
	query := `UPDATE comments SET is_deleted = TRUE, content = '[deleted]' WHERE id = $1`
	args := []interface{}{commentID}
	if !isAdmin {
		query += ` AND user_id = $2`
		args = append(args, userID)
	}
	tag, err := s.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("comment not found or not authorized")
	}
	return nil
}

func (s *CommentService) ToggleLike(ctx context.Context, commentID, userID string) (bool, error) {
	tag, err := s.db.Exec(ctx,
		`DELETE FROM comment_likes WHERE comment_id = $1 AND user_id = $2`, commentID, userID,
	)
	if err != nil {
		return false, err
	}
	if tag.RowsAffected() > 0 {
		return false, nil
	}
	_, err = s.db.Exec(ctx,
		`INSERT INTO comment_likes (user_id, comment_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, userID, commentID,
	)
	return true, err
}
