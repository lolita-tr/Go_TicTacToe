package auth

import (
	"context"
	"net/http"
)

type Service interface {
	SignUp(request SignUpRequest) (int, error)
	Authorization(ctx context.Context, request JwtRequest) (*JwtResponse, error)
	RefreshAccessToken(refreshToken string) (*JwtResponse, error)
	RefreshRefreshToken(refreshToken string) (*JwtResponse, error)
	ExtractUserUUIDFromRequest(r *http.Request) (string, error)
}

type JwtPr interface {
	GenerateAccessToken(userUUID string) (string, error)
	GenerateRefreshToken(userUUID string) (string, error)
	ValidateAccessToken(tokenString string) (string, error)
	ValidateRefreshToken(tokenString string) (bool, error)
	ExtractUserUUIDFromToken(tokenString string) (string, error)
}
