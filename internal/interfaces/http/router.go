package http

import (
	"log/slog"

	"github.com/danielng/kin-core-svc/internal/application/circle"
	"github.com/danielng/kin-core-svc/internal/application/user"
	"github.com/danielng/kin-core-svc/internal/infrastructure/auth"
	"github.com/danielng/kin-core-svc/internal/interfaces/http/handlers"
	"github.com/danielng/kin-core-svc/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type (
	HealthChecker = handlers.HealthChecker
	BuildInfo     = handlers.BuildInfo
)

type RouterConfig struct {
	Logger          *slog.Logger
	Auth0Validator  *auth.Auth0Validator
	UserService     *user.Service
	CircleService   *circle.Service
	ServiceName     string
	BuildInfo       BuildInfo
	TelemetryEnable bool
	HealthCheckers  []HealthChecker
}

func NewRouter(cfg RouterConfig) *gin.Engine {
	router := gin.New()

	router.Use(middleware.Recovery(cfg.Logger))

	if cfg.TelemetryEnable {
		router.Use(otelgin.Middleware(cfg.ServiceName))
	}

	healthHandler := handlers.NewHealthHandler(cfg.BuildInfo, cfg.HealthCheckers...)
	userHandler := handlers.NewUserHandler(cfg.UserService)
	circleHandler := handlers.NewCircleHandler(cfg.CircleService)

	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)

	v1 := router.Group("/api/v1")

	authMiddleware := middleware.Auth(cfg.Auth0Validator, cfg.UserService)

	users := v1.Group("/users")
	users.Use(authMiddleware)
	{
		users.GET("/me", userHandler.GetMe)
		users.PUT("/me", userHandler.UpdateProfile)
		users.PUT("/me/timezone", userHandler.UpdateTimezone)
		users.GET("/me/preferences", userHandler.GetPreferences)
		users.PUT("/me/preferences", userHandler.UpdatePreferences)
	}

	circles := v1.Group("/circles")
	circles.Use(authMiddleware)
	{
		circles.POST("", circleHandler.CreateCircle)
		circles.GET("", circleHandler.ListCircles)
		circles.POST("/join", circleHandler.AcceptInvitation)
		circles.GET("/:id", circleHandler.GetCircle)
		circles.PUT("/:id", circleHandler.UpdateCircle)
		circles.DELETE("/:id", circleHandler.DeleteCircle)
		circles.GET("/:id/members", circleHandler.ListMembers)
		circles.POST("/:id/members", circleHandler.AddMember)
		circles.DELETE("/:id/members/:memberId", circleHandler.RemoveMember)
		circles.POST("/:id/leave", circleHandler.LeaveCircle)
		circles.GET("/:id/sharing", circleHandler.GetSharingPreference)
		circles.PUT("/:id/sharing", circleHandler.UpdateSharingPreference)
		circles.POST("/:id/invitations", circleHandler.CreateInvitation)
	}

	return router
}
