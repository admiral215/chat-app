package usecases

import (
	"chat-app/internal/domain/entities"
	domain "chat-app/internal/domain/interfaces"
	"chat-app/internal/dto"
	"chat-app/pkg/jwt"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type authUseCase struct {
	userRepo   domain.UserRepository
	jwtService jwt.JWTService
}

func NewAuthUseCase(
	userRepo domain.UserRepository,
	jwtService jwt.JWTService,
) domain.AuthUseCase {
	return &authUseCase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (a authUseCase) Login(ctx context.Context, dto dto.LoginRequest) (string, error) {
	var existUser *entities.User
	existUser, err := a.userRepo.GetByUsernameOrEmail(ctx, dto.Username)

	if existUser == nil {
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(existUser.Password), []byte(dto.Password))
	if err != nil {
		return "", err
	}

	token, err := a.jwtService.GenerateToken(existUser)
	if err != nil {
		return "", err
	}
	return token, nil
}
