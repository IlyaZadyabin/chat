package service

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"chat/auth/internal/repository"
	"chat/auth/internal/service"
)

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) service.AuthService {
	return &authService{
		userRepo: userRepo,
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

	// TODO: Generate actual JWT refresh token
	refreshToken := fmt.Sprintf("refresh_token_for_user_%d", user.ID)
	
	return refreshToken, nil
}

func (s *authService) GetRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// TODO: Validate old refresh token and extract user info
	// For now, just return a new mock token
	newRefreshToken := fmt.Sprintf("new_%s", refreshToken)
	
	return newRefreshToken, nil
}

func (s *authService) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	// TODO: Validate refresh token and extract user info
	// For now, just return a mock access token
	accessToken := fmt.Sprintf("access_token_from_%s", refreshToken)
	
	return accessToken, nil
}