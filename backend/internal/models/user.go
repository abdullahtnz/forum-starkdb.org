package models

import "time"

type User struct {
	ID            string     `json:"id"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	PasswordHash  string     `json:"-"`
	Bio           string     `json:"bio"`
	AvatarURL     string     `json:"avatar_url"`
	IsAdmin       bool       `json:"is_admin"`
	AcceptedTerms bool       `json:"-"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	LastLogin     *time.Time `json:"last_login"`
}

type UserPublic struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Bio       string    `json:"bio"`
	AvatarURL string    `json:"avatar_url"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

type UserProfile struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Bio       string     `json:"bio"`
	AvatarURL string     `json:"avatar_url"`
	IsAdmin   bool       `json:"is_admin"`
	CreatedAt time.Time  `json:"created_at"`
	LastLogin *time.Time `json:"last_login"`
}
