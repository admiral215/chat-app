package repositories

import (
	"chat-app/internal/domain/entities"
	domain "chat-app/internal/domain/interfaces"
	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) domain.MessageRepository {
	return messageRepository{
		db: db,
	}
}

func (m messageRepository) SaveMessage(message *entities.Message) error {
	return m.db.Create(&message).Error
}

func (m messageRepository) GetUndeliveredMessages(userId string) ([]entities.Message, error) {
	var messages []entities.Message
	err := m.db.
		Joins("LEFT JOIN message_statuses ON messages.id = message_statuses.message_id").
		Where("message_statuses.recipient_id = ? AND (message_statuses.delivered_at IS NULL)", userId).
		Find(&messages).Error
	return messages, err
}

func (m messageRepository) CreateMessageStatus(status *entities.MessageStatus) error {
	return m.db.Create(status).Error
}

func (m messageRepository) UpdateMessageStatus(status *entities.MessageStatus) error {
	return m.db.Save(&status).Error
}
