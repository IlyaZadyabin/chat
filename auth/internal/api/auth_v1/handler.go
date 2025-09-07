package auth_v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"chat/auth/internal/service"
	desc "chat/auth/pkg/auth_v1"
)

type AuthV1Handler struct {
	desc.UnimplementedAuthV1Server
	authService service.AuthService
}

func NewAuthV1Handler(authService service.AuthService) *AuthV1Handler {
	return &AuthV1Handler{
		authService: authService,
	}
}

func (h *AuthV1Handler) Register(server *grpc.Server) {
	desc.RegisterAuthV1Server(server, h)
}

func (h *AuthV1Handler) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	refreshToken, err := h.authService.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	return &desc.LoginResponse{RefreshToken: refreshToken}, nil
}

func (h *AuthV1Handler) GetRefreshToken(ctx context.Context, req *desc.GetRefreshTokenRequest) (*desc.GetRefreshTokenResponse, error) {
	newRefreshToken, err := h.authService.GetRefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	return &desc.GetRefreshTokenResponse{RefreshToken: newRefreshToken}, nil
}

func (h *AuthV1Handler) GetAccessToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	accessToken, err := h.authService.GetAccessToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	return &desc.GetAccessTokenResponse{AccessToken: accessToken}, nil
}