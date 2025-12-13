package grpc

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	kinv1 "github.com/danielng/kin-core-svc/gen/proto/kin/v1"
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

type GatewayConfig struct {
	GRPCAddress    string
	Logger         *slog.Logger
	BuildInfo      BuildInfo
	HealthCheckers []HealthChecker
}

type GatewayServer struct {
	handler        http.Handler
	logger         *slog.Logger
	buildInfo      BuildInfo
	healthCheckers []HealthChecker
}

func NewGatewayServer(ctx context.Context, cfg GatewayConfig) (*GatewayServer, error) {
	gwMux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(customHeaderMatcher),
	)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := kinv1.RegisterUserServiceHandlerFromEndpoint(ctx, gwMux, cfg.GRPCAddress, opts); err != nil {
		return nil, err
	}

	if err := kinv1.RegisterCircleServiceHandlerFromEndpoint(ctx, gwMux, cfg.GRPCAddress, opts); err != nil {
		return nil, err
	}

	cfg.Logger.Info("gRPC-Gateway initialized", "grpc_address", cfg.GRPCAddress)

	server := &GatewayServer{
		logger:         cfg.Logger,
		buildInfo:      cfg.BuildInfo,
		healthCheckers: cfg.HealthCheckers,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", server.healthHandler)
	mux.HandleFunc("/ready", server.readyHandler)
	mux.Handle("/", gwMux)

	server.handler = mux

	return server, nil
}

func (s *GatewayServer) Handler() http.Handler {
	return s.handler
}

func customHeaderMatcher(key string) (string, bool) {
	switch key {
	case "Authorization", "Content-Type", "Accept":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
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

func (s *GatewayServer) healthHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *GatewayServer) readyHandler(w http.ResponseWriter, r *http.Request) {
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
