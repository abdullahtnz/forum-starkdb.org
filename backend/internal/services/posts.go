package services

import (
	"context"
	"fmt"

	"github.com/abdullahtnz/forum-starkdb/internal/models"
	"github.com/abdullahtnz/forum-starkdb/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostService struct {
	db       *pgxpool.Pool
	badWords []string
}

func NewPostService(db *pgxpool.Pool) *PostService {
	s := &PostService{db: db}
	go s.loadBadWords()
	return s
}

func (s *PostService) loadBadWords() {
	ctx := context.Background()
	rows, err := s.db.Query(ctx, `SELECT word FROM bad_words`)
	if err != nil {
		return
	}
	defer rows.Close()
	var words []string
	for rows.Next() {
		var w string
		if err := rows.Scan(&w); err == nil {
			words = append(words, w)
		}
	}
	s.badWords = words
}

func (s *PostService) Create(ctx context.Context, userID string, req models.CreatePostRequest) (*models.PostDetail, error) {
	req.Title = utils.SanitizeHTML(req.Title)
	req.Content = utils.FilterBadWords(utils.SanitizeHTML(req.Content), s.badWords)
	if len(req.Title) < 5 {
		return nil, fmt.Errorf("title must be at least 5 characters")
	}
	var post models.Post
	err := s.db.QueryRow(ctx,
		`INSERT INTO posts (user_id, title, content) VALUES ($1, $2, $3) RETURNING id, user_id, title, content, view_count, is_pinned, is_closed, created_at, updated_at`,
		userID, req.Title, req.Content,
	).Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.ViewCount, &post.IsPinned, &post.IsClosed, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}
	for _, tagID := range req.TagIDs {
		_, _ = s.db.Exec(ctx, `INSERT INTO post_tags (post_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, post.ID, tagID)
	}
	return s.GetByID(ctx, post.ID, userID)
}

func (s *PostService) GetByID(ctx context.Context, postID, viewerID string) (*models.PostDetail, error) {
	var detail models.PostDetail
	var avatarURL *string
	err := s.db.QueryRow(ctx,
		`SELECT p.id, p.user_id, u.username, u.avatar_url, p.title, p.content, p.view_count, p.is_pinned, p.is_closed, p.created_at, p.updated_at
		 FROM posts p JOIN users u ON p.user_id = u.id WHERE p.id = $1`, postID,
	).Scan(
		&detail.ID, &detail.UserID, &detail.Username, &avatarURL,
		&detail.Title, &detail.Content, &detail.ViewCount, &detail.IsPinned, &detail.IsClosed,
		&detail.CreatedAt, &detail.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}
	if avatarURL != nil {
		detail.AvatarURL = *avatarURL
	}
	detail.IsOwner = viewerID != "" && detail.UserID == viewerID

	_ = s.db.QueryRow(ctx, `SELECT COUNT(*) FROM post_likes WHERE post_id = $1`, postID).Scan(&detail.LikeCount)
	_ = s.db.QueryRow(ctx, `SELECT COUNT(*) FROM comments WHERE post_id = $1 AND is_deleted = FALSE`, postID).Scan(&detail.CommentCount)
	if viewerID != "" {
		var exists bool
		_ = s.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM post_likes WHERE post_id = $1 AND user_id = $2)`, postID, viewerID).Scan(&exists)
		detail.Liked = exists
	}

	detail.Tags = s.getPostTags(ctx, postID)
	detail.Attachments = s.getPostAttachments(ctx, postID)
	_, _ = s.db.Exec(ctx, `UPDATE posts SET view_count = view_count + 1 WHERE id = $1`, postID)
	return &detail, nil
}

