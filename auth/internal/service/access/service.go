package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	"google.golang.org/grpc/metadata"

	"chat/auth/internal/jwt"
	"chat/auth/internal/service"
)

const (
	authPrefix = "Bearer "
)

type accessService struct {
	tokenManager *jwt.TokenManager
}

func NewAccessService(tokenManager *jwt.TokenManager) service.AccessService {
	return &accessService{
		tokenManager: tokenManager,
	}
}

func (s *accessService) Check(ctx context.Context, endpointAddress string) error {
	if endpointAddress == "" {
		return fmt.Errorf("endpoint address is required")
	}

	log.Println("Checking access for endpoint:", endpointAddress)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return fmt.Errorf("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return fmt.Errorf("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	claims, err := s.tokenManager.VerifyAccessToken(accessToken)
	if err != nil {
		return fmt.Errorf("access token is invalid: %w", err)
	}

	accessibleRoles := s.getAccessibleRoles()

	requiredRole, ok := accessibleRoles[endpointAddress]
	if !ok {
		return nil
	}

	if requiredRole == claims.Role {
		return nil
	}

	return fmt.Errorf("access denied: role %s required, but user has role %s", requiredRole, claims.Role)
}

// getAccessibleRoles returns a map of endpoint addresses to required roles
func (s *accessService) getAccessibleRoles() map[string]string {
	return map[string]string{
		"/user_v1.UserV1/Create": "admin",
		"/user_v1.UserV1/Delete": "admin",
	}
}
