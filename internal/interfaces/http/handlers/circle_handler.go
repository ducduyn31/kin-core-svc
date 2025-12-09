package handlers

import (
	"time"

	"github.com/danielng/kin-core-svc/internal/application/circle"
	circleDomain "github.com/danielng/kin-core-svc/internal/domain/circle"
	"github.com/danielng/kin-core-svc/internal/domain/user"
	"github.com/danielng/kin-core-svc/pkg/ctxkey"
	"github.com/danielng/kin-core-svc/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CircleHandler struct {
	circleService *circle.Service
	validate      *validator.Validate
}

func NewCircleHandler(circleService *circle.Service) *CircleHandler {
	return &CircleHandler{
		circleService: circleService,
		validate:      validator.New(),
	}
}

type CreateCircleRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// POST /circles
func (h *CircleHandler) CreateCircle(c *gin.Context) {
	var req CreateCircleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, formatValidationErrors(err))
		return
	}

	userID := c.MustGet(ctxkey.UserID).(uuid.UUID)
	circ, err := h.circleService.CreateCircle(c.Request.Context(), circle.CreateCircleCommand{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   userID,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, circ)
}

// GET /circles
func (h *CircleHandler) ListCircles(c *gin.Context) {
	userID := c.MustGet(ctxkey.UserID).(uuid.UUID)

	circles, total, err := h.circleService.ListUserCircles(c.Request.Context(), circle.ListUserCirclesQuery{
		UserID: userID,
		Limit:  20,
		Offset: 0,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OKWithMeta(c, circles, response.PaginationMeta(1, 20, total))
}

// GET /circles/:id
func (h *CircleHandler) GetCircle(c *gin.Context) {
	circleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid circle ID")
		return
	}

	userID := c.MustGet(ctxkey.UserID).(uuid.UUID)
	circ, err := h.circleService.GetCircle(c.Request.Context(), circle.GetCircleQuery{
		CircleID: circleID,
		UserID:   userID,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, circ)
}

type UpdateCircleRequest struct {
	Name        string  `json:"name" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description" validate:"omitempty,max=500"`
}

// PUT /circles/:id
func (h *CircleHandler) UpdateCircle(c *gin.Context) {
	circleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid circle ID")
		return
	}

	var req UpdateCircleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, formatValidationErrors(err))
		return
	}

	userID := c.MustGet(ctxkey.UserID).(uuid.UUID)
	circ, err := h.circleService.UpdateCircle(c.Request.Context(), circle.UpdateCircleCommand{
		CircleID:    circleID,
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, circ)
}

// DELETE /circles/:id
func (h *CircleHandler) DeleteCircle(c *gin.Context) {
	circleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid circle ID")
		return
	}

	userID := c.MustGet(ctxkey.UserID).(uuid.UUID)
	err = h.circleService.DeleteCircle(c.Request.Context(), circle.DeleteCircleCommand{
		CircleID: circleID,
		UserID:   userID,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.NoContent(c)
}

// GET /circles/:id/members
func (h *CircleHandler) ListMembers(c *gin.Context) {
	circleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid circle ID")
		return
	}

	userID := c.MustGet(ctxkey.UserID).(uuid.UUID)
	members, err := h.circleService.ListMembers(c.Request.Context(), circle.ListCircleMembersQuery{
		CircleID: circleID,
		UserID:   userID,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, members)
}

type AddMemberRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	Role   string `json:"role" validate:"required,oneof=admin member"`
}

// POST /circles/:id/members
func (h *CircleHandler) AddMember(c *gin.Context) {
	circleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid circle ID")
		return
	}

	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, formatValidationErrors(err))
		return
	}

	memberID, _ := uuid.Parse(req.UserID)
	userID := c.MustGet(ctxkey.UserID).(uuid.UUID)
	member, err := h.circleService.AddMember(c.Request.Context(), circle.AddMemberCommand{
		CircleID: circleID,
		UserID:   userID,
		MemberID: memberID,
		Role:     circleDomain.MemberRole(req.Role),
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, member)
}

