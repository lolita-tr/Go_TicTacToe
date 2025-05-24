package middleware

import (
	"Go_TicTacToe/internal/auth"
	"context"
	"net/http"
	"strings"
)

type UserAuthenticator struct {
	authService auth.Service
	jwt         auth.JwtPr
}

func NewUserAuthenticator(authService auth.Service, jwt auth.JwtPr) *UserAuthenticator {
	return &UserAuthenticator{
		authService: authService,
		jwt:         jwt,
	}
}

func (u *UserAuthenticator) AuthMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/signup" || r.URL.Path == "/login" {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := u.jwt.ValidateAccessToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
