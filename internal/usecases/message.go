package usecases

import (
	"chat-app/internal/domain/entities"
	domain "chat-app/internal/domain/interfaces"
	"chat-app/internal/dto"
	"time"
)

type messageUseCase struct {
	messageRepo domain.MessageRepository
	groupRepo   domain.GroupRepository
}

func NewMessageUseCase(
	messageRepo domain.MessageRepository,
	groupRepo domain.GroupRepository,
) domain.MessageUseCase {
	return &messageUseCase{
		messageRepo: messageRepo,
		groupRepo:   groupRepo,
	}
}

func (m messageUseCase) GetGroup(groupId string) (*entities.Group, error) {
	return m.groupRepo.GetGroup(groupId)
}

func (m messageUseCase) HandleNewConnection(userID string) ([]entities.Message, error) {
	messages, err := m.messageRepo.GetUndeliveredMessages(userID)
	if err != nil {
		return nil, err
	}

	// Update status delivered untuk semua pesan
	currentTime := time.Now()
	for _, msg := range messages {
		status := entities.MessageStatus{
			MessageId:   msg.Id,
			RecipientId: userID,
			DeliveredAt: &currentTime,
		}
		if err := m.messageRepo.UpdateMessageStatus(&status); err != nil {
			return nil, err
		}
	}

	return messages, nil
}

func (m messageUseCase) HandleReadReceipt(messageId uint32, recipientId string) error {
	currentTime := time.Now()
	status := entities.MessageStatus{
		MessageId:   messageId,
		RecipientId: recipientId,
		ReadAt:      &currentTime,
	}
	return m.messageRepo.UpdateMessageStatus(&status)
}

func (m messageUseCase) HandleNewMessage(senderId string, wsMsg dto.WSMessage) (*entities.Message, error) {
	msg := &entities.Message{
		Type:        wsMsg.Type,
		SenderId:    senderId,
		Content:     wsMsg.Content,
		RecipientId: wsMsg.RecipientId,
	}

	if err := m.messageRepo.SaveMessage(msg); err != nil {
		return nil, err
	}

	if msg.Type != entities.TypeGroup {
		status := entities.MessageStatus{
			MessageId:   msg.Id,
			RecipientId: msg.RecipientId,
		}

		if err := m.messageRepo.CreateMessageStatus(&status); err != nil {
			return nil, err
		}
	} else {
		group, err := m.groupRepo.GetGroup(msg.RecipientId)
		if err != nil {
			return nil, err
		}
		for _, member := range group.Members {
			if msg.SenderId == member.UserId {
				continue
			}

			status := entities.MessageStatus{
				MessageId:   msg.Id,
				RecipientId: member.UserId,
			}

			if err := m.messageRepo.CreateMessageStatus(&status); err != nil {
				return nil, err
			}
		}
	}

	return msg, nil
}

func (m messageUseCase) HandleMessageDelivery(msg *entities.Message, recipientId string) error {
	currentTime := time.Now()
	status := entities.MessageStatus{
		MessageId:   msg.Id,
		RecipientId: recipientId,
		DeliveredAt: &currentTime,
	}
	return m.messageRepo.UpdateMessageStatus(&status)
}
