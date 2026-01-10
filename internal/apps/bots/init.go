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
	"github.com/goriiin/kotyari-bots_backend/internal/adapters/auth"
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
	"github.com/goriiin/kotyari-bots_backend/pkg/user"
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

// securityHandler реализует интерфейс ogen SecurityHandler
type securityHandler struct {
	authClient *auth.Client
}

func (s *securityHandler) HandleSessionAuth(ctx context.Context, operationName string, t gen.SessionAuth) (context.Context, error) {
	// t.APIKey содержит значение куки session_id
	userID, err := s.authClient.VerifySession(ctx, t.APIKey)
	if err != nil {
		return nil, err
	}
	return user.WithID(ctx, userID), nil
}

func (b *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	l := logger.NewLogger("bots-app", &b.config.ConfigBase)

	pool, err := postgres.GetPool(ctx, b.config.Database)
	if err != nil {
		return fmt.Errorf("postgres.GetPool: %w", err)
	}
	defer pool.Close()

	// Profiles Client
	conn, err := grpc.NewClient(b.config.ProfilesSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("grpc.NewClient: %w", err)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			log.Println("failed to close grpc connection")
		}
	}(conn)
	profilesClient := profiles.NewProfilesServiceClient(conn)

	// Auth Client
	authClient, err := auth.NewClient(b.config.Auth, l)
	if err != nil {
		return fmt.Errorf("auth.NewClient: %w", err)
	}

	botsRepo := repo.NewBotsRepository(pool)
	profileValidator := profiles_validator.NewGrpcValidator(profilesClient)
	profileGateway := profiles_getter.NewProfileGateway(profilesClient)
	botsUsecase := usecase.NewService(botsRepo, profileValidator, profileGateway)
	botsHandler := delivery.NewHandler(botsUsecase, l)

	// Инициализация Ogen Server с Security Handler
	secHandler := &securityHandler{authClient: authClient}
	svr, err := gen.NewServer(botsHandler, secHandler)
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
	botGrpcServer := bots.NewServer(botsUsecase, l)
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

	grpcServer.GracefulStop()
	if shutErr := b.server.Shutdown(shutdownCtx); shutErr != nil {
		return err
	}

	return nil
}
