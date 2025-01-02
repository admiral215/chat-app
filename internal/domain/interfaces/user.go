package domain

import (
	"chat-app/internal/domain/entities"
	"chat-app/internal/dto"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	GetByUsernameOrEmail(ctx context.Context, username string) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
}

type UserUseCase interface {
	CreateUser(ctx context.Context, user *dto.UserCreate) error
}