// DELETE /circles/:id/members/:memberId
func (h *CircleHandler) RemoveMember(c *gin.Context) {
	circleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid circle ID")
		return
	}

	memberID, err := uuid.Parse(c.Param("memberId"))
	if err != nil {
		response.BadRequest(c, "invalid member ID")
		return
	}

	userID := c.MustGet(ctxkey.UserID).(uuid.UUID)
	err = h.circleService.RemoveMember(c.Request.Context(), circle.RemoveMemberCommand{
		CircleID: circleID,
		UserID:   userID,
		MemberID: memberID,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.NoContent(c)
}

// POST /circles/:id/leave
func (h *CircleHandler) LeaveCircle(c *gin.Context) {
	circleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid circle ID")
		return
	}

	userID := c.MustGet(ctxkey.UserID).(uuid.UUID)
	err = h.circleService.LeaveCircle(c.Request.Context(), circle.LeaveCircleCommand{
		CircleID: circleID,
		UserID:   userID,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.NoContent(c)
}

// GET /circles/:id/sharing
func (h *CircleHandler) GetSharingPreference(c *gin.Context) {
	circleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid circle ID")
		return
	}

	userID := c.MustGet(ctxkey.UserID).(uuid.UUID)
	pref, err := h.circleService.GetSharingPreference(c.Request.Context(), circle.GetSharingPreferenceQuery{
		CircleID: circleID,
		UserID:   userID,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, pref)
}

type UpdateSharingPreferenceRequest struct {
	PrivacyLevel      *string `json:"privacy_level"`
	ShareTimezone     *bool   `json:"share_timezone"`
	ShareAvailability *bool   `json:"share_availability"`
	ShareLocation     *bool   `json:"share_location"`
	LocationPrecision *string `json:"location_precision"`
	ShareActivity     *bool   `json:"share_activity"`
}

// PUT /circles/:id/sharing
func (h *CircleHandler) UpdateSharingPreference(c *gin.Context) {
	circleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid circle ID")
		return
	}

	var req UpdateSharingPreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	userID := c.MustGet(ctxkey.UserID).(uuid.UUID)

	cmd := circle.UpdateSharingPreferenceCommand{
		CircleID:          circleID,
		UserID:            userID,
		ShareTimezone:     req.ShareTimezone,
		ShareAvailability: req.ShareAvailability,
		ShareLocation:     req.ShareLocation,
		ShareActivity:     req.ShareActivity,
	}

	if req.PrivacyLevel != nil {
		level := user.PrivacyLevel(*req.PrivacyLevel)
		cmd.PrivacyLevel = &level
	}

	if req.LocationPrecision != nil {
		precision := circleDomain.LocationPrecision(*req.LocationPrecision)
		cmd.LocationPrecision = &precision
	}

	pref, err := h.circleService.UpdateSharingPreference(c.Request.Context(), cmd)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, pref)
}

type CreateInvitationRequest struct {
	Type      string  `json:"type" validate:"required,oneof=direct link"`
	InviteeID *string `json:"invitee_id" validate:"omitempty,uuid"`
	MaxUses   *int    `json:"max_uses" validate:"omitempty,min=1"`
	ExpiresIn *int    `json:"expires_in_hours" validate:"omitempty,min=1"`
}

// POST /circles/:id/invitations
func (h *CircleHandler) CreateInvitation(c *gin.Context) {
	circleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid circle ID")
		return
	}

	var req CreateInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, formatValidationErrors(err))
		return
	}

	userID := c.MustGet(ctxkey.UserID).(uuid.UUID)

	cmd := circle.CreateInvitationCommand{
		CircleID:  circleID,
		InviterID: userID,
		Type:      circleDomain.InvitationType(req.Type),
		MaxUses:   req.MaxUses,
	}

	if req.InviteeID != nil {
		inviteeID, _ := uuid.Parse(*req.InviteeID)
		cmd.InviteeID = &inviteeID
	}

	if req.ExpiresIn != nil {
		duration := time.Duration(*req.ExpiresIn) * time.Hour
		cmd.ExpiresIn = &duration
	}

	inv, err := h.circleService.CreateInvitation(c.Request.Context(), cmd)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, inv)
}

// POST /circles/join
func (h *CircleHandler) AcceptInvitation(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.BadRequest(c, "invitation code is required")
		return
	}

	userID := c.MustGet(ctxkey.UserID).(uuid.UUID)
	circ, err := h.circleService.AcceptInvitation(c.Request.Context(), circle.AcceptInvitationCommand{
		Code:   code,
		UserID: userID,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, circ)
}
