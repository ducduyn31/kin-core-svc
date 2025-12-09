package circle

import (
	"net/http"

	"github.com/danielng/kin-core-svc/pkg/apperror"
)

var (
	ErrCircleNotFound = apperror.New(
		apperror.CodeCircleNotFound,
		"circle not found",
		http.StatusNotFound,
	)

	ErrNotCircleMember = apperror.New(
		apperror.CodeNotCircleMember,
		"not a member of this circle",
		http.StatusForbidden,
	)

	ErrNotCircleAdmin = apperror.New(
		apperror.CodeNotCircleAdmin,
		"admin privileges required",
		http.StatusForbidden,
	)

	ErrMemberNotFound = apperror.New(
		apperror.CodeCircleMemberNotFound,
		"circle member not found",
		http.StatusNotFound,
	)

	ErrAlreadyMember = apperror.New(
		apperror.CodeConflict,
		"already a member of this circle",
		http.StatusConflict,
	)

	ErrInvitationNotFound = apperror.New(
		apperror.CodeInvitationNotFound,
		"invitation not found",
		http.StatusNotFound,
	)

	ErrInvitationExpired = apperror.New(
		apperror.CodeInvitationExpired,
		"invitation has expired",
		http.StatusGone,
	)

	ErrInvitationInvalid = apperror.New(
		apperror.CodeBadRequest,
		"invitation is not valid",
		http.StatusBadRequest,
	)

	ErrCannotRemoveLastAdmin = apperror.New(
		apperror.CodeBadRequest,
		"cannot remove the last admin from circle",
		http.StatusBadRequest,
	)

	ErrCannotLeaveAsLastAdmin = apperror.New(
		apperror.CodeBadRequest,
		"cannot leave circle as the last admin",
		http.StatusBadRequest,
	)

	ErrInvalidRole = apperror.New(
		apperror.CodeValidation,
		"invalid member role",
		http.StatusBadRequest,
	)

	ErrSharingPreferenceNotFound = apperror.New(
		apperror.CodeNotFound,
		"sharing preference not found",
		http.StatusNotFound,
	)
)
