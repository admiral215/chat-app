package dto

import "chat-app/internal/domain/entities"

type WSMessage struct {
	Type        entities.TypeMessage `json:"type"`
	Content     string               `json:"content"`
	RecipientId string               `json:"recipient_id"`
	Action      string               `json:"action,omitempty"`
	MessageId   uint32               `json:"message_id,omitempty"`
}
