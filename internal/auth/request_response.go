package auth

import "time"

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type SignUpResponse struct {
	Status int `json:"status"`
}

type AuthResponse struct {
	UserUUID string `json:"user_uuid"`
}

type UserInfoResponse struct {
	Login     string    `json:"login"`
	CreatedAt time.Time `json:"created_at"`
}

type JwtRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type JwtResponse struct {
	Type         string `json:"type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshJwtRequest struct {
	RefreshToken string `json:"refresh_token"`
}
