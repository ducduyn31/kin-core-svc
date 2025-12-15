package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danielng/kin-core-svc/internal/application/circle"
	"github.com/danielng/kin-core-svc/internal/application/user"
	"github.com/danielng/kin-core-svc/internal/config"
	"github.com/danielng/kin-core-svc/internal/infrastructure/auth"
	"github.com/danielng/kin-core-svc/internal/infrastructure/postgres"
	"github.com/danielng/kin-core-svc/internal/infrastructure/redis"
	"github.com/danielng/kin-core-svc/internal/infrastructure/telemetry"
	connectServer "github.com/danielng/kin-core-svc/internal/interfaces/connect"
)

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger := setupLogger(cfg.Logging)

	logger.Info("starting Kin API server",
		"version", Version,
		"git_commit", GitCommit,
		"build_time", BuildTime,
		"port", cfg.Server.Port,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	otel, err := telemetry.New(ctx, telemetry.Config{
		ServiceName:    cfg.Telemetry.ServiceName,
		ServiceVersion: Version,
		Environment:    os.Getenv("KIN_ENV"),
		OTLPEndpoint:   cfg.Telemetry.OTLPEndpoint,
		Enabled:        cfg.Telemetry.Enabled,
	})
	if err != nil {
		logger.Error("failed to initialize telemetry", "error", err)
		os.Exit(1)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := otel.Shutdown(shutdownCtx); err != nil {
			logger.Error("failed to shutdown telemetry", "error", err)
		}
	}()
	if cfg.Telemetry.Enabled {
		logger.Info("telemetry initialized", "endpoint", cfg.Telemetry.OTLPEndpoint)
	}

	db, err := postgres.NewDB(ctx, cfg.Database, postgres.DBOptions{
		EnableTracing: otel.Enabled(),
	})
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("connected to PostgreSQL")

	redisClient, err := redis.NewClient(ctx, cfg.Redis, redis.ClientOptions{
		EnableTracing: otel.Enabled(),
	})
	if err != nil {
		logger.Error("failed to connect to Redis", "error", err)
		os.Exit(1)
	}
	defer func() { _ = redisClient.Close() }()
	logger.Info("connected to Redis")

	auth0Validator := auth.NewAuth0Validator(cfg.Auth.Domain, cfg.Auth.Audience)

	userRepo := postgres.NewUserRepository(db)
	circleRepo := postgres.NewCircleRepository(db)
	_ = redis.NewPresenceRepository(redisClient)

	userService := user.NewService(userRepo, logger)
	circleService := circle.NewService(circleRepo, logger)

	server := connectServer.NewServer(connectServer.ServerConfig{
		Logger:         logger,
		Auth0Validator: auth0Validator,
		UserService:    userService,
		CircleService:  circleService,
		BuildInfo: connectServer.BuildInfo{
			Version:   Version,
			GitCommit: GitCommit,
			BuildTime: BuildTime,
		},
		HealthCheckers:   []connectServer.HealthChecker{db, redisClient},
		EnableTracing:    otel.Enabled(),
		EnableReflection: cfg.Server.EnableReflection,
	})

	httpServer := &http.Server{
		Addr:         cfg.Server.Address(),
		Handler:      server.Handler(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	errCh := make(chan error, 1)

	go func() {
		logger.Info("Connect server listening", "address", cfg.Server.Address())
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("server error: %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	var exitCode int
	select {
	case sig := <-quit:
		logger.Info("received shutdown signal", "signal", sig)
	case err := <-errCh:
		logger.Error("server error", "error", err)
		exitCode = 1
	}

	logger.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}

	logger.Info("server stopped")

	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

func setupLogger(cfg config.LoggingConfig) *slog.Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{}

	switch cfg.Level {
	case "debug":
		opts.Level = slog.LevelDebug
	case "warn":
		opts.Level = slog.LevelWarn
	case "error":
		opts.Level = slog.LevelError
	default:
		opts.Level = slog.LevelInfo
	}

	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}
