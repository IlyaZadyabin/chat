package service

import (
	"context"

	"chat/auth/internal/model"
)

type UserService interface {
	Create(ctx context.Context, userCreate *model.UserCreate) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, userUpdate *model.UserUpdate) error
	Delete(ctx context.Context, id int64) error
}
