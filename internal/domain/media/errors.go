package media

import (
	"net/http"

	"github.com/danielng/kin-core-svc/pkg/apperror"
)

var (
	ErrMediaNotFound = apperror.New(
		apperror.CodeNotFound,
		"media not found",
		http.StatusNotFound,
	)

	ErrInvalidMediaType = apperror.New(
		apperror.CodeInvalidMediaType,
		"invalid media type",
		http.StatusBadRequest,
	)

	ErrMediaTooLarge = apperror.New(
		apperror.CodeMediaTooLarge,
		"media file too large",
		http.StatusRequestEntityTooLarge,
	)

	ErrUnsupportedMimeType = apperror.New(
		apperror.CodeValidation,
		"unsupported file type",
		http.StatusBadRequest,
	)

	ErrUploadFailed = apperror.New(
		apperror.CodeInternal,
		"failed to upload media",
		http.StatusInternalServerError,
	)

	ErrDownloadFailed = apperror.New(
		apperror.CodeInternal,
		"failed to download media",
		http.StatusInternalServerError,
	)
)
