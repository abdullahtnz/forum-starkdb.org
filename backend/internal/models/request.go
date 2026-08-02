package models

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Username      string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	AcceptedTerms bool   `json:"accepted_terms"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type UpdateProfileRequest struct {
	Username string `json:"username"`
	Bio      string `json:"bio"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

type PaginationQuery struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Sort     string `json:"sort"`
	Tag      string `json:"tag,omitempty"`
	Search   string `json:"q,omitempty"`
}

type AdminActionRequest struct {
	Action   string `json:"action"`
	Reason   string `json:"reason,omitempty"`
	TargetID string `json:"target_id,omitempty"`
}
