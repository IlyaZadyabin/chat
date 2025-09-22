package user_v1

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"chat/auth/internal/converter"
	"chat/auth/internal/logger"
	"chat/auth/internal/service"
	desc "chat/auth/pkg/user_v1"
)

type UserV1Handler struct {
	desc.UnimplementedUserV1Server
	userService service.UserService
}

func NewUserV1Handler(userService service.UserService) *UserV1Handler {
	return &UserV1Handler{
		userService: userService,
	}
}

func (h *UserV1Handler) Register(server *grpc.Server) {
	desc.RegisterUserV1Server(server, h)
}

func (h *UserV1Handler) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	logger.Info("Creating user...", zap.String("name", req.GetInfo().GetName()), zap.String("email", req.GetInfo().GetEmail()))

	userCreate := converter.ToUserCreateFromDesc(req)

	id, err := h.userService.Create(ctx, userCreate)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &desc.CreateResponse{Id: id}, nil
}

func (h *UserV1Handler) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	logger.Info("Getting user...", zap.Int64("id", req.GetId()))

	user, err := h.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &desc.GetResponse{
		User: converter.ToUserFromService(user),
	}, nil
}

func (h *UserV1Handler) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	logger.Info("Updating user...", zap.Int64("id", req.GetId()))

	err := h.userService.Update(ctx, converter.ToUserUpdateFromDesc(req))
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (h *UserV1Handler) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	logger.Info("Deleting user...", zap.Int64("id", req.GetId()))

	err := h.userService.Delete(ctx, req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	return &emptypb.Empty{}, nil
}
