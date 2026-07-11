// cmd/migrate is a standalone binary that applies the embedded SQL
// migrations (internal/db/postgres/migration) against Postgres. Built into
// its own small image (Dockerfile.migrate) and run as a Helm post-install/
// post-upgrade hook Job from shorturl-gitops - see
// helm/shorturl/templates/postgres-migrate-job.yaml there.
//
// Deliberately not run automatically by the main server binary (main.go):
// keeping it a separate one-shot process means migrations run exactly
// once per deploy (as a Job), not once per replica/pod restart, with no
// need for an advisory-lock dance.
package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/mykyta-kravchenko98/ShortUrl/internal/db/postgres/migration"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		mustEnv("POSTGRES_USER"),
		mustEnv("POSTGRES_PASSWORD"),
		mustEnv("POSTGRES_HOST"),
		envOr("POSTGRES_PORT", "5432"),
		mustEnv("POSTGRES_DB"),
		envOr("POSTGRES_SSLMODE", "disable"),
	)

	src, err := iofs.New(migration.Files, ".")
	if err != nil {
		slog.Error("failed to load embedded migrations", "error", err)
		os.Exit(1)
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, dsn)
	if err != nil {
		slog.Error("failed to initialize migrate", "error", err)
		os.Exit(1)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}

	slog.Info("migrations applied (or already up to date)")
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		slog.Error("missing required env var", "key", key)
		os.Exit(1)
	}
	return v
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
