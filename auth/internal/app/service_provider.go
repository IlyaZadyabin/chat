package app

import (
	"context"
	"log"
	"sync"

	"chat/auth/internal/api/access_v1"
	"chat/auth/internal/api/auth_v1"
	"chat/auth/internal/api/user_v1"
	"chat/auth/internal/config"
	"chat/auth/internal/database"
	"chat/auth/internal/jwt"
	"chat/auth/internal/repository"
	userRepository "chat/auth/internal/repository/user"
	"chat/auth/internal/service"
	accessService "chat/auth/internal/service/access"
	authService "chat/auth/internal/service/auth"
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

	authServiceOnce sync.Once
	authService     service.AuthService

	accessServiceOnce sync.Once
	accessService     service.AccessService

	userHandlerOnce sync.Once
	userHandler     *user_v1.UserV1Handler

	authHandlerOnce sync.Once
	authHandler     *auth_v1.AuthV1Handler

	accessHandlerOnce sync.Once
	accessHandler     *access_v1.AccessV1Handler

	jwtTokenManagerOnce sync.Once
	jwtTokenManager     *jwt.TokenManager

	swaggerConfigOnce sync.Once
	swaggerConfig     config.SwaggerConfig
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

func (s *ServiceProvider) GetAuthService(ctx context.Context) service.AuthService {
	s.authServiceOnce.Do(func() {
		s.authService = authService.NewAuthService(
			s.GetUserRepository(ctx),
			s.GetJWTTokenManager(ctx),
		)
	})
	return s.authService
}

func (s *ServiceProvider) GetAccessService(ctx context.Context) service.AccessService {
	s.accessServiceOnce.Do(func() {
		s.accessService = accessService.NewAccessService(s.GetJWTTokenManager(ctx))
	})
	return s.accessService
}

func (s *ServiceProvider) GetAuthHandler(ctx context.Context) *auth_v1.AuthV1Handler {
	s.authHandlerOnce.Do(func() {
		s.authHandler = auth_v1.NewAuthV1Handler(s.GetAuthService(ctx))
	})
	return s.authHandler
}

func (s *ServiceProvider) GetAccessHandler(ctx context.Context) *access_v1.AccessV1Handler {
	s.accessHandlerOnce.Do(func() {
		s.accessHandler = access_v1.NewAccessV1Handler(s.GetAccessService(ctx))
	})
	return s.accessHandler
}

func (s *ServiceProvider) GetJWTTokenManager(ctx context.Context) *jwt.TokenManager {
	s.jwtTokenManagerOnce.Do(func() {
		jwtConfig := config.NewJWTConfig()
		s.jwtTokenManager = jwt.NewTokenManager(
			jwtConfig.RefreshSecretKey,
			jwtConfig.AccessSecretKey,
			jwtConfig.RefreshTokenExpiry,
			jwtConfig.AccessTokenExpiry,
		)
	})
	return s.jwtTokenManager
}

func (s *ServiceProvider) GetSwaggerConfig() config.SwaggerConfig {
	s.swaggerConfigOnce.Do(func() {
		s.swaggerConfig = config.NewSwaggerConfig()
	})
	return s.swaggerConfig
}
