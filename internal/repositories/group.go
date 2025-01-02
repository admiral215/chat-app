package repositories

import (
	"chat-app/internal/domain/entities"
	domain "chat-app/internal/domain/interfaces"
	"context"
	"fmt"
	"gorm.io/gorm"
)

type groupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) domain.GroupRepository {
	return &groupRepository{
		db,
	}
}

func (g groupRepository) GetGroup(groupId string) (*entities.Group, error) {
	var groups entities.Group
	err := g.db.Preload("Members").
		Where("id = ?", groupId).
		First(&groups).Error
	return &groups, err
}

func (g groupRepository) CreateGroup(ctx context.Context, group entities.Group) error {
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&group).Error; err != nil {
			return fmt.Errorf("failed to create group: %w", err)
		}
		return nil
	})
	return err
}
