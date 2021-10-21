package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/benhall-1/appealscc/api/internal/authentication"
	"github.com/benhall-1/appealscc/api/internal/db"
	"github.com/benhall-1/appealscc/api/internal/models/authmodel"
	"github.com/benhall-1/appealscc/api/internal/models/discordmodel"
	"github.com/benhall-1/appealscc/api/internal/models/model"
	"github.com/benhall-1/appealscc/api/internal/oauth"
	"github.com/benhall-1/appealscc/api/internal/request"
	"github.com/getsentry/sentry-go"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user model.User
	var requestUser authmodel.LoginRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestUser); err != nil {
		sentryError := sentry.CaptureException(err)
		request.Respond(w, http.StatusBadRequest, fmt.Sprintf("😢 Request failed - Please try again. Error code: '%s'", *sentryError))
	} else {
		defer r.Body.Close()

		user.Email = requestUser.Email
		user.Password = requestUser.Password

		if status, _ := authentication.RegisterAccount(&user, nil); status {
			request.Respond(w, http.StatusOK, "Account Registered")
		} else {
			request.Respond(w, http.StatusBadRequest, "😢 Request failed - Please try again")
		}
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest authmodel.LoginRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginRequest); err != nil {
		sentryError := sentry.CaptureException(err)
		request.Respond(w, http.StatusBadRequest, fmt.Sprintf("😢 Request failed - Please try again. Error code: '%s'", *sentryError))
	} else {
		defer r.Body.Close()

		user := model.User{}
		if err := db.DB.First(&user, "Email = ?", loginRequest.Email); err.Error != nil {
			sentry.CaptureException(err.Error)
			request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("😢 %s", err.Error))
		} else {
			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
				sentryError := sentry.CaptureException(err)
				request.Respond(w, http.StatusUnauthorized, fmt.Sprintf("🚫 Incorrect username or password. Error code: %s", *sentryError))
			} else {
				tokenResponse, err := authentication.GenerateToken(user)
				if err != nil {
					sentryError := sentry.CaptureException(err)
					request.Respond(w, http.StatusUnauthorized, fmt.Sprintf("🚫 Incorrect username or password. Error code: %s", *sentryError))
				} else {
					request.Respond(w, http.StatusOK, tokenResponse)
				}
			}
		}
	}
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	request.RefreshToken(w, r)
}

var state = "random"

func LoginWithDiscord(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, oauth.DiscordOAuth().AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func AuthCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != state {
		w.WriteHeader(http.StatusBadRequest)
		request.Respond(w, http.StatusBadRequest, "State does not match.")
		return
	}
	// Step 3: We exchange the code we got for an access token
	// Then we can use the access token to do actions, limited to scopes we requested
	token, err := oauth.DiscordOAuth().Exchange(context.Background(), r.FormValue("code"))

	if err != nil {
		sentryError := sentry.CaptureException(err)
		request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Error whilst contacting discord. Error code '%s'", *sentryError))
		return
	}

	// Step 4: Use the access token, here we use it to get the logged in user's info.
	res, err := oauth.DiscordOAuth().Client(context.Background(), token).Get("https://discord.com/api/users/@me")

	if err != nil || res.StatusCode != 200 {
		if err != nil {
			sentryError := sentry.CaptureException(err)
			request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Error whilst fetching your details from discord. Error code '%s'", *sentryError))
		} else {
			request.Respond(w, http.StatusInternalServerError, res.Status)
		}
		return
	}

	defer res.Body.Close()

	var discordUser discordmodel.DiscordUser
	decoder := json.NewDecoder(res.Body)
	decoder.Decode(&discordUser)

	if err != nil {
		sentryError := sentry.CaptureException(err)
		request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Error whilst fetching your details from discord. Error code '%s'", *sentryError))
		return
	}

	if status, user := authentication.RegisterAccount(nil, &discordUser); status {
		token, err := authentication.GenerateToken(*user)

		if err != nil {
			sentryError := sentry.CaptureException(err)
			request.Respond(w, http.StatusInternalServerError, fmt.Sprintf("Error whilst fetching your details. Error code '%s'", *sentryError))
			return
		} else {
			request.Respond(w, http.StatusOK, token)
		}
	} else {
		request.Respond(w, http.StatusBadRequest, "😢 Request failed - Please try again")
	}
}
