package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/abdullahtnz/forum-starkdb/internal/models"
	"github.com/abdullahtnz/forum-starkdb/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthService struct {
	db         *pgxpool.Pool
	jwtSecret  string
	jwtRefresh string
}

func NewAuthService(db *pgxpool.Pool, jwtSecret, jwtRefresh string) *AuthService {
	return &AuthService{db: db, jwtSecret: jwtSecret, jwtRefresh: jwtRefresh}
}

func (s *AuthService) Signup(ctx context.Context, req models.SignupRequest) (*models.TokenResponse, error) {
	if !req.AcceptedTerms {
		return nil, errors.New("you must accept the terms and privacy policy")
	}
	username := utils.SanitizeUsername(req.Username)
	if len(username) < 3 || len(username) > 30 {
		return nil, errors.New("username must be between 3 and 30 characters")
	}
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	var userID string
	err = s.db.QueryRow(ctx,
		`INSERT INTO users (username, email, password_hash, accepted_terms) VALUES ($1, $2, $3, $4) RETURNING id`,
		username, req.Email, hash, true,
	).Scan(&userID)
	if err != nil {
		if isDuplicate(err) {
			return nil, errors.New("username or email already exists")
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return s.generateTokens(userID, username, false)
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.TokenResponse, error) {
	var user models.User
	err := s.db.QueryRow(ctx,
		`SELECT id, username, email, password_hash, is_admin FROM users WHERE email = $1`,
		req.Email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.IsAdmin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("invalid email or password")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}
	_, _ = s.db.Exec(ctx, `UPDATE users SET last_login = NOW() WHERE id = $1`, user.ID)
	return s.generateTokens(user.ID, user.Username, user.IsAdmin)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*models.TokenResponse, error) {
	tokenHash := utils.HashToken(refreshToken)
	var userID, username string
	var isAdmin bool
	err := s.db.QueryRow(ctx,
		`DELETE FROM refresh_tokens WHERE token_hash = $1 AND expires_at > NOW() RETURNING user_id`,
		tokenHash,
	).Scan(&userID)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}
	err = s.db.QueryRow(ctx,
		`SELECT username, is_admin FROM users WHERE id = $1`, userID,
	).Scan(&username, &isAdmin)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return s.generateTokens(userID, username, isAdmin)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := utils.HashToken(refreshToken)
	_, err := s.db.Exec(ctx, `DELETE FROM refresh_tokens WHERE token_hash = $1`, tokenHash)
	return err
}

func (s *AuthService) generateTokens(userID, username string, isAdmin bool) (*models.TokenResponse, error) {
	accessToken, err := utils.GenerateAccessToken(userID, username, s.jwtSecret, isAdmin)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	tokenHash := utils.HashToken(refreshToken)
	_, err = s.db.Exec(context.Background(),
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
		userID, tokenHash, time.Now().Add(7*24*time.Hour),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}
	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900,
	}, nil
}

func isDuplicate(err error) bool {
	return err != nil && (errors.Is(err, pgx.ErrNoRows) == false && contains(err.Error(), "unique") || contains(err.Error(), "duplicate"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
