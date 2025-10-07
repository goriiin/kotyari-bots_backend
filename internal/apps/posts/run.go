package posts

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-faster/errors"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts"
)

func (p *PostsApp) Run() error {
	if err := p.startHTTPServer(p.http); err != nil {
		log.Printf("Error happened starting server %v", err)
		return err
	}

	return nil
}

func (p *PostsApp) startHTTPServer(handler gen.Handler) error {
	svr, err := gen.NewServer(handler)
	if err != nil {
		return fmt.Errorf("ogen.NewServer: %w", err)
	}

	httpAddr := fmt.Sprintf("%s:%d", p.appCfg.API.Host, p.appCfg.API.Port)
	httpServer := &http.Server{
		Addr:    httpAddr,
		Handler: svr,
	}

	log.Printf("PostsApp HTTP service listening on %s", httpAddr)
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server exited with error: %w", err)
	}
	return nil
}
