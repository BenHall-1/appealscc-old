package model

import (
	"encoding/json"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type User struct {
	Base
	Email              string           `json:"Email" gorm:"uniqueIndex;type:varchar(256);"`
	Password           string           `json:"-"`
	OrganisationMods   []*Organisation  `json:"Organisations" gorm:"many2many:organisation_moderators;"`
	OwnedOrganisations []*Organisation  `json:"OwnedOrganisations" gorm:"foreignKey:OwnerID;references:ID;constraint:OnDelete:RESTRICT"`
	GlobalAdmin        bool             `json:"GlobalAdmin" gorm:"default:false;"`
	Appeals            []Appeal         `json:"Appeal" gorm:"foreignKey:Creator;references:ID;constraint:OnDelete:CASCADE"`
	AppealResponses    []AppealResponse `json:"AppealResponses" gorm:"foreignKey:Author;references:ID;constraint:OnDelete:CASCADE"`
	PremiumType        int              `json:"PremiumType"  gorm:"type:tinyint;default:0;"`
}

type Organisation struct {
	Base
	Name            string           `json:"Name"`
	Url             string           `json:"Url" gorm:"uniqueIndex;type:char(50);"`
	IconHash        *string          `json:"IconHash"`
	Description     string           `json:"Description"`
	Moderators      []*User          `json:"Moderators" gorm:"many2many:organisation_moderators;"`
	OwnerID         uuid.UUID        `json:"Owner"`
	Verified        bool             `json:"Verified"`
	AppealTemplates []AppealTemplate `json:"AppealTemplates" gorm:"foreignKey:Organisation;references:ID;constraint:OnDelete:CASCADE"`
	Appeals         []Appeal         `json:"Appeal" gorm:"foreignKey:Organisation;references:ID;constraint:OnDelete:CASCADE"`
}

type AppealTemplate struct {
	Base
	Organisation         uuid.UUID             `json:"Organisation"`
	Name                 string                `json:"Name"`
	Appeals              []Appeal              `json:"Appeal" gorm:"foreignKey:Template;references:ID;constraint:OnDelete:CASCADE"`
	AppealTemplateFields []AppealTemplateField `json:"AppealTemplateFields" gorm:"foreignKey:Template;references:ID;constraint:OnDelete:CASCADE"`
}

type AppealTemplateField struct {
	Base
	Template       uuid.UUID      `json:"Template"`
	Title          string         `json:"Title"`
	Type           string         `json:"Type"`
	CharacterLimit int            `json:"CharacterLimit"`
	Description    string         `json:"Description"`
	Placeholder    string         `json:"Placeholder"`
	AppealAnswers  []AppealAnswer `json:"AppealAnswers" gorm:"foreignKey:Field;references:ID;constraint:OnDelete:CASCADE"`
}

type Appeal struct {
	Base
	Organisation  uuid.UUID        `json:"Organisation"`
	Creator       uuid.UUID        `json:"Creator"`
	Responded     bool             `json:"Responded"`
	Responses     []AppealResponse `json:"Responses" gorm:"foreignKey:Appeal;references:ID;constraint:OnDelete:CASCADE"`
	Content       json.RawMessage  `json:"Content"`
	Template      uuid.UUID        `json:"Template"`
	AppealStatus  int              `json:"AppealStatus" gorm:"type:tinyint;default:0;"`
	AppealAnswers []AppealAnswer   `json:"AppealAnswers" gorm:"foreignKey:Appeal;references:ID;constraint:OnDelete:CASCADE"`
}

type AppealResponse struct {
	Base
	Appeal   uuid.UUID `json:"Appeal"`
	Author   uuid.UUID `json:"Author"`
	Content  string    `json:"Content"`
	Decision int       `json:"Decision" gorm:"type:tinyint;default:0;"`
}

type AppealAnswer struct {
	Base
	Appeal  uuid.UUID `json:"Appeal"`
	Field   uuid.UUID `json:"Field"`
	Type    string    `json:"Type"`
	Content string    `json:"Content"`
}

type Base struct {
	gorm.Model
	ID uuid.UUID `json:"ID" gorm:"type:char(36);primary_key;uniqueIndex"`
}

func (base *Base) BeforeCreate(tx *gorm.DB) (err error) {
	base.ID = uuid.New()
	return
}
