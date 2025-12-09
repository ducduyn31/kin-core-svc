package presence

import (
	"net/http"

	"github.com/danielng/kin-core-svc/pkg/apperror"
)

var (
	ErrPresenceNotFound = apperror.New(
		apperror.CodeNotFound,
		"presence not found",
		http.StatusNotFound,
	)

	ErrActivityNotFound = apperror.New(
		apperror.CodeNotFound,
		"activity not found",
		http.StatusNotFound,
	)

	ErrInvalidDeviceType = apperror.New(
		apperror.CodeValidation,
		"invalid device type",
		http.StatusBadRequest,
	)

	ErrInvalidActivityType = apperror.New(
		apperror.CodeValidation,
		"invalid activity type",
		http.StatusBadRequest,
	)
)
