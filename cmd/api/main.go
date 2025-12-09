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
	httpRouter "github.com/danielng/kin-core-svc/internal/interfaces/http"
	"github.com/gin-gonic/gin"
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
		"host", cfg.Server.Host,
		"port", cfg.Server.Port,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := postgres.NewDB(ctx, cfg.Database)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("connected to PostgreSQL")

	redisClient, err := redis.NewClient(ctx, cfg.Redis)
	if err != nil {
		logger.Error("failed to connect to Redis", "error", err)
		os.Exit(1)
	}
	defer func() { _ = redisClient.Close() }()
	logger.Info("connected to Redis")

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

	auth0Validator := auth.NewAuth0Validator(cfg.Auth.Domain, cfg.Auth.Audience)

	userRepo := postgres.NewUserRepository(db)
	circleRepo := postgres.NewCircleRepository(db)
	_ = redis.NewPresenceRepository(redisClient)

	userService := user.NewService(userRepo, logger)
	circleService := circle.NewService(circleRepo, logger)

	if os.Getenv("KIN_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := httpRouter.NewRouter(httpRouter.RouterConfig{
		Logger:         logger,
		Auth0Validator: auth0Validator,
		UserService:    userService,
		CircleService:  circleService,
		ServiceName:    cfg.Telemetry.ServiceName,
		BuildInfo: httpRouter.BuildInfo{
			Version:   Version,
			GitCommit: GitCommit,
			BuildTime: BuildTime,
		},
		TelemetryEnable: cfg.Telemetry.Enabled,
		HealthCheckers:  []httpRouter.HealthChecker{db, redisClient},
	})

	server := &http.Server{
		Addr:         cfg.Server.Address(),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		logger.Info("HTTP server listening", "address", cfg.Server.Address())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped")
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
