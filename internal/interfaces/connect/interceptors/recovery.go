package interceptors

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"

	"connectrpc.com/connect"
)

type RecoveryInterceptor struct {
	logger *slog.Logger
}

// NewRecoveryInterceptor constructs a RecoveryInterceptor that recovers from panics in Connect handlers.
// It uses the provided logger to record panic details (panic value, procedure, and stack trace) and causes panics to be converted into internal server errors.
func NewRecoveryInterceptor(logger *slog.Logger) *RecoveryInterceptor {
	return &RecoveryInterceptor{
		logger: logger,
	}
}

func (i *RecoveryInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (resp connect.AnyResponse, err error) {
		defer func() {
			if r := recover(); r != nil {
				i.logger.Error("panic recovered in unary handler",
					"panic", r,
					"procedure", req.Spec().Procedure,
					"stack", string(debug.Stack()),
				)
				err = connect.NewError(connect.CodeInternal, fmt.Errorf("internal server error"))
			}
		}()
		return next(ctx, req)
	}
}

func (i *RecoveryInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (i *RecoveryInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(ctx context.Context, conn connect.StreamingHandlerConn) (err error) {
		defer func() {
			if r := recover(); r != nil {
				i.logger.Error("panic recovered in streaming handler",
					"panic", r,
					"procedure", conn.Spec().Procedure,
					"stack", string(debug.Stack()),
				)
				err = connect.NewError(connect.CodeInternal, fmt.Errorf("internal server error"))
			}
		}()
		return next(ctx, conn)
	}
}