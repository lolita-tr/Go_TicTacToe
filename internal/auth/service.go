package auth

import (
	"Go_TicTacToe/internal/game"
	"Go_TicTacToe/internal/storage"
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

type ServiceImpl struct {
	userRepository storage.UserRepository
	jwtProvider    JwtPr
}

const SUCCESS = 1

func NewAuthService(userRepo storage.UserRepository, jwtPr JwtPr) Service {
	return &ServiceImpl{
		userRepository: userRepo,
		jwtProvider:    jwtPr,
	}
}

func (as *ServiceImpl) SignUp(request SignUpRequest) (int, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, errors.New("failed to hash password")
	}

	userID := uuid.New().String()
	err = as.userRepository.CreateUser(context.Background(), userID, request.Login, string(hashPassword))

	if err != nil {
		return 0, errors.New("failed to register user")
	}

	return SUCCESS, nil
}

func (as *ServiceImpl) Authorization(ctx context.Context, request JwtRequest) (*JwtResponse, error) {
	userUUID, storeHash, err := as.userRepository.GetUserByLogin(ctx, request.Login)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if bcrypt.CompareHashAndPassword([]byte(storeHash), []byte(request.Password)) != nil {
		return nil, errors.New("invalid password")
	}

	accessToken, err := as.jwtProvider.GenerateAccessToken(userUUID)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := as.jwtProvider.GenerateRefreshToken(userUUID)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	response := &JwtResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Type:         "Bearer ",
	}

	return response, nil
}

func (as *ServiceImpl) RefreshAccessToken(refreshToken string) (*JwtResponse, error) {
	if ok, err := as.jwtProvider.ValidateRefreshToken(refreshToken); !ok {
		return nil, err
	}

	userUUID, err := as.jwtProvider.ExtractUserUUIDFromToken(refreshToken)
	if err != nil {
		return nil, err
	}

	accessToken, err := as.jwtProvider.GenerateAccessToken(userUUID)
	if err != nil {
		return nil, err
	}

	response := &JwtResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Type:         "Bearer ",
	}

	return response, nil
}

func (as *ServiceImpl) RefreshRefreshToken(refreshToken string) (*JwtResponse, error) {
	if ok, err := as.jwtProvider.ValidateRefreshToken(refreshToken); !ok {
		return nil, err
	}

	userUUID, err := as.jwtProvider.ExtractUserUUIDFromToken(refreshToken)
	if err != nil {
		return nil, err
	}

	accessToken, err := as.jwtProvider.GenerateAccessToken(userUUID)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := as.jwtProvider.GenerateRefreshToken(userUUID)
	if err != nil {
		return nil, err
	}

	response := &JwtResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		Type:         "Bearer ",
	}

	return response, nil
}

func (as *ServiceImpl) ExtractUserUUIDFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.New("missing or invalid token")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	return as.jwtProvider.ExtractUserUUIDFromToken(tokenStr)
}

func (as *ServiceImpl) GetUserInfoByToken(r *http.Request) (*UserInfoResponse, error) {
	userUUID, err := as.ExtractUserUUIDFromRequest(r)
	if err != nil {
		return nil, err
	}

	login, createdAt, err := as.userRepository.GetUserInfoByID(r.Context(), userUUID)
	if err != nil {
		return nil, err
	}

	response := &UserInfoResponse{
		login,
		createdAt,
	}

	return response, nil
}

func (as *ServiceImpl) GetFinishedGamesByToken(r *http.Request) (*game.FinishedGamesResponse, error) {
	userUUID, err := as.ExtractUserUUIDFromRequest(r)
	if err != nil {
		return nil, err
	}

	games, err := as.userRepository.GetFinishedGamesByToken(r.Context(), userUUID)
	if err != nil {
		return nil, errors.New("game not found")
	}

	response := &game.FinishedGamesResponse{
		games,
	}

	return response, nil
}
