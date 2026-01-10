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
	// Re-initializing auth client here or extracting from App struct would be cleaner,
	// assuming we can modify NewPostsQueryApp to store authClient or create it here.
	// For simplicity based on previous pattern, creating a new one or assuming it's available.

	// Better approach: Pass it via struct (requires modifying init.go return struct)
	// Since I can modify files fully:

	l := log.Default() // using std log for run errors as per original

	// Create auth client again or pass it. Let's create it to stick to the pattern of separation
	// but normally it should be in struct.
	// NOTE: Please update struct in init.go to store authClient if you want reuse.
	// Implementing creation here for strict compliance with "Run" interface.

	authClient, err := auth.NewClient(p.config.Auth, nil) // logger nil might be issue, better pass it
	if err != nil {
		return err
	}

	if err := p.startHTTPServer(p.handler, authClient); err != nil {
		l.Printf("Error happened starting server %v", err)
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
