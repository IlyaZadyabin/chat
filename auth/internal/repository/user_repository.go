package repository

import (
	"context"

	"chat/auth/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, info *model.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.UserUpdate) error
	Delete(ctx context.Context, id int64) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByName(ctx context.Context, name string) (*model.User, error)
}
