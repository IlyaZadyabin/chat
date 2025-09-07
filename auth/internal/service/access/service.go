package service

import (
	"context"
	"fmt"

	"chat/auth/internal/service"
)

type accessService struct {
	// TODO: Add JWT token validation logic here when implementing actual JWT
}

func NewAccessService() service.AccessService {
	return &accessService{}
}

func (s *accessService) Check(ctx context.Context, endpointAddress string) error {
	// TODO: Implement actual access check logic with JWT token validation
	// For now, just do a basic validation that the endpoint is not empty
	if endpointAddress == "" {
		return fmt.Errorf("endpoint address is required")
	}

	// TODO: Extract JWT from request context and validate it
	// TODO: Check user permissions for the given endpoint
	
	// For this basic implementation, we'll allow access to all endpoints
	return nil
}