package service

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"chat/auth/internal/repository"
	desc "chat/auth/pkg/user_v1"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	id, err := s.userRepo.Create(ctx, req.Info, string(passwordHash))
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &desc.CreateResponse{Id: id}, nil
}

func (s *UserService) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	user, err := s.userRepo.Get(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &desc.GetResponse{
		User: user,
	}, nil
}

func (s *UserService) Update(ctx context.Context, req *desc.UpdateRequest) error {
	_, err := s.userRepo.Get(ctx, req.Id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if req.Info.Email != nil {
		existingUser, err := s.userRepo.GetByEmail(ctx, req.Info.Email.GetValue())
		if err == nil && existingUser.Id != req.Id {
			return fmt.Errorf("email is already taken by another user")
		}
	}

	err = s.userRepo.Update(ctx, req.Id, req.Info)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *UserService) Delete(ctx context.Context, req *desc.DeleteRequest) error {
	err := s.userRepo.Delete(ctx, req.Id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
