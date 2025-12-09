package availability

import (
	"net/http"

	"github.com/danielng/kin-core-svc/pkg/apperror"
)

var (
	ErrAvailabilityNotFound = apperror.New(
		apperror.CodeNotFound,
		"availability not found",
		http.StatusNotFound,
	)

	ErrWindowNotFound = apperror.New(
		apperror.CodeNotFound,
		"availability window not found",
		http.StatusNotFound,
	)

	ErrAutoRuleNotFound = apperror.New(
		apperror.CodeNotFound,
		"auto rule not found",
		http.StatusNotFound,
	)

	ErrInvalidStatus = apperror.New(
		apperror.CodeValidation,
		"invalid availability status",
		http.StatusBadRequest,
	)

	ErrInvalidTimeFormat = apperror.New(
		apperror.CodeValidation,
		"invalid time format, expected HH:MM",
		http.StatusBadRequest,
	)

	ErrInvalidWeekday = apperror.New(
		apperror.CodeValidation,
		"invalid weekday",
		http.StatusBadRequest,
	)

	ErrWindowOverlap = apperror.New(
		apperror.CodeConflict,
		"availability window overlaps with existing window",
		http.StatusConflict,
	)
)
