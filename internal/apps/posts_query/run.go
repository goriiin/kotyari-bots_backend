package posts_query

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-faster/errors"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_query"
	"github.com/goriiin/kotyari-bots_backend/pkg/cors"
)

func (p *PostsQueryApp) Run() error {
	if err := p.startHTTPServer(p.handler); err != nil {
		log.Printf("Error happened starting server %v", err)
		return err
	}
	return nil
}

func (p *PostsQueryApp) startHTTPServer(handler gen.Handler) error {
	svr, err := gen.NewServer(handler)
	if err != nil {
		return fmt.Errorf("ogen.NewServer: %w", err)
	}

	httpAddr := fmt.Sprintf("%s:%d", p.config.API.Host, p.config.API.Port)
	httpServer := &http.Server{
		Addr:         httpAddr,
		Handler:      cors.New().Handler(svr),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("PostsQueryApp HTTP service listening on %s", httpAddr)
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server exited with error: %w", err)
	}
	return nil
}
