package appeals

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/benhall-1/appealscc/api/internal/authentication"
	"github.com/benhall-1/appealscc/api/internal/db"
	"github.com/benhall-1/appealscc/api/internal/models/model"
	"github.com/benhall-1/appealscc/api/internal/request"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetAllAppealsForOrganisation(w http.ResponseWriter, r *http.Request) {
	if request.Authorize(w, r) {
		vars := mux.Vars(r)
		organisationId, _ := uuid.Parse(vars["organisationId"])

		appeals := []model.Appeal{}

		if err := db.DB.Find(&appeals, "Organisation = ?", organisationId); err.Error != nil {
			sentryError := sentry.CaptureException(err.Error)
			request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Error getting Appeals. Error code '%s'", *sentryError))
		} else {
			request.Respond(w, http.StatusOK, appeals)
		}
	}
}

func GetSingleAppeal(w http.ResponseWriter, r *http.Request) {
	if request.Authorize(w, r) {
		vars := mux.Vars(r)
		organisationId, _ := uuid.Parse(vars["organisationId"])
		appealId, _ := uuid.Parse(vars["appealId"])

		appeal := model.Appeal{}

		if err := db.DB.First(&appeal, "Id = ? AND Organisation = ?", appealId, organisationId); err.Error != nil {
			sentryError := sentry.CaptureException(err.Error)
			request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Appeal not found. Error code '%s'", *sentryError))
		} else {
			request.Respond(w, http.StatusOK, appeal)
		}
	}
}

func CreateAppeal(w http.ResponseWriter, r *http.Request) {
	if request.Authorize(w, r) {
		var appeal model.Appeal
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&appeal); err != nil {
			sentryError := sentry.CaptureException(err)
			request.Respond(w, http.StatusBadRequest, fmt.Sprintf("Invalid body in request. Error code '%s'", *sentryError))
		} else {
			defer r.Body.Close()
			vars := mux.Vars(r)

			organisationId, _ := uuid.Parse(vars["organisationId"])
			currentUserId, _ := uuid.Parse(authentication.GetCurrentUser(w, r)["Id"].(string))

			var tempOrg model.Organisation
			var tempAppealTemplate model.AppealTemplate
			var tempAppeal *model.Appeal

			if err := db.DB.Find(&tempOrg, "Id = ?", &organisationId); err.Error != nil {
				sentryError := sentry.CaptureException(err.Error)
				request.Respond(w, http.StatusBadRequest, fmt.Sprintf("Organisation not found. Error code '%s'", *sentryError))
			} else {
				if err := db.DB.Find(&tempAppealTemplate, "Id = ? AND organisation = ? ", appeal.Template, organisationId); err.Error != nil {
					sentryError := sentry.CaptureException(err.Error)
					request.Respond(w, http.StatusBadRequest, fmt.Sprintf("Template not found. Error code '%s'", *sentryError))
				} else {
					if err := db.DB.First(&tempAppeal, "creator = ? AND template = ? AND appeal_status = 0", appeal.Template, organisationId); err.Error != nil && tempAppeal != nil {
						sentryError := sentry.CaptureException(err.Error)
						request.Respond(w, http.StatusBadRequest, fmt.Sprintf("Appeal creation failed - You already have an open appeal for this form. Error code '%s'", *sentryError))
					} else {
						appeal.Creator = currentUserId
						appeal.Organisation = organisationId
						if err := db.DB.Create(&appeal); err.Error != nil {
							sentryError := sentry.CaptureException(err.Error)
							request.Respond(w, http.StatusBadRequest, fmt.Sprintf("Appeal creation failed. Error code '%s'", *sentryError))
						} else {
							request.Respond(w, http.StatusOK, appeal)
						}
					}
				}
			}
		}
	}
}

func AddAppealResponse(w http.ResponseWriter, r *http.Request) {
	if request.Authorize(w, r) {
		vars := mux.Vars(r)
		appealId, _ := uuid.Parse(vars["appealId"])

		var appealResponse model.AppealResponse
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&appealResponse); err != nil {
			sentryError := sentry.CaptureException(err)
			request.Respond(w, http.StatusBadRequest, fmt.Sprintf("Invalid body in request. Error code '%s'", *sentryError))
		} else {
			defer r.Body.Close()

			currentUserId := authentication.GetCurrentUser(w, r)["Id"].(string)

			appealResponse.Author, _ = uuid.Parse(currentUserId)
			appealResponse.Appeal = appealId

			if err := db.DB.Create(&appealResponse); err.Error != nil {
				sentryError := sentry.CaptureException(err.Error)
				request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Error whilst creating new Appeal Response. Error code '%s'", *sentryError))
			} else {
				if appealResponse.Decision != 0 {
					var appeal model.Appeal
					if err := db.DB.First(&appeal, "Id = ?", appealId); err.Error != nil {
						appeal.AppealStatus = appealResponse.Decision
						if err := db.DB.Save(&appeal); err.Error != nil {
							sentryError := sentry.CaptureException(err.Error)
							request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Error whilst creating new Appeal Response. Error code '%s'", *sentryError))
						} else {
							request.Respond(w, http.StatusOK, "Appeal response created & decision sent")
						}
					}
				}
				request.Respond(w, http.StatusOK, appealResponse)
			}
		}
	}
}
