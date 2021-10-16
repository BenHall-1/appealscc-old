package authentication

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/benhall-1/appealscc/api/internal/db"
	"github.com/benhall-1/appealscc/api/internal/models/discordmodel"
	"github.com/benhall-1/appealscc/api/internal/models/model"
	"github.com/getsentry/sentry-go"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"github.com/sethvargo/go-password/password"
)

var jwtKey = []byte(os.Getenv("SECRET"))

func RegisterAccount(user *model.User, discord *discordmodel.DiscordUser) (bool, *model.User) {
	var userExists *model.User
	if user == nil {
		user = &model.User{}
	}

	if len(user.Email) == 0 && discord != nil {
		fmt.Println(user)
		user.Email = discord.Email
	} else {
		return false, nil
	}

	if len(user.Password) == 0 {
		pass, _ := password.Generate(64, 10, 10, false, false)
		user.Password = pass
	}

	if result := db.DB.First(&userExists, "Email = ?", user.Email); result.RowsAffected > 0 {
		if discord != nil {
			return true, userExists
		} else {
			return false, nil
		}
	}

	plainPassword := user.Password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), 14)
	user.Password = string(hashedPassword)

	err := db.DB.Create(&user)
	if err.Error != nil {
		sentry.CaptureException(err.Error)
		return false, nil
	} else {
		return true, user
	}
}

func GenerateToken(user model.User) *model.TokenResponse {
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &model.Claims{
		Email: user.Email,
		Id:    user.ID.String(),
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    os.Getenv("TOKEN_ISSUER"),
			Audience:  os.Getenv("BASE_URL"),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		sentry.CaptureException(err)
		return nil
	} else {
		return &model.TokenResponse{Token: tokenString, Expiration: expirationTime}
	}
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) jwt.MapClaims {
	token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	tkn, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	claims, _ := tkn.Claims.(jwt.MapClaims)

	return claims
}
