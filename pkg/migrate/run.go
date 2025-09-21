package migrate

import (
	"flag"
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/goriiin/kotyari-bots_backend/pkg/postgres"
)

func Run(cfg postgres.Config) {
	dsn := postgres.ToUrl(cfg)
	if !strings.HasPrefix(dsn, "pgx5://") {
		dsn = "pgx5://" + strings.TrimPrefix(dsn, "postgres://")
	}

	log.Printf("Using DSN for migration: %s", dsn)

	m, err := migrate.New(
		"file:///app/migrations",
		dsn,
	)
	if err != nil {
		log.Fatalf("failed to create migrate instance: %v", err)
	}

	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("usage: migrate <up|down>")
	}

	command := args[0]
	switch command {
	case "up":
		log.Println("Applying migrations up...")
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("failed to apply migrations up: %v", err)
		}
		log.Println("Migrations applied successfully")
	case "down":
		log.Println("Applying migrations down...")
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("failed to apply migrations down: %v", err)
		}
		log.Println("Migrations rolled back successfully")
	default:
		log.Fatalf("unknown command: %s", command)
	}
}
