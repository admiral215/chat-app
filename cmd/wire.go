//go:build wireinject
// +build wireinject

package cmd

import (
	"chat-app/config"
	"chat-app/internal/delivery/http"
	"chat-app/internal/delivery/middleware"
	websocket2 "chat-app/internal/delivery/websocket"
	"chat-app/internal/repositories"
	"chat-app/internal/usecases"
	"chat-app/pkg/database"
	"chat-app/pkg/jwt"
	"github.com/google/wire"
)

func provideDatabaseConfig(cfg *config.Config) *config.DatabaseConfig {
	return &cfg.Database
}

func provideJwtConfig(cfg *config.Config) *config.JWTConfig {
	return &cfg.JWT
}

var configSet = wire.NewSet(
	provideDatabaseConfig,
	provideJwtConfig,
)

var authSet = wire.NewSet(
	jwt.NewJWTService,
	usecases.NewAuthUseCase,
)

var middlewareSet = wire.NewSet(
	middleware.NewAuthMiddleware,
)

var repositorySet = wire.NewSet(
	repositories.NewUserRepository,
	repositories.NewGroupRepository,
	repositories.NewMessageRepository,
)

var usecaseSet = wire.NewSet(
	usecases.NewUserUseCase,
	usecases.NewMessageUseCase,
	usecases.NewGroupUseCase,
)

var handlerSet = wire.NewSet(
	http.NewUserHandler,
	http.NewAuthHandler,
	http.NewGroupHandler,
)

var webSocketSet = wire.NewSet(
	websocket2.NewHub,
)

func InitializeApp(cfg *config.Config) (*App, error) {
	wire.Build(
		database.ConnectDB,
		authSet,
		webSocketSet,
		middlewareSet,

		configSet,
		repositorySet,
		usecaseSet,
		handlerSet,

		NewApp,
	)

	return nil, nil
}
