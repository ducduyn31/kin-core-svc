package response

import (
	"net/http"

	"github.com/danielng/kin-core-svc/pkg/apperror"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Data any   `json:"data"`
	Meta *Meta `json:"meta,omitempty"`
}

type Meta struct {
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"per_page,omitempty"`
	Total   int64  `json:"total,omitempty"`
	Cursor  string `json:"cursor,omitempty"`
	HasMore bool   `json:"has_more,omitempty"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{Data: data})
}

func OKWithMeta(c *gin.Context, data any, meta *Meta) {
	c.JSON(http.StatusOK, Response{Data: data, Meta: meta})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Response{Data: data})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Error(c *gin.Context, err error) {
	if appErr, ok := apperror.AsAppError(err); ok {
		c.JSON(appErr.HTTPStatus, ErrorResponse{
			Error: ErrorBody{
				Code:    string(appErr.Code),
				Message: appErr.Message,
				Details: appErr.Details,
			},
		})
		return
	}

	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: ErrorBody{
			Code:    string(apperror.CodeInternal),
			Message: "internal server error",
		},
	})
}

func BadRequest(c *gin.Context, message string, details ...string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error: ErrorBody{
			Code:    string(apperror.CodeBadRequest),
			Message: message,
			Details: details,
		},
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{
		Error: ErrorBody{
			Code:    string(apperror.CodeUnauthorized),
			Message: message,
		},
	})
}

func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, ErrorResponse{
		Error: ErrorBody{
			Code:    string(apperror.CodeForbidden),
			Message: message,
		},
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Error: ErrorBody{
			Code:    string(apperror.CodeNotFound),
			Message: message,
		},
	})
}

func ValidationError(c *gin.Context, details []string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error: ErrorBody{
			Code:    string(apperror.CodeValidation),
			Message: "validation failed",
			Details: details,
		},
	})
}

func InternalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: ErrorBody{
			Code:    string(apperror.CodeInternal),
			Message: "internal server error",
		},
	})
}

func PaginationMeta(page, perPage int, total int64) *Meta {
	return &Meta{
		Page:    page,
		PerPage: perPage,
		Total:   total,
	}
}

func CursorMeta(cursor string, hasMore bool) *Meta {
	return &Meta{
		Cursor:  cursor,
		HasMore: hasMore,
	}
}
