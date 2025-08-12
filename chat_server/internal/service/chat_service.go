package service

import (
	"context"

	"chat/chat_server/internal/model"
)

type ChatService interface {
	Create(ctx context.Context, req *model.ChatCreate) (int64, error)
	Delete(ctx context.Context, id int64) error
	SendMessage(ctx context.Context, msg *model.Message) error
}
