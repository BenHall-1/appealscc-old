package model

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type Response struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

type User struct {
	gorm.Model
	ID                 uuid.UUID       `json:"ID" gorm:"type:char(36);primary_key;"`
	Email              string          `json:"Email"`
	Password           string          `json:"Password"`
	OrganisationMods   []*Organisation `json:"Organisations" gorm:"many2many:organisation_moderators;"`
	OwnedOrganisations []*Organisation `json:"OwnedOrganisations" gorm:"foreignKey:OwnerID;"`
}

type Organisation struct {
	gorm.Model
	ID          uuid.UUID `json:"ID" gorm:"type:char(36);primary_key;"`
	Name        string    `json:"Name"`
	IconHash    *string   `json:"IconHash"`
	Description string    `json:"Description"`
	Moderators  []*User   `json:"Moderators" gorm:"many2many:organisation_moderators;"`
	OwnerID     string    `json:"Owner"`
}

type LoginRequest struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type Claims struct {
	Email string `json:"Email"`
	Id    string `json:"Id"`
	jwt.StandardClaims
}

type TokenResponse struct {
	Token      string    `json:"token"`
	Expiration time.Time `json:"expiration"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New()
	return
}

func (organisation *Organisation) BeforeCreate(tx *gorm.DB) (err error) {
	organisation.ID = uuid.New()
	return
}
