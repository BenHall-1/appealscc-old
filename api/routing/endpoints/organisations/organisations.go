package organisations

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/benhall-1/appealscc/api/internal/authentication"
	"github.com/benhall-1/appealscc/api/internal/db"
	"github.com/benhall-1/appealscc/api/internal/models/model"
	"github.com/benhall-1/appealscc/api/internal/request"
	"github.com/getsentry/sentry-go"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetAllOrganisations(w http.ResponseWriter, r *http.Request) {
	if request.Authorize(w, r) {
		organisations := []model.Organisation{}
		db.DB.Preload("Moderators").Find(&organisations)
		request.Respond(w, http.StatusOK, &organisations)
	}
}

func GetSingleOrganisation(w http.ResponseWriter, r *http.Request) {
	if request.Authorize(w, r) {
		vars := mux.Vars(r)
		organisationId, _ := uuid.Parse(vars["id"])

		organisation := model.Organisation{}

		if err := db.DB.First(&organisation, "Id = ?", organisationId); err.Error != nil {
			sentryError := sentry.CaptureException(err.Error)
			request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Organisation not found. Error code '%s'", *sentryError))
		} else {
			request.Respond(w, http.StatusOK, organisation)
		}
	}
}

func GetAllOrganisationsForUser(w http.ResponseWriter, r *http.Request) {
	if request.Authorize(w, r) {
		vars := mux.Vars(r)
		userId, _ := uuid.Parse(vars["userId"])

		organisation := []model.Organisation{}

		if err := db.DB.Find(&organisation, "owner_id = ?", userId); err.Error != nil {
			sentryError := sentry.CaptureException(err.Error)
			request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("No organisations found for User ID '%s'. Error code '%s'", userId, *sentryError))
		} else {
			request.Respond(w, http.StatusOK, organisation)
		}
	}
}

func CreateOrganisation(w http.ResponseWriter, r *http.Request) {
	if request.Authorize(w, r) {
		var organisation model.Organisation
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&organisation); err != nil {
			sentryError := sentry.CaptureException(err)
			request.Respond(w, http.StatusBadRequest, fmt.Sprintf("Invalid body in request. Error code '%s'", *sentryError))
		} else {
			defer r.Body.Close()

			currentUserId := authentication.GetCurrentUser(w, r)["Id"].(string)

			organisation.OwnerID, _ = uuid.Parse(currentUserId)

			if err := db.DB.Create(&organisation); err.Error != nil {
				sentryError := sentry.CaptureException(err.Error)
				request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Error whilst creating new Organisation. Error code '%s'", *sentryError))
			} else {
				request.Respond(w, http.StatusOK, organisation)
			}
		}
	}
}

func UpdateOrganisation(w http.ResponseWriter, r *http.Request) {
	if request.Authorize(w, r) {
		vars := mux.Vars(r)
		organisationId, _ := uuid.Parse(vars["id"])
		currentUser := authentication.GetCurrentUser(w, r)

		if isOrganisationOwnerOrGlobalAdmin(organisationId, currentUser) {
			organisation := model.Organisation{}

			if err := db.DB.First(&organisation, "Id = ?", organisationId); err.Error != nil {
				sentryError := sentry.CaptureException(err.Error)
				request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Organisation not found. Error code '%s'", *sentryError))
			} else {
				var bodyOrganisation model.Organisation
				decoder := json.NewDecoder(r.Body)
				if err := decoder.Decode(&bodyOrganisation); err != nil {
					sentryError := sentry.CaptureException(err)
					request.Respond(w, http.StatusBadRequest, fmt.Sprintf("Error whilst getting organisation. Error code '%s'", *sentryError))
				} else {
					defer r.Body.Close()

					if bodyOrganisation.Name != "" {
						organisation.Name = bodyOrganisation.Name
					}
					if bodyOrganisation.IconHash != nil {
						organisation.IconHash = bodyOrganisation.IconHash
					}
					if bodyOrganisation.Description != "" {
						organisation.Description = bodyOrganisation.Description
					}
					db.DB.Save(&organisation)
					request.Respond(w, http.StatusOK, organisation)
				}
			}
		} else {
			request.Respond(w, http.StatusForbidden, "Access Denied - You are not the owner of the organisation")
		}
	}
}

