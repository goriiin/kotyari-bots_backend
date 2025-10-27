package posts_command_producer

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-faster/errors"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_command"
)

func (p *PostsCommandProducerApp) Run() error {
	if err := p.startHTTPServer(p.handler); err != nil {
		log.Printf("Error happened starting server %v", err)
		return err
	}
	return nil
}

func (p *PostsCommandProducerApp) startHTTPServer(handler gen.Handler) error {
	svr, err := gen.NewServer(handler)
	if err != nil {
		return fmt.Errorf("ogen.NewServer: %w", err)
	}

	// TODO: XDDD
	// httpAddr := fmt.Sprintf("%s:%d", p.appCfg.API.Host, p.appCfg.API.Port)
	httpServer := &http.Server{
		Addr:         "0.0.0.0:8088",
		Handler:      svr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("PostsApp HTTP service listening on %s", "0.0.0.0:8088")
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server exited with error: %w", err)
	}
	return nil
}
