package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "chat/auth/pkg/user_v1"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, info *desc.UserInfo, passwordHash string) (int64, error) {
	query := `
		INSERT INTO users (name, email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $5)
		RETURNING id`

	now := time.Now()

	var id int64
	err := r.pool.QueryRow(ctx, query,
		info.Name,
		info.Email,
		passwordHash,
		info.Role.String(),
		now,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return id, nil
}

func (r *UserRepository) Get(ctx context.Context, id int64) (*desc.User, error) {
	query := `
		SELECT id, name, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE id = $1`

	row := r.pool.QueryRow(ctx, query, id)

	var user desc.User
	var name, email, passwordHash, roleStr string
	var createdAt, updatedAt time.Time

	err := row.Scan(&user.Id, &name, &email, &passwordHash, &roleStr, &createdAt, &updatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	var role desc.Role
	switch roleStr {
	case "USER":
		role = desc.Role_USER
	case "ADMIN":
		role = desc.Role_ADMIN
	default:
		role = desc.Role_USER
	}

	user.Info = &desc.UserInfo{Name: name, Email: email, Role: role}
	user.CreatedAt = timestamppb.New(createdAt)
	user.UpdatedAt = timestamppb.New(updatedAt)
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, id int64, info *desc.UpdateUserInfo) error {
	query := `
		UPDATE users
		SET name = COALESCE($2, name),
		    email = COALESCE($3, email),
		    updated_at = $4
		WHERE id = $1`

	var name, email *string
	if info.Name != nil {
		v := info.Name.GetValue()
		name = &v
	}
	if info.Email != nil {
		v := info.Email.GetValue()
		email = &v
	}

	ct, err := r.pool.Exec(ctx, query, id, name, email, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*desc.User, error) {
	query := `
		SELECT id, name, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE email = $1`

	row := r.pool.QueryRow(ctx, query, email)

	var user desc.User
	var name, passwordHash, roleStr string
	var createdAt, updatedAt time.Time

	err := row.Scan(&user.Id, &name, &email, &passwordHash, &roleStr, &createdAt, &updatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	var role desc.Role
	switch roleStr {
	case "USER":
		role = desc.Role_USER
	case "ADMIN":
		role = desc.Role_ADMIN
	default:
		role = desc.Role_USER
	}

	user.Info = &desc.UserInfo{Name: name, Email: email, Role: role}
	user.CreatedAt = timestamppb.New(createdAt)
	user.UpdatedAt = timestamppb.New(updatedAt)
	return &user, nil
}
