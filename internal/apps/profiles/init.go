package profiles

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
	deliverygrpc "github.com/goriiin/kotyari-bots_backend/internal/delivery_grpc/profiles"
	deliveryhttp "github.com/goriiin/kotyari-bots_backend/internal/delivery_http/profiles"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/profiles"
	"github.com/goriiin/kotyari-bots_backend/internal/logger"
	repo "github.com/goriiin/kotyari-bots_backend/internal/repo/profiles"
	usecase "github.com/goriiin/kotyari-bots_backend/internal/usecase/profiles"
	"github.com/goriiin/kotyari-bots_backend/pkg/cors"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

const serviceName = "profiles-service"

type ProfilesApp struct {
	config     *ProfilesAppConfig
	httpServer *http.Server
	grpcServer *grpc.Server
	log        *logger.Logger
}

func NewApp(config *ProfilesAppConfig) *ProfilesApp {
	return &ProfilesApp{
		config: config,
		log:    logger.NewLogger(serviceName, &config.ConfigBase),
	}
}

func (p *ProfilesApp) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := postgres.GetPool(ctx, p.config.Database)
	if err != nil {
		return fmt.Errorf("postgres.GetPool: %w", err)
	}
	defer pool.Close()

	profilesRepo := repo.NewRepository(pool)
	profilesUsecase := usecase.NewService(profilesRepo)

	grpcHandler := deliverygrpc.NewGRPCHandler(profilesUsecase, p.log)
	httpHandler := deliveryhttp.NewHTTPHandler(profilesUsecase, p.log)

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return p.startGRPCServer(gCtx, grpcHandler)
	})

	g.Go(func() error {
		return p.startHTTPServer(gCtx, httpHandler)
	})

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sig:
		p.log.Info("Received shutdown signal")
	case <-gCtx.Done():
		p.log.Info("Context cancelled, shutting down...")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if p.httpServer != nil {
		if err := p.httpServer.Shutdown(shutdownCtx); err != nil {
			p.log.Error(err, false, "HTTP server shutdown error")
		}
	}
	if p.grpcServer != nil {
		p.grpcServer.GracefulStop()
	}

	p.log.Info("Profiles service stopped gracefully.")
	return g.Wait()
}

func (p *ProfilesApp) startGRPCServer(ctx context.Context, handler profiles.ProfilesServiceServer) error {
	grpcAddr := fmt.Sprintf(":%d", p.config.API.GRPCPort)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen gRPC on %s: %w", grpcAddr, err)
	}

	p.grpcServer = grpc.NewServer()
	profiles.RegisterProfilesServiceServer(p.grpcServer, handler)

	p.log.Info(fmt.Sprintf("Profiles gRPC service listening on %s", grpcAddr))
	return p.grpcServer.Serve(lis)
}

func (p *ProfilesApp) startHTTPServer(ctx context.Context, handler gen.Handler) error {
	svr, err := gen.NewServer(handler)
	if err != nil {
		return fmt.Errorf("ogen.NewServer: %w", err)
	}

	httpAddr := fmt.Sprintf("%s:%d", p.config.API.Host, p.config.API.Port)
	p.httpServer = &http.Server{
		Addr:    httpAddr,
		Handler: cors.New().Handler(svr),
	}

	p.log.Info(fmt.Sprintf("Profiles HTTP service listening on %s", httpAddr))
	if err := p.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("http server exited with error: %w", err)
	}
	return nil
}
