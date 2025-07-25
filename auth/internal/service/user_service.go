package service

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"chat/auth/internal/model"
	"chat/auth/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Create(ctx context.Context, userCreate *model.UserCreate) (int64, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userCreate.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	userCreate.Info.Password = string(passwordHash)

	id, err := s.userRepo.Create(ctx, userCreate.Info)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (s *UserService) Get(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *UserService) Update(ctx context.Context, userUpdate *model.UserUpdate) error {
	_, err := s.userRepo.Get(ctx, userUpdate.ID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if userUpdate.Info.Email != "" {
		existingUser, err := s.userRepo.GetByEmail(ctx, userUpdate.Info.Email)
		if err == nil && existingUser.ID != userUpdate.ID {
			return fmt.Errorf("email is already taken by another user")
		}
	}

	err = s.userRepo.Update(ctx, userUpdate)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
