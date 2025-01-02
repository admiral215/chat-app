package usecases

import (
	"chat-app/internal/domain/entities"
	domain "chat-app/internal/domain/interfaces"
	"chat-app/internal/dto"
	"context"
)

type groupUseCase struct {
	groupRepo domain.GroupRepository
}

func NewGroupUseCase(groupRepo domain.GroupRepository) domain.GroupUseCase {
	return &groupUseCase{
		groupRepo: groupRepo,
	}
}

func (g groupUseCase) CreateGroup(ctx context.Context, creator string, dto *dto.GroupCreate) error {
	var groupMembers []*entities.GroupMember

	adminGroupMember := entities.GroupMember{
		UserId: creator,
		Role:   entities.Admin,
	}

	groupMembers = append(groupMembers, &adminGroupMember)

	for _, member := range dto.Members {
		newGroupMember := entities.GroupMember{
			UserId: member,
			Role:   entities.Member,
		}
		groupMembers = append(groupMembers, &newGroupMember)
	}

	newGroup := entities.Group{
		Name:    dto.Name,
		Members: groupMembers,
	}

	return g.groupRepo.CreateGroup(ctx, newGroup)
}
