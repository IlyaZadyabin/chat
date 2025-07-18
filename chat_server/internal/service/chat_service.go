package service

import (
	"context"
	"fmt"
	"time"

	"chat/chat_server/internal/repository"
	desc "chat/chat_server/pkg/chat_v1"
)

type ChatService struct {
	chatRepo *repository.ChatRepository
}

func NewChatService(chatRepo *repository.ChatRepository) *ChatService {
	return &ChatService{
		chatRepo: chatRepo,
	}
}

func (s *ChatService) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	if len(req.Usernames) == 0 {
		return nil, fmt.Errorf("at least one username is required")
	}

	if len(req.Usernames) > 10 {
		return nil, fmt.Errorf("maximum 10 users allowed per chat")
	}

	chatID, err := s.chatRepo.CreateChat(ctx, req.Usernames)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat: %w", err)
	}

	return &desc.CreateResponse{
		Id: chatID,
	}, nil
}

func (s *ChatService) Delete(ctx context.Context, req *desc.DeleteRequest) error {
	exists, err := s.chatRepo.ChatExists(ctx, req.Id)
	if err != nil {
		return fmt.Errorf("failed to check if chat exists: %w", err)
	}

	if !exists {
		return fmt.Errorf("chat not found")
	}

	err = s.chatRepo.DeleteChat(ctx, req.Id)
	if err != nil {
		return fmt.Errorf("failed to delete chat: %w", err)
	}

	return nil
}

func (s *ChatService) SendMessage(ctx context.Context, req *desc.SendMessageRequest) error {
	if req.Text == "" {
		return fmt.Errorf("message text cannot be empty")
	}

	if len(req.Text) > 1000 {
		return fmt.Errorf("message text too long (max 1000 characters)")
	}

	var timestamp time.Time
	if req.Timestamp != nil {
		timestamp = req.Timestamp.AsTime()
	} else {
		timestamp = time.Now()
	}

	err := s.chatRepo.SendMessage(ctx, req.From, req.Text, timestamp)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
