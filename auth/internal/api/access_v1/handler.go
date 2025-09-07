package access_v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"chat/auth/internal/service"
	desc "chat/auth/pkg/access_v1"
)

type AccessV1Handler struct {
	desc.UnimplementedAccessV1Server
	accessService service.AccessService
}

func NewAccessV1Handler(accessService service.AccessService) *AccessV1Handler {
	return &AccessV1Handler{
		accessService: accessService,
	}
}

func (h *AccessV1Handler) Register(server *grpc.Server) {
	desc.RegisterAccessV1Server(server, h)
}

func (h *AccessV1Handler) Check(ctx context.Context, req *desc.CheckRequest) (*emptypb.Empty, error) {
	err := h.accessService.Check(ctx, req.GetEndpointAddress())
	if err != nil {
		return nil, fmt.Errorf("access denied: %w", err)
	}

	return &emptypb.Empty{}, nil
}