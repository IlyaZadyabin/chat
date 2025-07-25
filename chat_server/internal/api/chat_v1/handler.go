package chat_v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"chat/chat_server/internal/converter"
	"chat/chat_server/internal/service"
	desc "chat/chat_server/pkg/chat_v1"
)

type ChatV1Handler struct {
	desc.UnimplementedChatV1Server
	chatService service.ChatService
}

func NewChatV1Handler(chatService service.ChatService) *ChatV1Handler {
	return &ChatV1Handler{
		chatService: chatService,
	}
}

func (h *ChatV1Handler) Register(server *grpc.Server) {
	desc.RegisterChatV1Server(server, h)
}

func (h *ChatV1Handler) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := h.chatService.Create(ctx, converter.ToChatCreateFromDesc(req))
	if err != nil {
		return nil, fmt.Errorf("failed to create chat: %w", err)
	}

	return &desc.CreateResponse{Id: id}, nil
}

func (h *ChatV1Handler) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := h.chatService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete chat: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (h *ChatV1Handler) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	err := h.chatService.SendMessage(ctx, converter.ToMessageFromDesc(req))
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return &emptypb.Empty{}, nil
}
