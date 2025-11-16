package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"gopkg.in/natefinch/lumberjack.v2"

	"chat/auth/internal/app"
	"chat/auth/internal/interceptor"
	"chat/auth/internal/logger"
	"chat/auth/internal/metric"
	accessDesc "chat/auth/pkg/access_v1"
	authDesc "chat/auth/pkg/auth_v1"
	userDesc "chat/auth/pkg/user_v1"
	_ "chat/auth/statik"
)

var logLevel = flag.String("l", "info", "log level")

const grpcPort = 50051
const httpPort = 8081

func main() {
	flag.Parse()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	logger.Init(getCore(getAtomicLevel()))

	err = metric.Init(context.Background())
	if err != nil {
		log.Fatalf("failed to init metrics: %v", err)
	}

	serviceProvider := app.NewServiceProvider()

	dbClient := serviceProvider.GetDbClient(context.Background())
	defer func() {
		if err := dbClient.Close(); err != nil {
			log.Printf("failed to close database connection: %v", err)
		}
	}()

	userHandler := serviceProvider.GetUserHandler(context.Background())
	authHandler := serviceProvider.GetAuthHandler(context.Background())
	accessHandler := serviceProvider.GetAccessHandler(context.Background())

	grpcAddr := fmt.Sprintf(":%d", grpcPort)
	httpAddr := fmt.Sprintf(":%d", httpPort)
	swaggerAddr := serviceProvider.GetSwaggerConfig().Address()
	metricsHost := os.Getenv("METRICS_HOST")
	metricsPort := os.Getenv("METRICS_PORT")
	if metricsHost == "" {
		metricsHost = "localhost"
	}
	if metricsPort == "" {
		metricsPort = "2112"
	}
	metricsAddr := metricsHost + ":" + metricsPort

	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer wg.Done()
		if err := runGRPCServer(userHandler, authHandler, accessHandler, grpcAddr); err != nil {
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

	go func() {
		defer wg.Done()
		if err := runPrometheusServer(metricsAddr); err != nil {
			log.Fatalf("failed to run Prometheus server: %v", err)
		}
	}()

	wg.Wait()
}

func runGRPCServer(userHandler userDesc.UserV1Server, authHandler authDesc.AuthV1Server, accessHandler accessDesc.AccessV1Server, addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcSrv := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				interceptor.MetricsInterceptor,
				interceptor.LogInterceptor,
				interceptor.ValidateInterceptor,
			),
		),
	)
	reflection.Register(grpcSrv)
	userDesc.RegisterUserV1Server(grpcSrv, userHandler)
	authDesc.RegisterAuthV1Server(grpcSrv, authHandler)
	accessDesc.RegisterAccessV1Server(grpcSrv, accessHandler)
	log.Printf("Auth gRPC server listening on %v", lis.Addr())
	return grpcSrv.Serve(lis)
}

func runHTTPServer(grpcAddr, httpAddr string) error {
	mux := runtime.NewServeMux()
	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := userDesc.RegisterUserV1HandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	httpServer := &http.Server{
		Addr:              httpAddr,
		Handler:           corsMiddleware.Handler(mux),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Printf("Auth HTTP gateway listening on %s", httpAddr)
	return httpServer.ListenAndServe()
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
		defer func() {
			if err := file.Close(); err != nil {
				log.Printf("failed to close swagger file: %v", err)
			}
		}()

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

	swaggerServer := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Printf("Auth Swagger server listening on %s", addr)
	log.Printf("Visit http://%s to view the Swagger UI", addr)
	return swaggerServer.ListenAndServe()
}

func getCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/auth.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)
}

func getAtomicLevel() zap.AtomicLevel {
	var level zapcore.Level
	if err := level.Set(*logLevel); err != nil {
		log.Fatalf("failed to set log level: %v", err)
	}

	return zap.NewAtomicLevelAt(level)
}

func runPrometheusServer(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	prometheusServer := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Printf("Auth Prometheus server listening on %s", addr)
	log.Printf("Visit http://%s/metrics to view metrics", addr)

	return prometheusServer.ListenAndServe()
}
