package bots

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	delivery "github.com/goriiin/kotyari-bots_backend/internal/delivery/bots"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/bots"
	"github.com/goriiin/kotyari-bots_backend/internal/middleware/security"
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

	// Init DB
	pool, err := postgres.GetPool(ctx, b.config.Database)
	if err != nil {
		return fmt.Errorf("postgres.GetPool: %w", err)
	}
	defer pool.Close()

	// Init gRPC client for Profiles service
	// conn, err := grpc.NewClient(b.config.ProfilesSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	//	return fmt.Errorf("grpc.NewClient: %w", err)
	// }
	// defer conn.Close()
	// profilesClient := profiles.NewProfilesServiceClient(conn)

	// Init dependencies
	botsRepo := repo.NewBotsRepository(pool)
	botsUsecase := usecase.NewService(botsRepo)
	botsHandler := delivery.NewHandler(botsUsecase, nil)

	// Init HTTP server
	svr, err := gen.NewServer(botsHandler, &security.NotImplemented{})
	if err != nil {
		return fmt.Errorf("ogen.NewServer: %w", err)
	}

	httpAddr := fmt.Sprintf("%s:%d", b.config.API.Host, b.config.API.Port)
	b.server = &http.Server{
		Addr:    httpAddr,
		Handler: svr,
	}

	// Run server
	go func() {
		log.Printf("Bots service listening on %s", httpAddr)
		if err := b.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server exited with error: %v", err)
		}
	}()

	// Graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Shutting down bots service...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	return b.server.Shutdown(shutdownCtx)
}
