package posts_query

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-faster/errors"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_query"
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

	// TODO: XDDD
	// httpAddr := fmt.Sprintf("%s:%d", p.appCfg.API.Host, p.appCfg.API.Port)
	httpServer := &http.Server{
		Addr:         "0.0.0.0:8089",
		Handler:      svr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("PostsApp HTTP service listening on %s", "0.0.0.0:8089")
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server exited with error: %w", err)
	}
	return nil
}
