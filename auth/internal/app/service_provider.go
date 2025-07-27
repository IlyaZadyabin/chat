package app

import (
	"context"
	"log"
	"sync"

	"chat/auth/internal/api/user_v1"
	"chat/auth/internal/database"
	"chat/auth/internal/repository"
	"chat/auth/internal/service"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ServiceProvider struct {
	dbOnce sync.Once
	dbPool *pgxpool.Pool

	userRepositoryOnce sync.Once
	userRepository     repository.UserRepository

	userServiceOnce sync.Once
	userService     *service.UserService

	userHandlerOnce sync.Once
	userHandler     *user_v1.UserV1Handler
}

func NewServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

func (s *ServiceProvider) GetDbPool(ctx context.Context) *pgxpool.Pool {
	s.dbOnce.Do(func() {
		var err error
		s.dbPool, err = database.NewConnection(database.NewConfig())
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
	})
	return s.dbPool
}

func (s *ServiceProvider) GetUserRepository(ctx context.Context) repository.UserRepository {
	s.userRepositoryOnce.Do(func() {
		s.userRepository = repository.NewUserRepository(s.GetDbPool(ctx))
	})
	return s.userRepository
}

func (s *ServiceProvider) GetUserService(ctx context.Context) *service.UserService {
	s.userServiceOnce.Do(func() {
		s.userService = service.NewUserService(s.GetUserRepository(ctx))
	})
	return s.userService
}

func (s *ServiceProvider) GetUserHandler(ctx context.Context) *user_v1.UserV1Handler {
	s.userHandlerOnce.Do(func() {
		s.userHandler = user_v1.NewUserV1Handler(s.GetUserService(ctx))
	})
	return s.userHandler
}
