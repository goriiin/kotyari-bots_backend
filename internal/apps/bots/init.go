package bots

import (
	"context"
	"fmt"
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
	"github.com/goriiin/kotyari-bots_backend/internal/logger"
	repo "github.com/goriiin/kotyari-bots_backend/internal/repo/bots"
	usecase "github.com/goriiin/kotyari-bots_backend/internal/usecase/bots"
	"github.com/goriiin/kotyari-bots_backend/pkg/cors"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const serviceName = "bots-service"

type App struct {
	config *AppConfig
	server *http.Server
	log    *logger.Logger
}

func NewApp(config *AppConfig) *App {
	return &App{
		config: config,
		log:    logger.NewLogger(serviceName, &config.ConfigBase),
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
			b.log.Error(err, false, "failed to close grpc connection")
		}
	}(conn)
	profilesClient := profiles.NewProfilesServiceClient(conn)

	botsRepo := repo.NewBotsRepository(pool)
	profileValidator := profiles_validator.NewGrpcValidator(profilesClient)
	profileGateway := profiles_getter.NewProfileGateway(profilesClient)
	botsUsecase := usecase.NewService(botsRepo, profileValidator, profileGateway)

	// Внедряем логгер
	botsHandler := delivery.NewHandler(botsUsecase, b.log)

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
		b.log.Info(fmt.Sprintf("Bots service listening on %s", httpAddr))
		if err := b.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			b.log.Fatal(err, true, "http server exited with error")
		}
	}()

	// gRPC Server Setup
	grpcAddr := fmt.Sprintf("%s:%d", b.config.GRPC.Host, b.config.GRPC.Port)
	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen for grpc: %w", err)
	}
	grpcServer := grpc.NewServer()

	// Внедряем логгер
	botGrpcServer := bots.NewServer(botsUsecase, b.log)

	botgrpc.RegisterBotServiceServer(grpcServer, botGrpcServer)

	go func() {
		b.log.Info(fmt.Sprintf("Bots gRPC service listening on %s", grpcAddr))
		if err = grpcServer.Serve(listener); err != nil {
			b.log.Fatal(err, true, "gRPC server exited with error")
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	b.log.Info("Shutting down bots services...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	b.log.Info("Stopping gRPC server...")
	grpcServer.GracefulStop()
	b.log.Info("gRPC server stopped.")

	b.log.Info("Stopping HTTP server...")
	if err = b.server.Shutdown(shutdownCtx); err != nil {
		b.log.Error(err, false, "HTTP server shutdown error")
		return err
	}
	b.log.Info("HTTP server stopped.")

	return nil
}
