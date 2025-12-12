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
	grpcServer "github.com/danielng/kin-core-svc/internal/interfaces/grpc"
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
		"grpc_port", cfg.GRPC.Port,
		"gateway_port", cfg.GRPC.GatewayPort,
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

	grpc := grpcServer.NewServer(grpcServer.ServerConfig{
		Logger:           logger,
		Auth0Validator:   auth0Validator,
		UserService:      userService,
		CircleService:    circleService,
		EnableReflection: cfg.GRPC.EnableReflection,
	})

	go func() {
		if err := grpc.Serve(cfg.GRPC.Address()); err != nil {
			logger.Error("gRPC server error", "error", err)
			os.Exit(1)
		}
	}()

	gatewayCtx := context.Background()
	gateway, err := grpcServer.NewGatewayServer(gatewayCtx, grpcServer.GatewayConfig{
		GRPCAddress: cfg.GRPC.Address(),
		Logger:      logger,
		BuildInfo: grpcServer.BuildInfo{
			Version:   Version,
			GitCommit: GitCommit,
			BuildTime: BuildTime,
		},
		HealthCheckers: []grpcServer.HealthChecker{db, redisClient},
	})
	if err != nil {
		logger.Error("failed to create gRPC gateway", "error", err)
		os.Exit(1)
	}

	gatewayServer := &http.Server{
		Addr:         cfg.GRPC.GatewayAddress(),
		Handler:      gateway.Handler(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		logger.Info("gRPC-Gateway listening", "address", cfg.GRPC.GatewayAddress())
		if err := gatewayServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("gRPC-Gateway error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down servers...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer shutdownCancel()

	if err := gatewayServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("gRPC-Gateway forced to shutdown", "error", err)
	}

	grpc.GracefulStop()

	logger.Info("servers stopped")
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
