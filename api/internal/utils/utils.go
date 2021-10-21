package utils

import (
	"github.com/benhall-1/appealscc/api/internal/db"
	"github.com/benhall-1/appealscc/api/internal/models/model"
	"github.com/getsentry/sentry-go"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func IsOrganisationOwnerOrGlobalAdmin(orgId uuid.UUID, user jwt.MapClaims) bool {
	var organisation *model.Organisation
	if err := db.DB.First(&organisation, "Id = ?", orgId); err.Error != nil {
		sentry.CaptureException(err.Error)
		return false
	} else {
		currentUserId, _ := uuid.Parse(user["Id"].(string))
		if organisation.OwnerID == currentUserId {
			return true
		} else {
			return user["GlobalAdmin"].(bool)
		}
	}
}
