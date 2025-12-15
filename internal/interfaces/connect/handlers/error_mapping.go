package handlers

import (
	"context"
	"errors"
	"net/http"

	"connectrpc.com/connect"
	"github.com/danielng/kin-core-svc/pkg/apperror"
)

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
