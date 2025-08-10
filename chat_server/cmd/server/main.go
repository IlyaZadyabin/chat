package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"chat/chat_server/internal/app"
	desc "chat/chat_server/pkg/chat_v1"
)

const grpcPort = 50052

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	serviceProvider := app.NewServiceProvider()

	dbClient := serviceProvider.GetDbClient(context.Background())
	defer dbClient.Close()

	chatHandler := serviceProvider.GetChatHandler(context.Background())

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
