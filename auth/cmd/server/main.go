package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"chat/auth/internal/app"
	"chat/auth/internal/interceptor"
	desc "chat/auth/pkg/user_v1"
	_ "chat/auth/statik"
)

const grpcPort = 50051
const httpPort = 8081

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	serviceProvider := app.NewServiceProvider()

	dbClient := serviceProvider.GetDbClient(context.Background())
	defer dbClient.Close()

	userHandler := serviceProvider.GetUserHandler(context.Background())

	grpcAddr := fmt.Sprintf(":%d", grpcPort)
	httpAddr := fmt.Sprintf(":%d", httpPort)
	swaggerAddr := serviceProvider.GetSwaggerConfig().Address()

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		if err := runGRPCServer(userHandler, grpcAddr); err != nil {
			log.Fatalf("failed to run gRPC server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := runHTTPServer(grpcAddr, httpAddr); err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := runSwaggerServer(swaggerAddr); err != nil {
			log.Fatalf("failed to run Swagger server: %v", err)
		}
	}()

	wg.Wait()
}

func runGRPCServer(handler desc.UserV1Server, addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcSrv := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.ValidateInterceptor),
	)
	reflection.Register(grpcSrv)
	desc.RegisterUserV1Server(grpcSrv, handler)
	log.Printf("Auth gRPC server listening on %v", lis.Addr())
	return grpcSrv.Serve(lis)
}

func runHTTPServer(grpcAddr, httpAddr string) error {
	mux := runtime.NewServeMux()
	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := desc.RegisterUserV1HandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	log.Printf("Auth HTTP gateway listening on %s", httpAddr)
	return http.ListenAndServe(httpAddr, corsMiddleware.Handler(mux))
}

func runSwaggerServer(addr string) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))

	mux.HandleFunc("/api.swagger.json", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving swagger file: /api.swagger.json")

		file, err := statikFs.Open("/api.swagger.json")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Successfully served swagger file")
	})

	log.Printf("Auth Swagger server listening on %s", addr)
	log.Printf("Visit http://%s to view the Swagger UI", addr)
	return http.ListenAndServe(addr, mux)
}