func (s *PostService) List(ctx context.Context, page, pageSize int, sort, tag, search, viewerID string) (*models.PostListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	baseQuery := `FROM posts p JOIN users u ON p.user_id = u.id`
	args := []interface{}{}
	argIdx := 1
	whereClause := ""

	if tag != "" {
		baseQuery += ` JOIN post_tags pt ON p.id = pt.post_id JOIN tags t ON pt.tag_id = t.id`
		whereClause = fmt.Sprintf(` WHERE t.slug = $%d`, argIdx)
		args = append(args, tag)
		argIdx++
	}
	if search != "" {
		if whereClause != "" {
			whereClause += ` AND`
		} else {
			whereClause = ` WHERE`
		}
		whereClause += fmt.Sprintf(` p.search_vector @@ plainto_tsquery('english', $%d)`, argIdx)
		args = append(args, search)
		argIdx++
	}

	var totalCount int
	countQuery := `SELECT COUNT(*) ` + baseQuery + whereClause
	if err := s.db.QueryRow(ctx, countQuery, args...).Scan(&totalCount); err != nil {
		return nil, fmt.Errorf("failed to count posts: %w", err)
	}

	orderClause := "ORDER BY p.is_pinned DESC, p.created_at DESC"
	if sort == "oldest" {
		orderClause = "ORDER BY p.is_pinned DESC, p.created_at ASC"
	} else if sort == "popular" {
		orderClause = "ORDER BY p.is_pinned DESC, p.view_count DESC"
	}

	dataQuery := fmt.Sprintf(`SELECT p.id, p.user_id, u.username, u.avatar_url, p.title, substring(p.content for 300) as content, p.view_count, p.is_pinned, p.is_closed, p.created_at, p.updated_at %s %s %s LIMIT $%d OFFSET $%d`,
		baseQuery, whereClause, orderClause, argIdx, argIdx+1)
	args = append(args, pageSize, offset)

	rows, err := s.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}
	defer rows.Close()

	var posts []models.PostDetail
	for rows.Next() {
		var d models.PostDetail
		var av *string
		if err := rows.Scan(&d.ID, &d.UserID, &d.Username, &av, &d.Title, &d.Content, &d.ViewCount, &d.IsPinned, &d.IsClosed, &d.CreatedAt, &d.UpdatedAt); err != nil {
			continue
		}
		if av != nil {
			d.AvatarURL = *av
		}
		_ = s.db.QueryRow(ctx, `SELECT COUNT(*) FROM post_likes WHERE post_id = $1`, d.ID).Scan(&d.LikeCount)
		_ = s.db.QueryRow(ctx, `SELECT COUNT(*) FROM comments WHERE post_id = $1 AND is_deleted = FALSE`, d.ID).Scan(&d.CommentCount)
		if viewerID != "" {
			var exists bool
			_ = s.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM post_likes WHERE post_id = $1 AND user_id = $2)`, d.ID, viewerID).Scan(&exists)
			d.Liked = exists
		}
		d.IsOwner = viewerID != "" && d.UserID == viewerID
		d.Tags = s.getPostTags(ctx, d.ID)
		posts = append(posts, d)
	}
	if posts == nil {
		posts = []models.PostDetail{}
	}

	return &models.PostListResponse{
		Posts:      posts,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (s *PostService) Update(ctx context.Context, postID, userID string, req models.UpdatePostRequest) (*models.PostDetail, error) {
	req.Title = utils.SanitizeHTML(req.Title)
	req.Content = utils.FilterBadWords(utils.SanitizeHTML(req.Content), s.badWords)
	_, err := s.db.Exec(ctx,
		`UPDATE posts SET title = $1, content = $2, updated_at = NOW() WHERE id = $3 AND user_id = $4`,
		req.Title, req.Content, postID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}
	return s.GetByID(ctx, postID, userID)
}

func (s *PostService) Delete(ctx context.Context, postID, userID string, isAdmin bool) error {
	query := `DELETE FROM posts WHERE id = $1`
	args := []interface{}{postID}
	if !isAdmin {
		query += ` AND user_id = $2`
		args = append(args, userID)
	}
	tag, err := s.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("post not found or not authorized")
	}
	return nil
}

func (s *PostService) ToggleLike(ctx context.Context, postID, userID string) (bool, error) {
	tag, err := s.db.Exec(ctx,
		`DELETE FROM post_likes WHERE post_id = $1 AND user_id = $2`, postID, userID,
	)
	if err != nil {
		return false, err
	}
	if tag.RowsAffected() > 0 {
		return false, nil
	}
	_, err = s.db.Exec(ctx,
		`INSERT INTO post_likes (user_id, post_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, userID, postID,
	)
	return true, err
}

func (s *PostService) PinPost(ctx context.Context, postID string, pinned bool) error {
	_, err := s.db.Exec(ctx, `UPDATE posts SET is_pinned = $1 WHERE id = $2`, pinned, postID)
	return err
}

func (s *PostService) ClosePost(ctx context.Context, postID string, closed bool) error {
	_, err := s.db.Exec(ctx, `UPDATE posts SET is_closed = $1 WHERE id = $2`, closed, postID)
	return err
}

func (s *PostService) getPostTags(ctx context.Context, postID string) []models.Tag {
	rows, err := s.db.Query(ctx,
		`SELECT t.id, t.name, t.slug, t.description, t.created_at FROM tags t JOIN post_tags pt ON t.id = pt.tag_id WHERE pt.post_id = $1`, postID,
	)
	if err != nil {
		return nil
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
	return tags
}

func (s *PostService) getPostAttachments(ctx context.Context, postID string) []models.Attachment {
	rows, err := s.db.Query(ctx,
		`SELECT id, post_id, file_name, file_path, file_size, mime_type, created_at FROM post_attachments WHERE post_id = $1`, postID,
	)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var attachments []models.Attachment
	for rows.Next() {
		var a models.Attachment
		if err := rows.Scan(&a.ID, &a.PostID, &a.FileName, &a.FilePath, &a.FileSize, &a.MimeType, &a.CreatedAt); err == nil {
			attachments = append(attachments, a)
		}
	}
	if attachments == nil {
		attachments = []models.Attachment{}
	}
	return attachments
}
