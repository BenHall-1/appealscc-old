package templates

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/benhall-1/appealscc/api/internal/authentication"
	"github.com/benhall-1/appealscc/api/internal/db"
	"github.com/benhall-1/appealscc/api/internal/models/model"
	"github.com/benhall-1/appealscc/api/internal/request"
	"github.com/benhall-1/appealscc/api/internal/utils"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetAllTemplates(w http.ResponseWriter, r *http.Request) {
	if request.Authorize(w, r) {
		vars := mux.Vars(r)
		organisationId, _ := uuid.Parse(vars["organisationId"])

		var templates []model.AppealTemplate

		if err := db.DB.Find(&templates, "organisation = ?", organisationId); err.Error != nil {
			sentryError := sentry.CaptureException(err.Error)
			request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Error whilst getting all Appeal Templates. Error code '%s'", *sentryError))
		} else {
			request.Respond(w, http.StatusOK, templates)
		}
	}
}

func GetTemplateById(w http.ResponseWriter, r *http.Request) {
	if request.Authorize(w, r) {
		vars := mux.Vars(r)
		organisationId, _ := uuid.Parse(vars["organisationId"])
		templateId, _ := uuid.Parse(vars["templateId"])

		var template model.AppealTemplate

		if err := db.DB.Preload("AppealTemplateFields").First(&template, "organisation = ? AND Id = ?", organisationId, templateId); err.Error != nil {
			sentryError := sentry.CaptureException(err.Error)
			request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Error whilst getting Appeal Template. Error code '%s'", *sentryError))
		} else {
			request.Respond(w, http.StatusOK, template)
		}
	}
}

func CreateTemplate(w http.ResponseWriter, r *http.Request) {
	if request.Authorize(w, r) {
		vars := mux.Vars(r)
		organisationId, _ := uuid.Parse(vars["organisationId"])

		currentUser := authentication.GetCurrentUser(w, r)

		if utils.IsOrganisationOwnerOrGlobalAdmin(organisationId, currentUser) {

			var tempOrg model.Organisation
			currentUserPremiumType := authentication.GetCurrentUser(w, r)["PremiumType"].(float64)

			if err := db.DB.Preload("AppealTemplates").First(&tempOrg, "Id = ?", organisationId); err.Error != nil {
				sentryError := sentry.CaptureException(err.Error)
				request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Error whilst creating a new Appeal Template. Error code '%s'", *sentryError))
			} else {
				if currentUserPremiumType == 0 && len(tempOrg.AppealTemplates) == 2 {
					request.Respond(w, http.StatusBadRequest, "Error whilst creating a new appeal template - You have reached the maximum number of appeal templates for the Free plan.")
				} else {
					var appealTemplate model.AppealTemplate
					decoder := json.NewDecoder(r.Body)
					if err := decoder.Decode(&appealTemplate); err != nil {
						sentryError := sentry.CaptureException(err)
						request.Respond(w, http.StatusBadRequest, fmt.Sprintf("Invalid body in request. Error code '%s'", *sentryError))
					} else {
						defer r.Body.Close()

						appealTemplate.Organisation = organisationId

						if err := db.DB.Create(&appealTemplate); err.Error != nil {
							sentryError := sentry.CaptureException(err.Error)
							request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Error whilst creating a new Appeal Template. Error code '%s'", *sentryError))
						} else {
							request.Respond(w, http.StatusOK, appealTemplate)
						}
					}
				}
			}
		} else {
			request.Respond(w, http.StatusForbidden, "Access Denied - You are not the owner of the organisation")
		}
	}
}

func UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	if request.Authorize(w, r) {
		vars := mux.Vars(r)
		organisationId, _ := uuid.Parse(vars["organisationId"])
		templateId, _ := uuid.Parse(vars["templateId"])

		currentUser := authentication.GetCurrentUser(w, r)

		if utils.IsOrganisationOwnerOrGlobalAdmin(organisationId, currentUser) {

			var appealTemplate model.AppealTemplate

			if err := db.DB.First(&appealTemplate, "Id = ?", templateId); err.Error != nil {
				sentryError := sentry.CaptureException(err.Error)
				request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Appeal template does not exist. Error code '%s'", *sentryError))
			} else {

				decoder := json.NewDecoder(r.Body)
				if err := decoder.Decode(&appealTemplate); err != nil {
					sentryError := sentry.CaptureException(err)
					request.Respond(w, http.StatusBadRequest, fmt.Sprintf("Invalid body in request. Error code '%s'", *sentryError))
				} else {
					defer r.Body.Close()

					appealTemplate.Organisation = organisationId

					if err := db.DB.Model(&appealTemplate).Omit("AppealTemplateFields.*").Save(&appealTemplate); err.Error != nil {
						sentryError := sentry.CaptureException(err.Error)
						request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Error whilst update the Appeal Template. Error code '%s'", *sentryError))
					} else {
						if err := db.DB.Model(&appealTemplate).Association("AppealTemplateFields").Replace(&appealTemplate.AppealTemplateFields); err != nil {
							sentryError := sentry.CaptureException(err)
							request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Error whilst update the Appeal Template. Error code '%s'", *sentryError))
						} else {
							if err := db.DB.Unscoped().Delete(model.AppealTemplateField{}, "template IS NULL"); err.Error != nil {
								sentryError := sentry.CaptureException(err.Error)
								request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Error whilst update the Appeal Template. Error code '%s'", *sentryError))
							} else {
								request.Respond(w, http.StatusOK, appealTemplate)
							}
						}
					}
				}
			}
		} else {
			request.Respond(w, http.StatusForbidden, "Access Denied - You are not the owner of the organisation")
		}
	}
}

func DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	if request.Authorize(w, r) {
		vars := mux.Vars(r)
		organisationId, _ := uuid.Parse(vars["organisationId"])
		templateId, _ := uuid.Parse(vars["templateId"])
		currentUser := authentication.GetCurrentUser(w, r)

		if utils.IsOrganisationOwnerOrGlobalAdmin(organisationId, currentUser) {
			template := model.AppealTemplate{}

			if err := db.DB.First(&template, "Id = ?", templateId); err.Error != nil {
				sentryError := sentry.CaptureException(err.Error)
				request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Template not found. Error code '%s'", *sentryError))
			} else {

				db.DB.Unscoped().Delete(&template)
				request.Respond(w, http.StatusOK, "Template deleted")
			}
		} else {
			request.Respond(w, http.StatusForbidden, "Access Denied - You are not the owner of the organisation")
		}
	}
}
