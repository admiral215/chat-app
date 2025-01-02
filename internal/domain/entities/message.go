package entities

import (
	"time"
)

type TypeMessage string

const (
	TypePrivate TypeMessage = "private"
	TypeGroup   TypeMessage = "group"
)

type Message struct {
	BaseModel
	Id          uint32      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Type        TypeMessage `json:"type" gorm:"type:varchar(10);check:type in ('private','group');not null"`
	SenderId    string      `json:"sender_id"`
	Content     string      `json:"content" gorm:"type:text"`
	RecipientId string      `json:"recipient_id"`

	Sender   User            `gorm:"foreignKey:SenderId"`
	Statuses []MessageStatus `json:"statuses" gorm:"foreignKey:MessageId"`
}

type MessageStatus struct {
	MessageId   uint32     `json:"message_id" gorm:"primary_key"`
	RecipientId string     `json:"recipient_id" gorm:"primary_key"`
	DeliveredAt *time.Time `json:"delivered_at"`
	ReadAt      *time.Time `json:"read_at"`

	Message Message `json:"message" gorm:"foreignKey:MessageId"`
}
