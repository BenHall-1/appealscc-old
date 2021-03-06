package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/benhall-1/appealscc/api/internal/models/authmodel"
	"github.com/getsentry/sentry-go"
	"github.com/golang-jwt/jwt"
)

type Response struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

// Create the JWT string
var jwtKey = []byte(os.Getenv("SECRET"))

func createResponse(status int, body interface{}) Response {
	return Response{Status: status, Body: body}
}

func Respond(w http.ResponseWriter, status int, body interface{}) error {
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(createResponse(status, body))
}

func Authorize(w http.ResponseWriter, r *http.Request) bool {
	token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
	if token == "" {
		Respond(w, http.StatusUnauthorized, "Access Denied - Token not found")
		return false
	}
	claims := &authmodel.Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		sentryError := sentry.CaptureException(err)
		if err == jwt.ErrSignatureInvalid {
			Respond(w, http.StatusUnauthorized, fmt.Sprintf("Invalid signature for provided token. Error code '%s'", *sentryError))
			return false
		}
		Respond(w, http.StatusBadRequest, fmt.Sprintf("Error whilst processing token. Error code '%s'", *sentryError))
		return false
	}
	if !tkn.Valid {
		Respond(w, http.StatusUnauthorized, "Expired token - Please try logging in again")
		return false
	}

	return true
}

func RefreshToken(w http.ResponseWriter, r *http.Request) bool {
	token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
	if token == "" {
		Respond(w, http.StatusUnauthorized, "Access Denied - Token not found")
		return false
	}
	claims := &authmodel.Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		sentryError := sentry.CaptureException(err)
		if err == jwt.ErrSignatureInvalid {
			Respond(w, http.StatusUnauthorized, fmt.Sprintf("Invalid signature for provided token. Error code '%s'", *sentryError))
			return false
		}
		Respond(w, http.StatusBadRequest, "Error whilst authenticating")
		return false
	}
	if !tkn.Valid {
		Respond(w, http.StatusUnauthorized, "Expired token - Please try logging in again")
		return false
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	nToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := nToken.SignedString(jwtKey)
	if err != nil {
		sentry.CaptureException(err)
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}

	Respond(w, http.StatusOK, authmodel.TokenResponse{Token: tokenString, Expiration: expirationTime})

	return true
}
