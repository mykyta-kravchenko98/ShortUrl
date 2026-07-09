package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/mykyta-kravchenko98/ShortUrl/internal/cache"
	"github.com/mykyta-kravchenko98/ShortUrl/internal/config"
	repositories "github.com/mykyta-kravchenko98/ShortUrl/internal/db/postgres"
	"github.com/mykyta-kravchenko98/ShortUrl/internal/handler"
	"github.com/mykyta-kravchenko98/ShortUrl/internal/observability"
	"github.com/mykyta-kravchenko98/ShortUrl/internal/router"
	"github.com/mykyta-kravchenko98/ShortUrl/internal/service"
	"github.com/mykyta-kravchenko98/ShortUrl/pkg/closeutil"
	"github.com/mykyta-kravchenko98/ShortUrl/pkg/generator"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	otelShutdown, err := observability.Setup(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to set up observability:", err)
		os.Exit(1)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := otelShutdown(shutdownCtx); err != nil {
			slog.Error("otel shutdown failed", "error", err)
		}
	}()

	env := os.Getenv("environment")
	if env == "" {
		env = "dev"
	}

	var conf *config.Config
	var confErr error
	switch env {
	case "dev":
		conf, confErr = config.LoadConfigJSON(env)
	case "prod":
		conf, confErr = config.LoadConfigYAML()
	default:
		confErr = fmt.Errorf("unknown environment %q", env)
	}
	if confErr != nil {
		slog.Error("config load failed", "error", confErr)
		os.Exit(1)
	}

	psqlConf := conf.PostgresDB
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		psqlConf.Host, psqlConf.Port, psqlConf.User, psqlConf.Password, psqlConf.DBName, psqlConf.SSLMode)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		slog.Error("failed to open postgres connection", "error", err)
		os.Exit(1)
	}
	defer closeutil.Close(db)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		slog.Error("failed to reach postgres", "error", err)
		os.Exit(1)
	}

	urlRepo := repositories.NewCurrencySnapshotDataService(db)
	c := cache.InitLRUCache(100)

	idGen, err := generator.NewSnowflake(int64(conf.Server.DataCenterID), int64(conf.Server.MashineID))
	if err != nil {
		slog.Error("failed to init id generator", "error", err)
		os.Exit(1)
	}

	urlService := service.NewURLService(idGen, c, urlRepo)

	r := router.New()
	v1 := r.Group("/api/v1")

	h := handler.NewHandler(urlService)
	h.Register(v1)
	h.RegisterHealth(r)

	// Bind on all interfaces: inside a k8s pod, 127.0.0.1-only binding is
	// unreachable from the Service/kubelet probes.
	addr := fmt.Sprintf("0.0.0.0:%s", conf.Server.RESTPort)

	go func() {
		slog.Info("starting server", "addr", addr, "environment", env)
		if err := r.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server stopped unexpectedly", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	slog.Info("shutdown signal received, draining connections")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := r.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	slog.Info("shutdown complete")
}
