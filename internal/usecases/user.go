package usecases

import (
	"chat-app/internal/domain/entities"
	domain "chat-app/internal/domain/interfaces"
	"chat-app/internal/dto"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

type userUseCase struct {
	userRepo domain.UserRepository
}

func NewUserUseCase(
	userRepo domain.UserRepository,
) domain.UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (u userUseCase) CreateUser(ctx context.Context, user *dto.UserCreate) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}

	newUser := entities.User{
		Username: user.Username,
		Password: string(hashedPass),
		Email:    user.Email,
	}

	err = u.userRepo.Create(ctx, &newUser)
	if err != nil {
		return err
	}
	return nil
}
