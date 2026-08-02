package services

import (
	"context"
	"fmt"

	"github.com/abdullahtnz/forum-starkdb/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TagService struct {
	db *pgxpool.Pool
}

func NewTagService(db *pgxpool.Pool) *TagService {
	return &TagService{db: db}
}

func (s *TagService) List(ctx context.Context) ([]models.Tag, error) {
	rows, err := s.db.Query(ctx, `SELECT id, name, slug, description, created_at FROM tags ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}
	defer rows.Close()
	var tags []models.Tag
	for rows.Next() {
		var t models.Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.Slug, &t.Description, &t.CreatedAt); err == nil {
			tags = append(tags, t)
		}
	}
	if tags == nil {
		tags = []models.Tag{}
	}
	return tags, nil
}

type UploadService struct {
	db        *pgxpool.Pool
	uploadDir string
	maxSize   int64
}

func NewUploadService(db *pgxpool.Pool, uploadDir string, maxSize int64) *UploadService {
	return &UploadService{db: db, uploadDir: uploadDir, maxSize: maxSize}
}

func (s *UploadService) GetUploadDir() string { return s.uploadDir }
func (s *UploadService) GetMaxSize() int64    { return s.maxSize }

func (s *UploadService) AttachToPost(ctx context.Context, postID, fileName, filePath, mimeType string, fileSize int64) (*models.Attachment, error) {
	var a models.Attachment
	err := s.db.QueryRow(ctx,
		`INSERT INTO post_attachments (post_id, file_name, file_path, file_size, mime_type) VALUES ($1, $2, $3, $4, $5) RETURNING id, post_id, file_name, file_path, file_size, mime_type, created_at`,
		postID, fileName, filePath, fileSize, mimeType,
	).Scan(&a.ID, &a.PostID, &a.FileName, &a.FilePath, &a.FileSize, &a.MimeType, &a.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to attach file: %w", err)
	}
	return &a, nil
}
