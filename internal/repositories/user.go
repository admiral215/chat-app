package repositories

import (
	"chat-app/internal/domain/entities"
	domain "chat-app/internal/domain/interfaces"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (u userRepository) GetByUsernameOrEmail(ctx context.Context, username string) (*entities.User, error) {
	var user entities.User
	err := u.db.WithContext(ctx).Where("username = ?", username).Or("email = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return &user, nil
}

func (u userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := u.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (u userRepository) Create(ctx context.Context, user *entities.User) error {
	err := u.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return err
	}
	return nil
}
