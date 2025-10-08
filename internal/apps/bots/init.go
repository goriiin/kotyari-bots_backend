package bots

import (
	"context"
	"fmt"
	profiles "github.com/goriiin/kotyari-bots_backend/api/protos/bot_profile/gen"
	adapter "github.com/goriiin/kotyari-bots_backend/internal/adapters/profiles"
	"github.com/goriiin/kotyari-bots_backend/pkg/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	delivery "github.com/goriiin/kotyari-bots_backend/internal/delivery_http/bots"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/bots"
	repo "github.com/goriiin/kotyari-bots_backend/internal/repo/bots"
	usecase "github.com/goriiin/kotyari-bots_backend/internal/usecase/bots"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
)

type BotsApp struct {
	config *BotsAppConfig
	server *http.Server
}

func NewApp(config *BotsAppConfig) *BotsApp {
	return &BotsApp{
		config: config,
	}
}

func (b *BotsApp) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Printf("info config: %v", b.config.Database)

	pool, err := postgres.GetPool(ctx, b.config.Database)
	if err != nil {
		return fmt.Errorf("postgres.GetPool: %w", err)
	}
	defer pool.Close()

	conn, err := grpc.NewClient(b.config.ProfilesSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("grpc.NewClient: %w", err)
	}
	defer conn.Close()
	profilesClient := profiles.NewProfilesServiceClient(conn)

	botsRepo := repo.NewBotsRepository(pool)
	profileValidator := adapter.NewGrpcValidator(profilesClient)
	botsUsecase := usecase.NewService(botsRepo, profileValidator)
	botsHandler := delivery.NewHandler(botsUsecase, profilesClient)

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
		log.Printf("Bots service listening on %s", httpAddr)
		if err := b.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server exited with error: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Shutting down bots service...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	return b.server.Shutdown(shutdownCtx)
}
