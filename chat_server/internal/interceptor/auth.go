package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	accesspb "chat/auth/pkg/access_v1"
)

type AuthInterceptor struct {
	accessClient accesspb.AccessV1Client
}

func NewAuthInterceptor(accessClient accesspb.AccessV1Client) *AuthInterceptor {
	return &AuthInterceptor{
		accessClient: accessClient,
	}
}

func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Printf("Auth interceptor: checking access for method %s", info.FullMethod)

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Printf("Auth interceptor: no metadata in context")
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		authHeader, ok := md["authorization"]
		if !ok || len(authHeader) == 0 {
			log.Printf("Auth interceptor: no authorization header")
			return nil, status.Errorf(codes.Unauthenticated, "missing authorization header")
		}

		outgoingCtx := metadata.NewOutgoingContext(ctx, md)

		checkReq := &accesspb.CheckRequest{
			EndpointAddress: info.FullMethod,
		}

		_, err := a.accessClient.Check(outgoingCtx, checkReq)
		if err != nil {
			log.Printf("Auth interceptor: access denied for method %s: %v", info.FullMethod, err)
			return nil, status.Errorf(codes.PermissionDenied, "access denied: %v", err)
		}

		log.Printf("Auth interceptor: access granted for method %s", info.FullMethod)

		return handler(ctx, req)
	}
}
