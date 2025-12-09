package user

import (
	"net/http"

	"github.com/danielng/kin-core-svc/pkg/apperror"
)

var (
	ErrUserNotFound = apperror.New(
		apperror.CodeUserNotFound,
		"user not found",
		http.StatusNotFound,
	)

	ErrUserAlreadyExists = apperror.New(
		apperror.CodeUserAlreadyExists,
		"user already exists",
		http.StatusConflict,
	)

	ErrInvalidTimezone = apperror.New(
		apperror.CodeValidation,
		"invalid timezone",
		http.StatusBadRequest,
	)

	ErrInvalidPrivacyLevel = apperror.New(
		apperror.CodeValidation,
		"invalid privacy level",
		http.StatusBadRequest,
	)

	ErrPreferencesNotFound = apperror.New(
		apperror.CodeNotFound,
		"user preferences not found",
		http.StatusNotFound,
	)
)
