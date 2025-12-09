package handlers

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
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

type HealthHandler struct {
	buildInfo BuildInfo
	checkers  []HealthChecker
}

func NewHealthHandler(buildInfo BuildInfo, checkers ...HealthChecker) *HealthHandler {
	return &HealthHandler{
		buildInfo: buildInfo,
		checkers:  checkers,
	}
}

type HealthResponse struct {
	Status       string                  `json:"status"`
	Dependencies map[string]HealthStatus `json:"dependencies,omitempty"`
}

type ReadyResponse struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
	BuildTime string `json:"build_time"`
}

type HealthStatus struct {
	Status  string `json:"status"`
	Latency string `json:"latency,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (h *HealthHandler) Health(c *gin.Context) {
	dependencies := make(map[string]HealthStatus)
	var mu sync.Mutex
	var wg sync.WaitGroup

	allHealthy := true
	var healthyMu sync.Mutex

	for _, checker := range h.checkers {
		wg.Add(1)
		go func(chk HealthChecker) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
			defer cancel()

			start := time.Now()
			err := chk.Ping(ctx)
			latency := time.Since(start)

			status := HealthStatus{
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

	c.JSON(httpStatus, HealthResponse{
		Status:       overallStatus,
		Dependencies: dependencies,
	})
}

func (h *HealthHandler) Ready(c *gin.Context) {
	c.JSON(http.StatusOK, ReadyResponse{
		Status:    "ready",
		Version:   h.buildInfo.Version,
		GitCommit: h.buildInfo.GitCommit,
		BuildTime: h.buildInfo.BuildTime,
	})
}
