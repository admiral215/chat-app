package entities

import (
	"github.com/google/uuid"
)

type User struct {
	BaseModel
	Id       uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	Email    string    `json:"email" gorm:"unique"`
	Username string    `json:"username" gorm:"unique"`
	Password string    `json:"password"`

	Groups []*GroupMember `json:"groups" gorm:"foreignKey:UserId"`
}
