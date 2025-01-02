package entities

import "github.com/google/uuid"

type Group struct {
	BaseModel
	Id   uuid.UUID `gorm:"primary_key;type:uuid;default:gen_random_uuid()"`
	Name string    `gorm:"type:varchar(255);not null"`

	Members []*GroupMember `gorm:"foreignKey:GroupId"`
}
