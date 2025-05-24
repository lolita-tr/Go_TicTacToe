package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"os"
	"strings"
	"time"
)

type JwtProvider struct {
	jwtKey        []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	signingMethod jwt.SigningMethod
}

func NewJwtProvider() JwtPr {
	jwtSecret := os.Getenv("JWT_SECRET")

	return &JwtProvider{
		jwtKey:        []byte(jwtSecret),
		accessExpiry:  15 * time.Minute,
		refreshExpiry: 60 * time.Minute,
		signingMethod: jwt.SigningMethodHS256,
	}
}

func (jp *JwtProvider) GenerateAccessToken(userUUID string) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userUUID,
		"iat":     time.Now().UTC().Unix(),
		"exp":     time.Now().Add(jp.accessExpiry).Unix(),
		"type":    "access",
	})

	accessString, err := accessToken.SignedString(jp.jwtKey)
	if err != nil {
		return "", errors.New("could not generate accessToken")
	}

	return accessString, nil
}

func (jp *JwtProvider) GenerateRefreshToken(userUUID string) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userUUID,
		"iat":     time.Now().UTC().Unix(),
		"exp":     time.Now().Add(60 * time.Minute).Unix(),
		"type":    "refresh",
	})

	refreshString, err := refreshToken.SignedString(jp.jwtKey)
	if err != nil {
		return "", errors.New("could not generate refreshToken")
	}

	return refreshString, nil
}

func (jp *JwtProvider) ValidateAccessToken(tokenString string) (string, error) {
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return "", errors.New("empty token")
	}

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return "", errors.New("token must have 3 parts separated by dots")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jp.jwtKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userUUID, ok := claims["user_id"].(string)
	if !ok || userUUID == "" {
		return "", errors.New("user_id not found in token")
	}

	return userUUID, nil
}

func (jp *JwtProvider) ValidateRefreshToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jp.jwtKey, nil
	})

	if err != nil {
		return false, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if tokenType, ok := claims["type"].(string); !ok || tokenType != "refresh" {
			return false, errors.New("invalid toke type")
		}
		return true, nil
	}

	return false, errors.New("invalid token")
}

func (jp *JwtProvider) ExtractUserUUIDFromToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jp.jwtKey, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	userUUID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("user_uuid not found in token")
	}

	return userUUID, nil
}
