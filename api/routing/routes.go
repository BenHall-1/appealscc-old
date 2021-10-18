package routing

import (
	"net/http"

	"github.com/benhall-1/appealscc/api/routing/endpoints/appeals"
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
	router.HandleFunc("/api/appeals", appeals.ReturnAllAppeals).Methods("GET")
	router.HandleFunc("/api/appeals/{id}", appeals.ReturnSingleAppeal).Methods("GET")
	router.HandleFunc("/api/appeals", appeals.CreateAppeal).Methods("POST")

	// Define Authentication API Routes
	router.HandleFunc("/api/auth/register", auth.Register).Methods("POST")
	router.HandleFunc("/api/auth/refresh", auth.RefreshToken).Methods("POST")
	router.HandleFunc("/api/auth/login", auth.Login).Methods("POST")
	router.HandleFunc("/api/auth/discord", auth.LoginWithDiscord).Methods("GET")
	router.HandleFunc("/api/auth/callback", auth.AuthCallback).Methods("GET")

	// Define Organisations API Routes
	router.HandleFunc("/api/organisations", organisations.CreateOrganisation).Methods("POST")
	router.HandleFunc("/api/organisations/{id}", organisations.UpdateOrganisation).Methods("PUT")
	router.HandleFunc("/api/organisations/{id}", organisations.DeleteOrganisation).Methods("DELETE")
	router.HandleFunc("/api/organisations", organisations.GetAllOrganisations).Methods("GET")
	router.HandleFunc("/api/organisations/byuser/{userId}", organisations.GetAllOrganisationsForUser).Methods("GET")
	router.HandleFunc("/api/organisations/{id}", organisations.GetSingleOrganisation).Methods("GET")
	router.HandleFunc("/api/organisations/{id}/moderators/{userId}", organisations.AddOrganisationModerator).Methods("POST")
	router.HandleFunc("/api/organisations/{id}/moderators/{userId}", organisations.RemoveOrganisationModerator).Methods("DELETE")

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
