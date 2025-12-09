package conversation

import (
	"net/http"

	"github.com/danielng/kin-core-svc/pkg/apperror"
)

var (
	ErrConversationNotFound = apperror.New(
		apperror.CodeConversationNotFound,
		"conversation not found",
		http.StatusNotFound,
	)

	ErrNotParticipant = apperror.New(
		apperror.CodeForbidden,
		"not a participant in this conversation",
		http.StatusForbidden,
	)

	ErrAlreadyParticipant = apperror.New(
		apperror.CodeConflict,
		"already a participant in this conversation",
		http.StatusConflict,
	)

	ErrParticipantNotFound = apperror.New(
		apperror.CodeNotFound,
		"participant not found",
		http.StatusNotFound,
	)

	ErrCannotRemoveFromDirect = apperror.New(
		apperror.CodeBadRequest,
		"cannot remove participant from direct conversation",
		http.StatusBadRequest,
	)

	ErrDirectConversationExists = apperror.New(
		apperror.CodeConflict,
		"direct conversation already exists",
		http.StatusConflict,
	)

	ErrInvalidConversationType = apperror.New(
		apperror.CodeBadRequest,
		"invalid conversation type",
		http.StatusBadRequest,
	)
)
