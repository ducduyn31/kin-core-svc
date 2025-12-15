package connect

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"github.com/danielng/kin-core-svc/gen/proto/kin/v1/kinv1connect"
	"github.com/danielng/kin-core-svc/internal/application/circle"
	"github.com/danielng/kin-core-svc/internal/application/user"
	"github.com/danielng/kin-core-svc/internal/infrastructure/auth"
	"github.com/danielng/kin-core-svc/internal/interfaces/connect/handlers"
	"github.com/danielng/kin-core-svc/internal/interfaces/connect/interceptors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type HealthChecker interface {
	Ping(ctx context.Context) error
	Name() string
}

type BuildInfo struct {
	Version   string
	GitCommit string
	BuildTime string
}

type ServerConfig struct {
	Logger           *slog.Logger
	Auth0Validator   *auth.Auth0Validator
	UserService      *user.Service
	CircleService    *circle.Service
	BuildInfo        BuildInfo
	HealthCheckers   []HealthChecker
	EnableTracing    bool
	EnableReflection bool // Enable gRPC reflection (for development only)
}

type Server struct {
	handler        http.Handler
	logger         *slog.Logger
	buildInfo      BuildInfo
	healthCheckers []HealthChecker
}

func NewServer(cfg ServerConfig) *Server {
	mux := http.NewServeMux()

	authInterceptor := interceptors.NewAuthInterceptor(cfg.Auth0Validator, cfg.UserService)
	recoveryInterceptor := interceptors.NewRecoveryInterceptor(cfg.Logger)

	handlerOpts := []connect.HandlerOption{
		connect.WithInterceptors(
			recoveryInterceptor,
			authInterceptor,
		),
	}

	userHandler := handlers.NewUserHandler(cfg.UserService)
	circleHandler := handlers.NewCircleHandler(cfg.CircleService)

	path, handler := kinv1connect.NewUserServiceHandler(userHandler, handlerOpts...)
	mux.Handle(path, handler)

	path, handler = kinv1connect.NewCircleServiceHandler(circleHandler, handlerOpts...)
	mux.Handle(path, handler)

	// Enable gRPC reflection for development (allows service discovery in gRPC clients)
	if cfg.EnableReflection {
		reflector := grpcreflect.NewStaticReflector(
			kinv1connect.UserServiceName,
			kinv1connect.CircleServiceName,
		)
		mux.Handle(grpcreflect.NewHandlerV1(reflector))
		mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
		cfg.Logger.Info("gRPC reflection enabled")
	}

	server := &Server{
		logger:         cfg.Logger,
		buildInfo:      cfg.BuildInfo,
		healthCheckers: cfg.HealthCheckers,
	}

	mux.HandleFunc("/health", server.healthHandler)
	mux.HandleFunc("/ready", server.readyHandler)

	server.handler = h2c.NewHandler(mux, &http2.Server{})

	return server
}

func (s *Server) Handler() http.Handler {
	return s.handler
}

type healthResponse struct {
	Status       string                  `json:"status"`
	Dependencies map[string]healthStatus `json:"dependencies,omitempty"`
}

type healthStatus struct {
	Status  string `json:"status"`
	Latency string `json:"latency,omitempty"`
	Error   string `json:"error,omitempty"`
}

type readyResponse struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
	BuildTime string `json:"build_time"`
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	dependencies := make(map[string]healthStatus)
	var mu sync.Mutex
	var wg sync.WaitGroup

	allHealthy := true
	var healthyMu sync.Mutex

	for _, checker := range s.healthCheckers {
		wg.Add(1)
		go func(chk HealthChecker) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
			defer cancel()

			start := time.Now()
			err := chk.Ping(ctx)
			latency := time.Since(start)

			status := healthStatus{
				Status:  "healthy",
				Latency: latency.String(),
			}

			if err != nil {
				status.Status = "unhealthy"
				status.Error = err.Error()
				healthyMu.Lock()
				allHealthy = false
				healthyMu.Unlock()
			}

			mu.Lock()
			dependencies[chk.Name()] = status
			mu.Unlock()
		}(checker)
	}

	wg.Wait()

	overallStatus := "healthy"
	httpStatus := http.StatusOK
	if !allHealthy {
		overallStatus = "unhealthy"
		httpStatus = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	if err := json.NewEncoder(w).Encode(healthResponse{
		Status:       overallStatus,
		Dependencies: dependencies,
	}); err != nil {
		slog.Error("failed to encode health response", "error", err)
	}
}

func (s *Server) readyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(readyResponse{
		Status:    "ready",
		Version:   s.buildInfo.Version,
		GitCommit: s.buildInfo.GitCommit,
		BuildTime: s.buildInfo.BuildTime,
	}); err != nil {
		slog.Error("failed to encode ready response", "error", err)
	}
}
