package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	// General errors
	CodeInternal       ErrorCode = "INTERNAL_ERROR"
	CodeValidation     ErrorCode = "VALIDATION_ERROR"
	CodeNotFound       ErrorCode = "NOT_FOUND"
	CodeConflict       ErrorCode = "CONFLICT"
	CodeUnauthorized   ErrorCode = "UNAUTHORIZED"
	CodeForbidden      ErrorCode = "FORBIDDEN"
	CodeBadRequest     ErrorCode = "BAD_REQUEST"
	CodeTooManyRequest ErrorCode = "TOO_MANY_REQUESTS"

	// Domain-specific errors
	CodeUserNotFound         ErrorCode = "USER_NOT_FOUND"
	CodeUserAlreadyExists    ErrorCode = "USER_ALREADY_EXISTS"
	CodeCircleNotFound       ErrorCode = "CIRCLE_NOT_FOUND"
	CodeCircleMemberNotFound ErrorCode = "CIRCLE_MEMBER_NOT_FOUND"
	CodeNotCircleMember      ErrorCode = "NOT_CIRCLE_MEMBER"
	CodeNotCircleAdmin       ErrorCode = "NOT_CIRCLE_ADMIN"
	CodeConversationNotFound ErrorCode = "CONVERSATION_NOT_FOUND"
	CodeMessageNotFound      ErrorCode = "MESSAGE_NOT_FOUND"
	CodeContactNotFound      ErrorCode = "CONTACT_NOT_FOUND"
	CodeContactAlreadyExists ErrorCode = "CONTACT_ALREADY_EXISTS"
	CodeInvitationNotFound   ErrorCode = "INVITATION_NOT_FOUND"
	CodeInvitationExpired    ErrorCode = "INVITATION_EXPIRED"
	CodeInvalidMediaType     ErrorCode = "INVALID_MEDIA_TYPE"
	CodeMediaTooLarge        ErrorCode = "MEDIA_TOO_LARGE"
)

type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    []string  `json:"details,omitempty"`
	HTTPStatus int       `json:"-"`
	Err        error     `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) WithDetails(details ...string) *AppError {
	e.Details = append(e.Details, details...)
	return e
}

func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

func New(code ErrorCode, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// Common error constructors

func Internal(message string) *AppError {
	return New(CodeInternal, message, http.StatusInternalServerError)
}

func Validation(message string) *AppError {
	return New(CodeValidation, message, http.StatusBadRequest)
}

func NotFound(message string) *AppError {
	return New(CodeNotFound, message, http.StatusNotFound)
}

func Conflict(message string) *AppError {
	return New(CodeConflict, message, http.StatusConflict)
}

func Unauthorized(message string) *AppError {
	return New(CodeUnauthorized, message, http.StatusUnauthorized)
}

func Forbidden(message string) *AppError {
	return New(CodeForbidden, message, http.StatusForbidden)
}

func BadRequest(message string) *AppError {
	return New(CodeBadRequest, message, http.StatusBadRequest)
}

func TooManyRequests(message string) *AppError {
	return New(CodeTooManyRequest, message, http.StatusTooManyRequests)
}

// Domain-specific error constructors

func UserNotFound() *AppError {
	return New(CodeUserNotFound, "user not found", http.StatusNotFound)
}

func UserAlreadyExists() *AppError {
	return New(CodeUserAlreadyExists, "user already exists", http.StatusConflict)
}

func CircleNotFound() *AppError {
	return New(CodeCircleNotFound, "circle not found", http.StatusNotFound)
}

func NotCircleMember() *AppError {
	return New(CodeNotCircleMember, "not a member of this circle", http.StatusForbidden)
}

func NotCircleAdmin() *AppError {
	return New(CodeNotCircleAdmin, "admin privileges required", http.StatusForbidden)
}

func ConversationNotFound() *AppError {
	return New(CodeConversationNotFound, "conversation not found", http.StatusNotFound)
}

func MessageNotFound() *AppError {
	return New(CodeMessageNotFound, "message not found", http.StatusNotFound)
}

func ContactNotFound() *AppError {
	return New(CodeContactNotFound, "contact not found", http.StatusNotFound)
}

func ContactAlreadyExists() *AppError {
	return New(CodeContactAlreadyExists, "contact already exists", http.StatusConflict)
}

func InvitationNotFound() *AppError {
	return New(CodeInvitationNotFound, "invitation not found", http.StatusNotFound)
}

func InvitationExpired() *AppError {
	return New(CodeInvitationExpired, "invitation has expired", http.StatusGone)
}

func InvalidMediaType() *AppError {
	return New(CodeInvalidMediaType, "invalid media type", http.StatusBadRequest)
}

func MediaTooLarge() *AppError {
	return New(CodeMediaTooLarge, "media file too large", http.StatusRequestEntityTooLarge)
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// AsAppError extracts AppError from an error
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// HTTPStatusFromError returns HTTP status from error, defaulting to 500
func HTTPStatusFromError(err error) int {
	if appErr, ok := AsAppError(err); ok {
		return appErr.HTTPStatus
	}
	return http.StatusInternalServerError
}
