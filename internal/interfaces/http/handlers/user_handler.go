package handlers

import (
	"github.com/danielng/kin-core-svc/internal/application/user"
	userDomain "github.com/danielng/kin-core-svc/internal/domain/user"
	"github.com/danielng/kin-core-svc/pkg/ctxkey"
	"github.com/danielng/kin-core-svc/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService *user.Service
	validate    *validator.Validate
}

func NewUserHandler(userService *user.Service) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validator.New(),
	}
}

// GET /users/me
func (h *UserHandler) GetMe(c *gin.Context) {
	u := c.MustGet(ctxkey.User).(*userDomain.User)
	response.OK(c, u)
}

type UpdateProfileRequest struct {
	DisplayName string  `json:"display_name" validate:"omitempty,min=1,max=100"`
	Bio         *string `json:"bio" validate:"omitempty,max=500"`
	Avatar      *string `json:"avatar" validate:"omitempty,url"`
}

// PUT /users/me
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, formatValidationErrors(err))
		return
	}

	userID := c.MustGet(ctxkey.UserID).(uuid.UUID)
	u, err := h.userService.UpdateProfile(c.Request.Context(), user.UpdateProfileCommand{
		UserID:      userID,
		DisplayName: req.DisplayName,
		Bio:         req.Bio,
		Avatar:      req.Avatar,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, u)
}

type UpdateTimezoneRequest struct {
	Timezone string `json:"timezone" validate:"required"`
}

// PUT /users/me/timezone
func (h *UserHandler) UpdateTimezone(c *gin.Context) {
	var req UpdateTimezoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, formatValidationErrors(err))
		return
	}

	userID := c.MustGet(ctxkey.UserID).(uuid.UUID)
	u, err := h.userService.UpdateTimezone(c.Request.Context(), user.UpdateTimezoneCommand{
		UserID:   userID,
		Timezone: req.Timezone,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, u)
}

// GET /users/me/preferences
func (h *UserHandler) GetPreferences(c *gin.Context) {
	userID := c.MustGet(ctxkey.UserID).(uuid.UUID)
	prefs, err := h.userService.GetPreferences(c.Request.Context(), user.GetUserPreferencesQuery{
		UserID: userID,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, prefs)
}

type UpdatePreferencesRequest struct {
	DefaultPrivacyLevel   *string `json:"default_privacy_level"`
	ShowOnlineStatus      *bool   `json:"show_online_status"`
	ShowLastSeen          *bool   `json:"show_last_seen"`
	ShowReadReceipts      *bool   `json:"show_read_receipts"`
	AllowContactDiscovery *bool   `json:"allow_contact_discovery"`
	PushNotifications     *bool   `json:"push_notifications"`
	EmailNotifications    *bool   `json:"email_notifications"`
	QuietHoursEnabled     *bool   `json:"quiet_hours_enabled"`
	QuietHoursStart       *string `json:"quiet_hours_start"`
	QuietHoursEnd         *string `json:"quiet_hours_end"`
}

// PUT /users/me/preferences
func (h *UserHandler) UpdatePreferences(c *gin.Context) {
	var req UpdatePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	userID := c.MustGet(ctxkey.UserID).(uuid.UUID)
	prefs, err := h.userService.UpdatePreferences(c.Request.Context(), user.UpdatePreferencesCommand{
		UserID:                userID,
		DefaultPrivacyLevel:   req.DefaultPrivacyLevel,
		ShowOnlineStatus:      req.ShowOnlineStatus,
		ShowLastSeen:          req.ShowLastSeen,
		ShowReadReceipts:      req.ShowReadReceipts,
		AllowContactDiscovery: req.AllowContactDiscovery,
		PushNotifications:     req.PushNotifications,
		EmailNotifications:    req.EmailNotifications,
		QuietHoursEnabled:     req.QuietHoursEnabled,
		QuietHoursStart:       req.QuietHoursStart,
		QuietHoursEnd:         req.QuietHoursEnd,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, prefs)
}

func formatValidationErrors(err error) []string {
	var errors []string
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			errors = append(errors, e.Field()+": "+e.Tag())
		}
	}
	return errors
}
