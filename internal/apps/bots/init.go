package bots

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	profiles "github.com/goriiin/kotyari-bots_backend/api/protos/bot_profile/gen"
	botgrpc "github.com/goriiin/kotyari-bots_backend/api/protos/bots/gen"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_grpc/bots"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_grpc/profiles_getter"
	"github.com/goriiin/kotyari-bots_backend/internal/delivery_grpc/profiles_validator"
	delivery "github.com/goriiin/kotyari-bots_backend/internal/delivery_http/bots"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/bots"
	repo "github.com/goriiin/kotyari-bots_backend/internal/repo/bots"
	usecase "github.com/goriiin/kotyari-bots_backend/internal/usecase/bots"
	"github.com/goriiin/kotyari-bots_backend/pkg/cors"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	config *AppConfig
	server *http.Server
}

func NewApp(config *AppConfig) *App {
	return &App{
		config: config,
	}
}

func (b *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := postgres.GetPool(ctx, b.config.Database)
	if err != nil {
		return fmt.Errorf("postgres.GetPool: %w", err)
	}
	defer pool.Close()

	conn, err := grpc.NewClient(b.config.ProfilesSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("grpc.NewClient: %w", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Println("failed to close grpc connection")
		}
	}(conn)
	profilesClient := profiles.NewProfilesServiceClient(conn)

	botsRepo := repo.NewBotsRepository(pool)
	profileValidator := profiles_validator.NewGrpcValidator(profilesClient)
	profileGateway := profiles_getter.NewProfileGateway(profilesClient)
	botsUsecase := usecase.NewService(botsRepo, profileValidator, profileGateway)
	botsHandler := delivery.NewHandler(botsUsecase)

	svr, err := gen.NewServer(botsHandler)
	if err != nil {
		return fmt.Errorf("ogen.NewServer: %w", err)
	}

	httpAddr := fmt.Sprintf("%s:%d", b.config.API.Host, b.config.API.Port)
	b.server = &http.Server{
		Addr:    httpAddr,
		Handler: cors.New().Handler(svr),
	}

	go func() {
		log.Printf("Bots service listening on %s\n", httpAddr)
		if err := b.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("http server exited with error")
		}
	}()

	// gRPC Server Setup
	grpcAddr := fmt.Sprintf("%s:%d", b.config.GRPC.Host, b.config.GRPC.Port)
	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen for grpc: %w", err)
	}
	grpcServer := grpc.NewServer()
	botGrpcServer := bots.NewServer(botsUsecase)
	botgrpc.RegisterBotServiceServer(grpcServer, botGrpcServer)

	go func() {
		log.Printf("Bots gRPC service listening on %s\n", grpcAddr)
		if err = grpcServer.Serve(listener); err != nil {
			log.Println("gRPC server exited with error")
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Shutting down bots services...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	log.Println("Stopping gRPC server...")
	grpcServer.GracefulStop()
	log.Println("gRPC server stopped.")

	log.Println("Stopping HTTP server...")
	if err = b.server.Shutdown(shutdownCtx); err != nil {
		log.Println("HTTP server shutdown error")
		return err
	}
	log.Println("HTTP server stopped.")

	return nil
}
