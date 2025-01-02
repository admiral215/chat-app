package domain

import (
	"chat-app/internal/dto"
	"context"
)

type AuthUseCase interface {
	Login(ctx context.Context, dto dto.LoginRequest) (string, error)
}
