package service

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"chat/auth/internal/jwt"
	"chat/auth/internal/repository"
	"chat/auth/internal/service"
)

type authService struct {
	userRepo     repository.UserRepository
	tokenManager *jwt.TokenManager
}

func NewAuthService(userRepo repository.UserRepository, tokenManager *jwt.TokenManager) service.AuthService {
	return &authService{
		userRepo:     userRepo,
		tokenManager: tokenManager,
	}
}

func (s *authService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepo.GetByName(ctx, username)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Info.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials: %w", err)
	}

	userInfo := jwt.UserInfo{
		UserID:   user.ID,
		Username: user.Info.Name,
		Role:     user.Info.Role,
	}

	refreshToken, err := s.tokenManager.GenerateRefreshToken(userInfo)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return refreshToken, nil
}

func (s *authService) GetRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := s.tokenManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	user, err := s.userRepo.GetByName(ctx, claims.Username)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	userInfo := jwt.UserInfo{
		UserID:   user.ID,
		Username: user.Info.Name,
		Role:     user.Info.Role,
	}

	newRefreshToken, err := s.tokenManager.GenerateRefreshToken(userInfo)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return newRefreshToken, nil
}

func (s *authService) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := s.tokenManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	user, err := s.userRepo.GetByName(ctx, claims.Username)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	userInfo := jwt.UserInfo{
		UserID:   user.ID,
		Username: user.Info.Name,
		Role:     user.Info.Role,
	}

	accessToken, err := s.tokenManager.GenerateAccessToken(userInfo)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, nil
}
