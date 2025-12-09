package contact

import (
	"net/http"

	"github.com/danielng/kin-core-svc/pkg/apperror"
)

var (
	ErrContactNotFound = apperror.New(
		apperror.CodeContactNotFound,
		"contact not found",
		http.StatusNotFound,
	)

	ErrContactAlreadyExists = apperror.New(
		apperror.CodeContactAlreadyExists,
		"contact already exists",
		http.StatusConflict,
	)

	ErrContactRequestNotFound = apperror.New(
		apperror.CodeNotFound,
		"contact request not found",
		http.StatusNotFound,
	)

	ErrContactRequestAlreadyExists = apperror.New(
		apperror.CodeConflict,
		"contact request already exists",
		http.StatusConflict,
	)

	ErrContactRequestNotPending = apperror.New(
		apperror.CodeBadRequest,
		"contact request is not pending",
		http.StatusBadRequest,
	)

	ErrCannotAddSelf = apperror.New(
		apperror.CodeBadRequest,
		"cannot add yourself as a contact",
		http.StatusBadRequest,
	)

	ErrContactBlocked = apperror.New(
		apperror.CodeForbidden,
		"contact is blocked",
		http.StatusForbidden,
	)
)
