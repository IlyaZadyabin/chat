package main

import (
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"chat/chat_server/internal/api/chat_v1"
	"chat/chat_server/internal/database"
	"chat/chat_server/internal/repository"
	"chat/chat_server/internal/service"
	desc "chat/chat_server/pkg/chat_v1"
)

const grpcPort = 50052

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db, err := database.NewConnection(database.NewConfig())
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	defer db.Close()

	chatRepo := repository.NewChatRepository(db)
	chatService := service.NewChatService(chatRepo)
	chatHandler := chat_v1.NewChatV1Handler(*chatService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}

	grpcSrv := grpc.NewServer()
	reflection.Register(grpcSrv)
	desc.RegisterChatV1Server(grpcSrv, chatHandler)
	log.Printf("Chat server listening on %v", lis.Addr())
	if err := grpcSrv.Serve(lis); err != nil {
		log.Fatalf("serve error: %v", err)
	}
}
