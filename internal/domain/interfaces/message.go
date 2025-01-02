package domain

import (
	"chat-app/internal/domain/entities"
	"chat-app/internal/dto"
)

type MessageRepository interface {
	SaveMessage(message *entities.Message) error
	GetUndeliveredMessages(userId string) ([]entities.Message, error)
	CreateMessageStatus(status *entities.MessageStatus) error
	UpdateMessageStatus(status *entities.MessageStatus) error
}

type MessageUseCase interface {
	HandleNewConnection(userID string) ([]entities.Message, error)
	HandleReadReceipt(messageId uint32, recipientId string) error
	HandleNewMessage(senderId string, wsMsg dto.WSMessage) (*entities.Message, error)
	HandleMessageDelivery(msg *entities.Message, recipientId string) error
	GetGroup(groupId string) (*entities.Group, error)
}
