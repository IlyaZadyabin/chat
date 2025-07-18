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

	"chat/auth/internal/database"
	"chat/auth/internal/repository"
	"chat/auth/internal/service"
	desc "chat/auth/pkg/user_v1"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedUserV1Server
	userService *service.UserService
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("[Create] Name:%s Email:%s Role:%s", req.Info.Name, req.Info.Email, req.Info.Role.String())
	return s.userService.Create(ctx, req)
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("[Get] ID:%d", req.Id)
	resp, err := s.userService.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	// if not found, generate fake to keep demo working
	if resp.User == nil {
		resp = &desc.GetResponse{User: &desc.User{Id: req.Id, Info: &desc.UserInfo{Name: gofakeit.Name(), Email: gofakeit.Email(), Role: desc.Role_USER}, CreatedAt: timestamppb.Now(), UpdatedAt: timestamppb.Now()}}
	}
	return resp, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("[Update] ID:%d", req.Id)
	return &emptypb.Empty{}, s.userService.Update(ctx, req)
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("[Delete] ID:%d", req.Id)
	return &emptypb.Empty{}, s.userService.Delete(ctx, req)
}

func main() {
	db, err := database.NewConnection(database.NewConfig())
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}

	grpcSrv := grpc.NewServer()
	reflection.Register(grpcSrv)
	desc.RegisterUserV1Server(grpcSrv, &server{userService: userService})
	log.Printf("Auth server listening on %v", lis.Addr())
	if err := grpcSrv.Serve(lis); err != nil {
		log.Fatalf("serve error: %v", err)
	}
}
