package app

import (
	"context"
	"log"
	"sync"

	"chat/chat_server/internal/api/chat_v1"
	"chat/chat_server/internal/database"
	"chat/chat_server/internal/repository"
	chatRepository "chat/chat_server/internal/repository/chat"
	"chat/chat_server/internal/service"
	chatService "chat/chat_server/internal/service/chat"
	"chat/pkg/database/client"
	"chat/pkg/database/pg"
	"chat/pkg/database/transaction"
)

type ServiceProvider struct {
	dbClientOnce sync.Once
	dbClient     client.Client

	txManagerOnce sync.Once
	txManager     client.TxManager

	chatRepositoryOnce sync.Once
	chatRepository     repository.ChatRepository

	chatServiceOnce sync.Once
	chatService     service.ChatService

	chatHandlerOnce sync.Once
	chatHandler     *chat_v1.ChatV1Handler
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

func (s *ServiceProvider) GetChatRepository(ctx context.Context) repository.ChatRepository {
	s.chatRepositoryOnce.Do(func() {
		s.chatRepository = chatRepository.NewChatRepository(s.GetDbClient(ctx))
	})
	return s.chatRepository
}

func (s *ServiceProvider) GetChatService(ctx context.Context) service.ChatService {
	s.chatServiceOnce.Do(func() {
		s.chatService = chatService.NewChatService(s.GetChatRepository(ctx))
	})
	return s.chatService
}

func (s *ServiceProvider) GetChatHandler(ctx context.Context) *chat_v1.ChatV1Handler {
	s.chatHandlerOnce.Do(func() {
		s.chatHandler = chat_v1.NewChatV1Handler(s.GetChatService(ctx))
	})
	return s.chatHandler
}
