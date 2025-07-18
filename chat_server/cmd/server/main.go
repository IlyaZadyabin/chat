package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"chat/chat_server/internal/database"
	"chat/chat_server/internal/repository"
	"chat/chat_server/internal/service"
	desc "chat/chat_server/pkg/chat_v1"
)

const grpcPort = 50052

type server struct {
	desc.UnimplementedChatV1Server
	chatService *service.ChatService
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("[CreateChat] usernames:%v", req.Usernames)
	return s.chatService.Create(ctx, req)
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("[DeleteChat] id:%d", req.Id)
	return &emptypb.Empty{}, s.chatService.Delete(ctx, req)
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf("[SendMessage] from:%s text:%s", req.From, req.Text)
	if req.Timestamp == nil {
		req.Timestamp = timestamppb.New(gofakeit.Date())
	}
	return &emptypb.Empty{}, s.chatService.SendMessage(ctx, req)
}

func main() {
	db, err := database.NewConnection(database.NewConfig())
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	chatRepo := repository.NewChatRepository(db)
	chatService := service.NewChatService(chatRepo)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}

	grpcSrv := grpc.NewServer()
	reflection.Register(grpcSrv)
	desc.RegisterChatV1Server(grpcSrv, &server{chatService: chatService})
	log.Printf("Chat server listening on %v", lis.Addr())
	if err := grpcSrv.Serve(lis); err != nil {
		log.Fatalf("serve error: %v", err)
	}
}
