package index

import (
	"fmt"
	"net/http"

	"github.com/benhall-1/appealscc/api/internal/request"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	if request.Authorize(w, r) {
		request.Respond(w, http.StatusOK, "🎉 Success! Welcome to the AppealsCC API")
	}
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	request.Respond(w, http.StatusNotFound, fmt.Sprintf("😢 The path '%s' is not found", r.URL.Path))
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	request.Respond(w, http.StatusMethodNotAllowed, fmt.Sprintf("🚫 The method '%s' is not allowed for '%s'", r.Method, r.URL.Path))
}
