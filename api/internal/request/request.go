package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/benhall-1/appealscc/api/internal/models/model"
	"github.com/getsentry/sentry-go"
	"github.com/golang-jwt/jwt"
)

// Create the JWT string
var jwtKey = []byte(os.Getenv("SECRET"))

func createResponse(status int, body interface{}) model.Response {
	return model.Response{Status: status, Body: body}
}

func Respond(w http.ResponseWriter, status int, body interface{}, formats ...string) error {
	w.WriteHeader(status)
	if fmt.Sprintf("%T", body) == "string" {
		return json.NewEncoder(w).Encode(createResponse(status, fmt.Sprintf(body.(string), formats)))
	} else {
		return json.NewEncoder(w).Encode(createResponse(status, body))
	}
}

func Authorize(w http.ResponseWriter, r *http.Request) bool {
	token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
	if token == "" {
		Respond(w, http.StatusUnauthorized, "Access Denied - Token not found")
		return false
	}
	claims := &model.Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		sentryError := sentry.CaptureException(err)
		if err == jwt.ErrSignatureInvalid {
			Respond(w, http.StatusUnauthorized, "Invalid signature for provided token. Error code '%s'", string(*sentryError))
			return false
		}
		Respond(w, http.StatusBadRequest, "Error whilst processing token. Error code '%s'", string(*sentryError))
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
	claims := &model.Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		sentryError := sentry.CaptureException(err)
		if err == jwt.ErrSignatureInvalid {
			Respond(w, http.StatusUnauthorized, "Invalid signature for provided token. Error code '%s'", string(*sentryError))
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

	Respond(w, http.StatusOK, model.TokenResponse{Token: tokenString, Expiration: expirationTime})

	return true
}
