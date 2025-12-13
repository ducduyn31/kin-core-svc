package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/danielng/kin-core-svc/pkg/apperror"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}

	// Check for context cancellation/timeout first
	if errors.Is(err, context.Canceled) {
		return status.Error(codes.Canceled, "request canceled")
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return status.Error(codes.DeadlineExceeded, "request deadline exceeded")
	}

	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		grpcCode := httpStatusToGRPCCode(appErr.HTTPStatus)
		st := status.New(grpcCode, appErr.Message)

		// Attach AppError.Code and AppError.Details as gRPC error details
		errorInfo := &errdetails.ErrorInfo{
			Reason: string(appErr.Code),
			Domain: "kin.api",
		}

		if len(appErr.Details) > 0 {
			errorInfo.Metadata = make(map[string]string)
			for i, detail := range appErr.Details {
				errorInfo.Metadata[fmt.Sprintf("detail_%d", i)] = detail
			}
		}

		stWithDetails, detailErr := st.WithDetails(errorInfo)
		if detailErr != nil {
			// Fall back to status without details if attachment fails
			return st.Err()
		}
		return stWithDetails.Err()
	}

	return status.Error(codes.Internal, "internal server error")
}

func httpStatusToGRPCCode(httpStatus int) codes.Code {
	switch httpStatus {
	case http.StatusOK:
		return codes.OK
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusConflict:
		return codes.AlreadyExists
	case http.StatusGone:
		return codes.NotFound
	case http.StatusTooManyRequests:
		return codes.ResourceExhausted
	case http.StatusRequestEntityTooLarge:
		return codes.InvalidArgument
	case http.StatusInternalServerError:
		return codes.Internal
	case http.StatusNotImplemented:
		return codes.Unimplemented
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	default:
		return codes.Unknown
	}
}
