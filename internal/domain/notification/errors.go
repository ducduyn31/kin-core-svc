package notification

import (
	"net/http"

	"github.com/danielng/kin-core-svc/pkg/apperror"
)

var (
	ErrNotificationNotFound = apperror.New(
		apperror.CodeNotFound,
		"notification not found",
		http.StatusNotFound,
	)

	ErrPreferencesNotFound = apperror.New(
		apperror.CodeNotFound,
		"notification preferences not found",
		http.StatusNotFound,
	)

	ErrPushFailed = apperror.New(
		apperror.CodeInternal,
		"failed to send push notification",
		http.StatusInternalServerError,
	)

	ErrInvalidNotificationType = apperror.New(
		apperror.CodeValidation,
		"invalid notification type",
		http.StatusBadRequest,
	)
)
