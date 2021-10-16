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
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetAllOrganisations(w http.ResponseWriter, r *http.Request) {
	if request.Authorize(w, r) {
		organisations := []model.Organisation{}
		db.DB.Find(&organisations)
		request.Respond(w, http.StatusOK, &organisations)
	}
}

func GetSingleOrganisation(w http.ResponseWriter, r *http.Request) {
	if request.Authorize(w, r) {
		vars := mux.Vars(r)
		organisationId, _ := uuid.Parse(vars["id"])

		organisation := model.Organisation{}

		if err := db.DB.First(&organisation, "Id = ?", organisationId); err.Error != nil {
			sentry.CaptureException(err.Error)
			request.Respond(w, http.StatusInternalServerError, "Organisation not found")
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
			sentry.CaptureException(err.Error)
			fmt.Println(err.Error.Error())
			request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("No organisations found for User ID '%s'", userId))
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
			request.Respond(w, http.StatusBadRequest, fmt.Sprintf("Invalid Body: %s", err.Error()))
		} else {
			defer r.Body.Close()

			currentUserId := authentication.GetCurrentUser(w, r)["Id"].(string)

			organisation.OwnerID = currentUserId

			if err := db.DB.Create(&organisation); err.Error != nil {
				sentry.CaptureException(err.Error)
				request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Error whilst creating new Organisation: %s", err.Error.Error()))
			} else {
				request.Respond(w, http.StatusOK, organisation)
			}
		}
	}
}

func UpdateOrganisation(w http.ResponseWriter, r *http.Request) {}

func DeleteOrganisation(w http.ResponseWriter, r *http.Request) {}
