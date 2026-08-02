package services

import (
	"context"
	"fmt"

	"github.com/abdullahtnz/forum-starkdb/internal/models"
	"github.com/abdullahtnz/forum-starkdb/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	db *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetProfile(ctx context.Context, userID string) (*models.UserProfile, error) {
	var p models.UserProfile
	err := s.db.QueryRow(ctx,
		`SELECT id, username, email, bio, COALESCE(avatar_url, ''), is_admin, created_at, last_login FROM users WHERE id = $1`, userID,
	).Scan(&p.ID, &p.Username, &p.Email, &p.Bio, &p.AvatarURL, &p.IsAdmin, &p.CreatedAt, &p.LastLogin)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &p, nil
}

func (s *UserService) GetPublicProfile(ctx context.Context, userID string) (*models.UserPublic, error) {
	var p models.UserPublic
	err := s.db.QueryRow(ctx,
		`SELECT id, username, bio, COALESCE(avatar_url, ''), is_admin, created_at FROM users WHERE id = $1`, userID,
	).Scan(&p.ID, &p.Username, &p.Bio, &p.AvatarURL, &p.IsAdmin, &p.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &p, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID string, req models.UpdateProfileRequest) (*models.UserProfile, error) {
	username := utils.SanitizeUsername(req.Username)
	if len(username) < 3 || len(username) > 30 {
		return nil, fmt.Errorf("username must be between 3 and 30 characters")
	}
	_, err := s.db.Exec(ctx,
		`UPDATE users SET username = $1, bio = $2, updated_at = NOW() WHERE id = $3`,
		username, utils.SanitizeHTML(req.Bio), userID,
	)
	if err != nil {
		if isDuplicate(err) {
			return nil, fmt.Errorf("username already taken")
		}
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}
	return s.GetProfile(ctx, userID)
}

func (s *UserService) GetUserPosts(ctx context.Context, userID string, page, pageSize int) ([]models.PostDetail, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	rows, err := s.db.Query(ctx,
		`SELECT p.id, p.user_id, u.username, COALESCE(u.avatar_url, ''), p.title, substring(p.content for 200), p.view_count, p.is_pinned, p.is_closed, p.created_at, p.updated_at
		 FROM posts p JOIN users u ON p.user_id = u.id
		 WHERE p.user_id = $1 ORDER BY p.created_at DESC LIMIT $2 OFFSET $3`,
		userID, pageSize, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.PostDetail
	for rows.Next() {
		var d models.PostDetail
		if err := rows.Scan(&d.ID, &d.UserID, &d.Username, &d.AvatarURL, &d.Title, &d.Content, &d.ViewCount, &d.IsPinned, &d.IsClosed, &d.CreatedAt, &d.UpdatedAt); err != nil {
			continue
		}
		_ = s.db.QueryRow(ctx, `SELECT COUNT(*) FROM post_likes WHERE post_id = $1`, d.ID).Scan(&d.LikeCount)
		_ = s.db.QueryRow(ctx, `SELECT COUNT(*) FROM comments WHERE post_id = $1 AND is_deleted = FALSE`, d.ID).Scan(&d.CommentCount)
		d.IsOwner = true
		d.Tags = s.getPostTags(ctx, d.ID)
		posts = append(posts, d)
	}
	if posts == nil {
		posts = []models.PostDetail{}
	}
	return posts, nil
}

func (s *UserService) getPostTags(ctx context.Context, postID string) []models.Tag {
	rows, _ := s.db.Query(ctx,
		`SELECT t.id, t.name, t.slug, t.description FROM tags t JOIN post_tags pt ON t.id = pt.tag_id WHERE pt.post_id = $1`, postID,
	)
	if rows == nil {
		return []models.Tag{}
	}
	defer rows.Close()
	var tags []models.Tag
	for rows.Next() {
		var t models.Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug, &t.Description); err == nil {
			tags = append(tags, t)
		}
	}
	if tags == nil {
		tags = []models.Tag{}
	}
	return tags
}
