package app

import (
	"context"
	"log"
	"sync"

	"chat/auth/internal/api/user_v1"
	"chat/auth/internal/database"
	"chat/auth/internal/repository"
	userRepository "chat/auth/internal/repository/user"
	"chat/auth/internal/service"
	userService "chat/auth/internal/service/user"
	"common/database/client"
	"common/database/pg"
	"common/database/transaction"
)

type ServiceProvider struct {
	dbClientOnce sync.Once
	dbClient     client.Client

	txManagerOnce sync.Once
	txManager     client.TxManager

	userRepositoryOnce sync.Once
	userRepository     repository.UserRepository

	userServiceOnce sync.Once
	userService     service.UserService

	userHandlerOnce sync.Once
	userHandler     *user_v1.UserV1Handler
}

func NewServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

func (s *ServiceProvider) GetDbClient(ctx context.Context) client.Client {
	s.dbClientOnce.Do(func() {
		cfg := database.NewConfig()
		dbClient, err := pg.New(ctx, cfg.GetDSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}
		s.dbClient = dbClient
	})
	return s.dbClient
}

func (s *ServiceProvider) GetTxManager(ctx context.Context) client.TxManager {
	s.txManagerOnce.Do(func() {
		s.txManager = transaction.NewTransactionManager(s.GetDbClient(ctx).DB())
	})
	return s.txManager
}

func (s *ServiceProvider) GetUserRepository(ctx context.Context) repository.UserRepository {
	s.userRepositoryOnce.Do(func() {
		s.userRepository = userRepository.NewUserRepository(s.GetDbClient(ctx))
	})
	return s.userRepository
}

func (s *ServiceProvider) GetUserService(ctx context.Context) service.UserService {
	s.userServiceOnce.Do(func() {
		s.userService = userService.NewUserService(s.GetUserRepository(ctx))
	})
	return s.userService
}

func (s *ServiceProvider) GetUserHandler(ctx context.Context) *user_v1.UserV1Handler {
	s.userHandlerOnce.Do(func() {
		s.userHandler = user_v1.NewUserV1Handler(s.GetUserService(ctx))
	})
	return s.userHandler
}
