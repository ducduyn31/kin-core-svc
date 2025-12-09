package location

import (
	"net/http"

	"github.com/danielng/kin-core-svc/pkg/apperror"
)

var (
	ErrLocationNotFound = apperror.New(
		apperror.CodeNotFound,
		"location not found",
		http.StatusNotFound,
	)

	ErrPlaceNotFound = apperror.New(
		apperror.CodeNotFound,
		"place not found",
		http.StatusNotFound,
	)

	ErrCheckInNotFound = apperror.New(
		apperror.CodeNotFound,
		"check-in not found",
		http.StatusNotFound,
	)

	ErrInvalidCoordinates = apperror.New(
		apperror.CodeValidation,
		"invalid coordinates",
		http.StatusBadRequest,
	)

	ErrInvalidPlaceType = apperror.New(
		apperror.CodeValidation,
		"invalid place type",
		http.StatusBadRequest,
	)

	ErrInvalidPrecision = apperror.New(
		apperror.CodeValidation,
		"invalid location precision",
		http.StatusBadRequest,
	)

	ErrLocationSharingDisabled = apperror.New(
		apperror.CodeForbidden,
		"location sharing is disabled",
		http.StatusForbidden,
	)
)
