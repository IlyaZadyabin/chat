package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"chat/auth/internal/app"
	desc "chat/auth/pkg/user_v1"
)

const grpcPort = 50051

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	serviceProvider := app.NewServiceProvider()

	dbPool := serviceProvider.GetDbPool(context.Background())
	defer dbPool.Close()

	userHandler := serviceProvider.GetUserHandler(context.Background())

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}

	grpcSrv := grpc.NewServer()
	reflection.Register(grpcSrv)
	desc.RegisterUserV1Server(grpcSrv, userHandler)
	log.Printf("Auth server listening on %v", lis.Addr())
	if err := grpcSrv.Serve(lis); err != nil {
		log.Fatalf("serve error: %v", err)
	}
}
