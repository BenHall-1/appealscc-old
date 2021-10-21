package routing

import (
	"net/http"

	"github.com/benhall-1/appealscc/api/routing/endpoints/appeals"
	"github.com/benhall-1/appealscc/api/routing/endpoints/appeals/templates"
	"github.com/benhall-1/appealscc/api/routing/endpoints/auth"
	"github.com/benhall-1/appealscc/api/routing/endpoints/index"
	"github.com/benhall-1/appealscc/api/routing/endpoints/organisations"

	"github.com/gorilla/mux"
)

func SetupRequests(router *mux.Router) {

	router.Use(commonMiddleware)

	// Define default API Routes
	router.HandleFunc("/", index.HomePage)

	// Define Appeals API Routes
	router.HandleFunc("/api/appeals/{organisationId}", appeals.GetAllAppealsForOrganisation).Methods("GET")
	router.HandleFunc("/api/appeals/{organisationId}/{appealId}", appeals.GetSingleAppeal).Methods("GET")
	router.HandleFunc("/api/appeals/{organisationId}/create", appeals.CreateAppeal).Methods("POST")
	router.HandleFunc("/api/appeals/{organisationId}/{appealId}/respond", appeals.AddAppealResponse).Methods("POST")
	router.HandleFunc("/api/appeals/{organisationId}/templates", templates.GetAllTemplates).Methods("GET")
	router.HandleFunc("/api/appeals/{organisationId}/templates/{templateId}", templates.GetTemplateById).Methods("GET")
	router.HandleFunc("/api/appeals/{organisationId}/templates/create", templates.CreateTemplate).Methods("POST")
	router.HandleFunc("/api/appeals/{organisationId}/templates/{templateId}/update", templates.UpdateTemplate).Methods("PUT")
	router.HandleFunc("/api/appeals/{organisationId}/templates/{templateId}/delete", templates.DeleteTemplate).Methods("DELETE")

	// Define Authentication API Routes
	router.HandleFunc("/api/auth/register", auth.Register).Methods("POST")
	router.HandleFunc("/api/auth/refresh", auth.RefreshToken).Methods("POST")
	router.HandleFunc("/api/auth/login", auth.Login).Methods("POST")
	router.HandleFunc("/api/auth/discord", auth.LoginWithDiscord).Methods("GET")
	router.HandleFunc("/api/auth/callback", auth.AuthCallback).Methods("GET")

	// Define Organisations API Routes
	router.HandleFunc("/api/organisations/create", organisations.CreateOrganisation).Methods("POST")
	router.HandleFunc("/api/organisations/{id}/update", organisations.UpdateOrganisation).Methods("PUT")
	router.HandleFunc("/api/organisations/{id}/delete", organisations.DeleteOrganisation).Methods("DELETE")
	router.HandleFunc("/api/organisations", organisations.GetAllOrganisations).Methods("GET")
	router.HandleFunc("/api/organisations/byuser/{userId}", organisations.GetAllOrganisationsForUser).Methods("GET")
	router.HandleFunc("/api/organisations/{id}", organisations.GetSingleOrganisation).Methods("GET")
	router.HandleFunc("/api/organisations/{id}/moderators/{userId}/add", organisations.AddOrganisationModerator).Methods("POST")
	router.HandleFunc("/api/organisations/{id}/moderators/{userId}/remove", organisations.RemoveOrganisationModerator).Methods("DELETE")

	// Handling Errors
	router.NotFoundHandler = http.HandlerFunc(index.NotFound)
	router.MethodNotAllowedHandler = http.HandlerFunc(index.MethodNotAllowed)
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
