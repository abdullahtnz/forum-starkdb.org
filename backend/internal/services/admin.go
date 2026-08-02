package services

import (
	"context"

	"github.com/abdullahtnz/forum-starkdb/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminService struct {
	db *pgxpool.Pool
}

func NewAdminService(db *pgxpool.Pool) *AdminService {
	return &AdminService{db: db}
}

func (s *AdminService) GetUsers(ctx context.Context, page, pageSize int) ([]models.UserPublic, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}
	offset := (page - 1) * pageSize

	var total int
	_ = s.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total)

	rows, err := s.db.Query(ctx,
		`SELECT id, username, bio, COALESCE(avatar_url, ''), is_admin, created_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []models.UserPublic
	for rows.Next() {
		var u models.UserPublic
		if err := rows.Scan(&u.ID, &u.Username, &u.Bio, &u.AvatarURL, &u.IsAdmin, &u.CreatedAt); err == nil {
			users = append(users, u)
		}
	}
	if users == nil {
		users = []models.UserPublic{}
	}
	return users, total, nil
}

func (s *AdminService) DeleteUser(ctx context.Context, userID string) error {
	_, err := s.db.Exec(ctx, `DELETE FROM users WHERE id = $1 AND is_admin = FALSE`, userID)
	return err
}

func (s *AdminService) ToggleAdmin(ctx context.Context, userID string, makeAdmin bool) error {
	_, err := s.db.Exec(ctx, `UPDATE users SET is_admin = $1 WHERE id = $2`, makeAdmin, userID)
	return err
}

func (s *AdminService) AddBadWord(ctx context.Context, word string) error {
	_, err := s.db.Exec(ctx, `INSERT INTO bad_words (word) VALUES ($1) ON CONFLICT (word) DO NOTHING`, word)
	return err
}

func (s *AdminService) RemoveBadWord(ctx context.Context, word string) error {
	_, err := s.db.Exec(ctx, `DELETE FROM bad_words WHERE word = $1`, word)
	return err
}

func (s *AdminService) GetBadWords(ctx context.Context) ([]string, error) {
	rows, err := s.db.Query(ctx, `SELECT word FROM bad_words ORDER BY word`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var words []string
	for rows.Next() {
		var w string
		if err := rows.Scan(&w); err == nil {
			words = append(words, w)
		}
	}
	if words == nil {
		words = []string{}
	}
	return words, nil
}

func (s *AdminService) DeletePost(ctx context.Context, postID string) error {
	_, err := s.db.Exec(ctx, `DELETE FROM posts WHERE id = $1`, postID)
	return err
}

func (s *AdminService) DeleteComment(ctx context.Context, commentID string) error {
	_, err := s.db.Exec(ctx, `UPDATE comments SET is_deleted = TRUE, content = '[removed by moderator]' WHERE id = $1`, commentID)
	return err
}
