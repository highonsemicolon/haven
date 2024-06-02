package models

import (
	"github.com/google/uuid"
)

type Repo struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey" notnull:"true"`
	Name    string    `json:"name" gorm:"unique" notnull:"true"`
	GitURL  string    `json:"git_url" notnull:"true"`
	Comment string    `json:"comment" gorm:"default:null"`
}
