package domain

import (
	"chat-app/internal/domain/entities"
	"chat-app/internal/dto"
	"context"
)

type GroupRepository interface {
	GetGroup(groupId string) (*entities.Group, error)
	CreateGroup(ctx context.Context, group entities.Group) error
}

type GroupUseCase interface {
	CreateGroup(ctx context.Context, creator string, dto *dto.GroupCreate) error
}
