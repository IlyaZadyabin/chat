package repository

import (
	"context"
	"time"
)

type ChatRepository interface {
	CreateChat(ctx context.Context, usernames []string) (int64, error)
	DeleteChat(ctx context.Context, chatID int64) error
	SendMessage(ctx context.Context, fromUser, text string, ts time.Time) error
	GetChatUsers(ctx context.Context, chatID int64) ([]string, error)
	ChatExists(ctx context.Context, chatID int64) (bool, error)
}
