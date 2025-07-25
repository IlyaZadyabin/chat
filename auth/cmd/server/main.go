package main

import (
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"chat/auth/internal/api/user_v1"
	"chat/auth/internal/database"
	"chat/auth/internal/repository"
	"chat/auth/internal/service"
	desc "chat/auth/pkg/user_v1"
)

const grpcPort = 50051

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

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := user_v1.NewUserV1Handler(*userService)

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
