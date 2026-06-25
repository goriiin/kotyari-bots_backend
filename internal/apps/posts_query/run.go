package posts_query

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-faster/errors"
	"github.com/goriiin/kotyari-bots_backend/internal/adapters/auth"
	gen "github.com/goriiin/kotyari-bots_backend/internal/gen/posts/posts_query"
	"github.com/goriiin/kotyari-bots_backend/pkg/cors"
)

func (p *PostsQueryApp) Run() error {
	// Reuse the auth client built (with a real logger) in NewPostsQueryApp instead
	// of constructing a second one with a nil logger, which panicked on the first
	// auth failure inside VerifySession.
	if err := p.startHTTPServer(p.handler, p.authClient); err != nil {
		log.Printf("Error happened starting server %v", err)
		return err
	}
	return nil
}

func (p *PostsQueryApp) startHTTPServer(handler gen.Handler, authClient *auth.Client) error {
	secHandler := &securityHandler{authClient: authClient}
	svr, err := gen.NewServer(handler, secHandler)
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
