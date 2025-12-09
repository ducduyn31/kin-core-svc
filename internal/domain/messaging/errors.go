package messaging

import (
	"net/http"

	"github.com/danielng/kin-core-svc/pkg/apperror"
)

var (
	ErrMessageNotFound = apperror.New(
		apperror.CodeMessageNotFound,
		"message not found",
		http.StatusNotFound,
	)

	ErrNotMessageSender = apperror.New(
		apperror.CodeForbidden,
		"only the sender can perform this action",
		http.StatusForbidden,
	)

	ErrMessageDeleted = apperror.New(
		apperror.CodeNotFound,
		"message has been deleted",
		http.StatusNotFound,
	)

	ErrCannotEditMessage = apperror.New(
		apperror.CodeBadRequest,
		"message cannot be edited",
		http.StatusBadRequest,
	)

	ErrEditWindowExpired = apperror.New(
		apperror.CodeBadRequest,
		"edit window has expired",
		http.StatusBadRequest,
	)

	ErrReactionAlreadyExists = apperror.New(
		apperror.CodeConflict,
		"reaction already exists",
		http.StatusConflict,
	)

	ErrReactionNotFound = apperror.New(
		apperror.CodeNotFound,
		"reaction not found",
		http.StatusNotFound,
	)

	ErrInvalidContentType = apperror.New(
		apperror.CodeValidation,
		"invalid content type",
		http.StatusBadRequest,
	)

	ErrEmptyMessage = apperror.New(
		apperror.CodeValidation,
		"message content cannot be empty",
		http.StatusBadRequest,
	)
)
