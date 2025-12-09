package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/danielng/kin-core-svc/pkg/response"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logAttrs := []any{
					"error", err,
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
					"stack", string(debug.Stack()),
				}

				spanCtx := trace.SpanContextFromContext(c.Request.Context())
				if spanCtx.HasTraceID() {
					logAttrs = append(logAttrs, "trace_id", spanCtx.TraceID().String())
				}

				logger.Error("panic recovered", logAttrs...)

				c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse{
					Error: response.ErrorBody{
						Code:    "INTERNAL_ERROR",
						Message: "internal server error",
					},
				})
			}
		}()

		c.Next()
	}
}
