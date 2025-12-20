package handlers

import (
	"context"
	"errors"
	"net/http"

	"connectrpc.com/connect"
	"github.com/danielng/kin-core-svc/pkg/apperror"
)

// mapError translates an input error into a Connect protocol error suitable for responses.
// It returns nil when the input is nil. Context cancellation and deadline errors are mapped
// to CodeCanceled and CodeDeadlineExceeded with corresponding short messages. If the error
// is an *apperror.AppError its HTTPStatus is converted to a Connect code and the AppError
// message is preserved. All other errors are mapped to CodeInternal with a generic message.
func mapError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, context.Canceled) {
		return connect.NewError(connect.CodeCanceled, errors.New("request canceled"))
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return connect.NewError(connect.CodeDeadlineExceeded, errors.New("request deadline exceeded"))
	}

	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		code := httpStatusToConnectCode(appErr.HTTPStatus)
		return connect.NewError(code, errors.New(appErr.Message))
	}

	return connect.NewError(connect.CodeInternal, errors.New("internal server error"))
}

// httpStatusToConnectCode maps an HTTP status code to the corresponding connect.Code.
// Common statuses are translated (400 → InvalidArgument, 401 → Unauthenticated, 403 → PermissionDenied,
// 404/410 → NotFound, 409 → AlreadyExists, 429 → ResourceExhausted, 413 → InvalidArgument,
// 500 → Internal, 501 → Unimplemented, 503 → Unavailable). Unknown or unmapped statuses return CodeUnknown.
func httpStatusToConnectCode(httpStatus int) connect.Code {
	switch httpStatus {
	case http.StatusBadRequest:
		return connect.CodeInvalidArgument
	case http.StatusUnauthorized:
		return connect.CodeUnauthenticated
	case http.StatusForbidden:
		return connect.CodePermissionDenied
	case http.StatusNotFound:
		return connect.CodeNotFound
	case http.StatusConflict:
		return connect.CodeAlreadyExists
	case http.StatusGone:
		return connect.CodeNotFound
	case http.StatusTooManyRequests:
		return connect.CodeResourceExhausted
	case http.StatusRequestEntityTooLarge:
		return connect.CodeInvalidArgument
	case http.StatusInternalServerError:
		return connect.CodeInternal
	case http.StatusNotImplemented:
		return connect.CodeUnimplemented
	case http.StatusServiceUnavailable:
		return connect.CodeUnavailable
	default:
		return connect.CodeUnknown
	}
}