package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/types/known/timestamppb"

	"chat/auth/internal/model"
	"chat/pkg/database/client"
)

type UserRepository interface {
	Create(ctx context.Context, info *model.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.UserUpdate) error
	Delete(ctx context.Context, id int64) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type userRepository struct {
	db client.Client
}

func NewUserRepository(db client.Client) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, info *model.UserInfo) (int64, error) {
	q := client.Query{
		Name: "user_repository.Create",
		QueryRaw: `
			INSERT INTO users (name, email, password_hash, role, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $5)
			RETURNING id`,
	}

	now := time.Now()

	var id int64
	err := r.db.DB().QueryRowContext(ctx, q,
		info.Name,
		info.Email,
		info.Password,
		info.Role,
		now,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return id, nil
}

func (r *userRepository) Get(ctx context.Context, id int64) (*model.User, error) {
	q := client.Query{
		Name: "user_repository.Get",
		QueryRaw: `
			SELECT id, name, email, password_hash, role, created_at, updated_at
			FROM users
			WHERE id = $1`,
	}

	row := r.db.DB().QueryRowContext(ctx, q, id)

	var user model.User
	var name, email, passwordHash, roleStr string
	var createdAt, updatedAt time.Time

	err := row.Scan(&user.ID, &name, &email, &passwordHash, &roleStr, &createdAt, &updatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.Info = &model.UserInfo{Name: name, Email: email, Role: roleStr}
	user.CreatedAt = timestamppb.New(createdAt)
	user.UpdatedAt = timestamppb.New(updatedAt)
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.UserUpdate) error {
	q := client.Query{
		Name: "user_repository.Update",
		QueryRaw: `
			UPDATE users
			SET name = COALESCE($2, name),
			    email = COALESCE($3, email),
			    role = COALESCE($4, role),
			    updated_at = $5
			WHERE id = $1`,
	}

	_, err := r.db.DB().ExecContext(ctx, q, user.ID, user.Info.Name, user.Info.Email, user.Info.Role, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	q := client.Query{
		Name:     "user_repository.Delete",
		QueryRaw: `DELETE FROM users WHERE id=$1`,
	}

	cmd, err := r.db.DB().ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	q := client.Query{
		Name: "user_repository.GetByEmail",
		QueryRaw: `
			SELECT id, name, email, password_hash, role, created_at, updated_at
			FROM users
			WHERE email = $1`,
	}

	row := r.db.DB().QueryRowContext(ctx, q, email)

	var user model.User
	var name, passwordHash, roleStr string
	var createdAt, updatedAt time.Time

	err := row.Scan(&user.ID, &name, &email, &passwordHash, &roleStr, &createdAt, &updatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	user.Info = &model.UserInfo{Name: name, Email: email, Role: roleStr}
	user.CreatedAt = timestamppb.New(createdAt)
	user.UpdatedAt = timestamppb.New(updatedAt)
	return &user, nil
}
