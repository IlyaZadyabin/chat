package service

import (
	"context"
	"fmt"

	"chat/chat_server/internal/model"
	"chat/chat_server/internal/repository"
)

type ChatService struct {
	chatRepo repository.ChatRepository
}

func NewChatService(chatRepo repository.ChatRepository) *ChatService {
	return &ChatService{
		chatRepo: chatRepo,
	}
}

func (s *ChatService) Create(ctx context.Context, req *model.ChatCreate) (int64, error) {
	if len(req.Usernames) == 0 {
		return 0, fmt.Errorf("at least one username is required")
	}

	if len(req.Usernames) > 10 {
		return 0, fmt.Errorf("maximum 10 users allowed per chat")
	}

	chatID, err := s.chatRepo.CreateChat(ctx, req.Usernames)
	if err != nil {
		return 0, fmt.Errorf("failed to create chat: %w", err)
	}

	return chatID, nil
}

func (s *ChatService) Delete(ctx context.Context, id int64) error {
	exists, err := s.chatRepo.ChatExists(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check if chat exists: %w", err)
	}

	if !exists {
		return fmt.Errorf("chat not found")
	}

	err = s.chatRepo.DeleteChat(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete chat: %w", err)
	}

	return nil
}

func (s *ChatService) SendMessage(ctx context.Context, msg *model.Message) error {
	if msg.Text == "" {
		return fmt.Errorf("message text cannot be empty")
	}

	if len(msg.Text) > 1000 {
		return fmt.Errorf("message text too long (max 1000 characters)")
	}

	err := s.chatRepo.SendMessage(ctx, msg.From, msg.Text, msg.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
