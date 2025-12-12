package converter

import (
	kinv1 "github.com/danielng/kin-core-svc/gen/proto/kin/v1"
	"github.com/danielng/kin-core-svc/internal/domain/circle"
	"github.com/danielng/kin-core-svc/internal/domain/user"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CircleToProto(c *circle.Circle) *kinv1.Circle {
	if c == nil {
		return nil
	}

	pb := &kinv1.Circle{
		Id:        c.ID.String(),
		Name:      c.Name,
		CreatedBy: c.CreatedBy.String(),
		CreatedAt: timestamppb.New(c.CreatedAt),
		UpdatedAt: timestamppb.New(c.UpdatedAt),
	}

	if c.Description != nil {
		pb.Description = c.Description
	}
	if c.Avatar != nil {
		pb.Avatar = c.Avatar
	}

	return pb
}

func CirclesToProto(circles []*circle.Circle) []*kinv1.Circle {
	result := make([]*kinv1.Circle, len(circles))
	for i, c := range circles {
		result[i] = CircleToProto(c)
	}
	return result
}

func MemberToProto(m *circle.Member) *kinv1.Member {
	if m == nil {
		return nil
	}

	pb := &kinv1.Member{
		Id:        m.ID.String(),
		CircleId:  m.CircleID.String(),
		UserId:    m.UserID.String(),
		Role:      MemberRoleToProto(m.Role),
		JoinedAt:  timestamppb.New(m.JoinedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
	}

	if m.Nickname != nil {
		pb.Nickname = m.Nickname
	}

	return pb
}

func MembersToProto(members []*circle.Member) []*kinv1.Member {
	result := make([]*kinv1.Member, len(members))
	for i, m := range members {
		result[i] = MemberToProto(m)
	}
	return result
}

func MemberRoleToProto(role circle.MemberRole) kinv1.MemberRole {
	switch role {
	case circle.MemberRoleAdmin:
		return kinv1.MemberRole_MEMBER_ROLE_ADMIN
	case circle.MemberRoleMember:
		return kinv1.MemberRole_MEMBER_ROLE_MEMBER
	default:
		return kinv1.MemberRole_MEMBER_ROLE_UNSPECIFIED
	}
}

func MemberRoleFromProto(role kinv1.MemberRole) circle.MemberRole {
	switch role {
	case kinv1.MemberRole_MEMBER_ROLE_ADMIN:
		return circle.MemberRoleAdmin
	case kinv1.MemberRole_MEMBER_ROLE_MEMBER:
		return circle.MemberRoleMember
	default:
		return circle.MemberRoleMember
	}
}

func InvitationToProto(i *circle.Invitation) *kinv1.Invitation {
	if i == nil {
		return nil
	}

	pb := &kinv1.Invitation{
		Id:        i.ID.String(),
		CircleId:  i.CircleID.String(),
		InviterId: i.InviterID.String(),
		Type:      InvitationTypeToProto(i.Type),
		Code:      i.Code,
		Status:    InvitationStatusToProto(i.Status),
		UseCount:  int32(i.UseCount),
		CreatedAt: timestamppb.New(i.CreatedAt),
		UpdatedAt: timestamppb.New(i.UpdatedAt),
	}

	if i.InviteeID != nil {
		inviteeID := i.InviteeID.String()
		pb.InviteeId = &inviteeID
	}
	if i.MaxUses != nil {
		maxUses := int32(*i.MaxUses)
		pb.MaxUses = &maxUses
	}
	if i.ExpiresAt != nil {
		pb.ExpiresAt = timestamppb.New(*i.ExpiresAt)
	}

	return pb
}

func InvitationTypeToProto(t circle.InvitationType) kinv1.InvitationType {
	switch t {
	case circle.InvitationTypeDirect:
		return kinv1.InvitationType_INVITATION_TYPE_DIRECT
	case circle.InvitationTypeLink:
		return kinv1.InvitationType_INVITATION_TYPE_LINK
	default:
		return kinv1.InvitationType_INVITATION_TYPE_UNSPECIFIED
	}
}

func InvitationTypeFromProto(t kinv1.InvitationType) circle.InvitationType {
	switch t {
	case kinv1.InvitationType_INVITATION_TYPE_DIRECT:
		return circle.InvitationTypeDirect
	case kinv1.InvitationType_INVITATION_TYPE_LINK:
		return circle.InvitationTypeLink
	default:
		return circle.InvitationTypeDirect
	}
}

func InvitationStatusToProto(s circle.InvitationStatus) kinv1.InvitationStatus {
	switch s {
	case circle.InvitationStatusPending:
		return kinv1.InvitationStatus_INVITATION_STATUS_PENDING
	case circle.InvitationStatusAccepted:
		return kinv1.InvitationStatus_INVITATION_STATUS_ACCEPTED
	case circle.InvitationStatusExpired:
		return kinv1.InvitationStatus_INVITATION_STATUS_EXPIRED
	case circle.InvitationStatusRevoked:
		return kinv1.InvitationStatus_INVITATION_STATUS_REVOKED
	default:
		return kinv1.InvitationStatus_INVITATION_STATUS_UNSPECIFIED
	}
}

func SharingPreferenceToProto(sp *circle.SharingPreference) *kinv1.SharingPreference {
	if sp == nil {
		return nil
	}

	return &kinv1.SharingPreference{
		Id:                sp.ID.String(),
		CircleId:          sp.CircleID.String(),
		UserId:            sp.UserID.String(),
		PrivacyLevel:      PrivacyLevelToProto(sp.PrivacyLevel),
		ShareTimezone:     sp.ShareTimezone,
		ShareAvailability: sp.ShareAvailability,
		ShareLocation:     sp.ShareLocation,
		LocationPrecision: LocationPrecisionToProto(sp.LocationPrecision),
		ShareActivity:     sp.ShareActivity,
		CreatedAt:         timestamppb.New(sp.CreatedAt),
		UpdatedAt:         timestamppb.New(sp.UpdatedAt),
	}
}

func LocationPrecisionToProto(p circle.LocationPrecision) kinv1.LocationPrecision {
	switch p {
	case circle.LocationPrecisionCountry:
		return kinv1.LocationPrecision_LOCATION_PRECISION_COUNTRY
	case circle.LocationPrecisionCity:
		return kinv1.LocationPrecision_LOCATION_PRECISION_CITY
	case circle.LocationPrecisionNeighborhood:
		return kinv1.LocationPrecision_LOCATION_PRECISION_NEIGHBORHOOD
	case circle.LocationPrecisionExact:
		return kinv1.LocationPrecision_LOCATION_PRECISION_EXACT
	default:
		return kinv1.LocationPrecision_LOCATION_PRECISION_UNSPECIFIED
	}
}

func LocationPrecisionFromProto(p kinv1.LocationPrecision) circle.LocationPrecision {
	switch p {
	case kinv1.LocationPrecision_LOCATION_PRECISION_COUNTRY:
		return circle.LocationPrecisionCountry
	case kinv1.LocationPrecision_LOCATION_PRECISION_CITY:
		return circle.LocationPrecisionCity
	case kinv1.LocationPrecision_LOCATION_PRECISION_NEIGHBORHOOD:
		return circle.LocationPrecisionNeighborhood
	case kinv1.LocationPrecision_LOCATION_PRECISION_EXACT:
		return circle.LocationPrecisionExact
	default:
		return circle.LocationPrecisionCity
	}
}

func PrivacyLevelToProtoFromCircle(level user.PrivacyLevel) kinv1.PrivacyLevel {
	return PrivacyLevelToProto(level)
}