func DeleteOrganisation(w http.ResponseWriter, r *http.Request) {
	if request.Authorize(w, r) {
		vars := mux.Vars(r)
		organisationId, _ := uuid.Parse(vars["id"])
		currentUser := authentication.GetCurrentUser(w, r)

		if isOrganisationOwnerOrGlobalAdmin(organisationId, currentUser) {
			organisation := model.Organisation{}

			if err := db.DB.First(&organisation, "Id = ?", organisationId); err.Error != nil {
				sentryError := sentry.CaptureException(err.Error)
				request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Organisation not found. Error code '%s'", *sentryError))
			} else {
				db.DB.Delete(&organisation)
				request.Respond(w, http.StatusOK, "Organisation deleted")
			}
		} else {
			request.Respond(w, http.StatusForbidden, "Access Denied - You are not the owner of the organisation")
		}
	}
}

func AddOrganisationModerator(w http.ResponseWriter, r *http.Request) {
	if request.Authorize(w, r) {
		vars := mux.Vars(r)
		organisationId, _ := uuid.Parse(vars["id"])
		userId, _ := uuid.Parse(vars["userId"])
		currentUser := authentication.GetCurrentUser(w, r)

		if isOrganisationOwnerOrGlobalAdmin(organisationId, currentUser) {
			organisation := model.Organisation{}

			if err := db.DB.First(&organisation, "Id = ?", organisationId); err.Error != nil {
				sentryError := sentry.CaptureException(err.Error)
				request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Organisation not found. Error code '%s'", *sentryError))
			} else {
				newModerator := model.User{}
				if err := db.DB.First(&newModerator, "Id = ?", userId); err.Error != nil {
					sentryError := sentry.CaptureException(err.Error)
					request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("User not found. Error code '%s'", *sentryError))
				} else {
					fmt.Println(newModerator.ID.String())
					if err := db.DB.Model(&organisation).Omit("Moderators.*").Association("Moderators").Append(&newModerator); err != nil {
						sentryError := sentry.CaptureException(err)
						request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Could not add user as a Moderator for %s. Error code '%s'", organisation.Name, *sentryError))
					} else {
						request.Respond(w, http.StatusOK, fmt.Sprintf("User Id '%s' added to the Moderators list of organisation '%s'", userId, organisation.Name))
					}
				}
			}
		} else {
			request.Respond(w, http.StatusForbidden, "Access Denied - You are not the owner of the organisation")
		}
	}
}

func RemoveOrganisationModerator(w http.ResponseWriter, r *http.Request) {
	if request.Authorize(w, r) {
		vars := mux.Vars(r)
		organisationId, _ := uuid.Parse(vars["id"])
		userId, _ := uuid.Parse(vars["userId"])
		currentUser := authentication.GetCurrentUser(w, r)

		if isOrganisationOwnerOrGlobalAdmin(organisationId, currentUser) {
			organisation := model.Organisation{}

			if err := db.DB.First(&organisation, "Id = ?", organisationId); err.Error != nil {
				sentryError := sentry.CaptureException(err.Error)
				request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Organisation not found. Error code '%s'", *sentryError))
			} else {
				newModerator := model.User{}
				if err := db.DB.First(&newModerator, "Id = ?", userId); err.Error != nil {
					sentryError := sentry.CaptureException(err.Error)
					request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("User not found. Error code '%s'", *sentryError))
				} else {
					if err := db.DB.Model(&organisation).Association("Moderators").Delete(&newModerator); err != nil {
						sentryError := sentry.CaptureException(err)
						request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Could not remove user from being a Moderator. Error code '%s'", *sentryError))
					} else {
						request.Respond(w, http.StatusOK, fmt.Sprintf("User Id '%s' removed from Moderators list of organisation '%s'", userId, organisation.Name))
					}
				}
			}
		} else {
			request.Respond(w, http.StatusForbidden, "Access Denied - You are not the owner of the organisation")
		}
	}
}

func isOrganisationOwnerOrGlobalAdmin(orgId uuid.UUID, user jwt.MapClaims) bool {
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
