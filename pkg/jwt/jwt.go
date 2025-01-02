package jwt

import (
	"chat-app/config"
	"chat-app/internal/domain/entities"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTService interface {
	GenerateToken(user *entities.User) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	GetUserIDFromToken(token string) (string, error)
}

type jwtService struct {
	config *config.Config
}

func NewJWTService(config *config.Config) JWTService {
	return &jwtService{
		config: config,
	}
}

func (j jwtService) GenerateToken(user *entities.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.Id,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Duration(j.config.JWT.ExpireTime) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.config.JWT.Secret))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.config.JWT.Secret), nil
	})
}

func (j jwtService) GetUserIDFromToken(tokenString string) (string, error) {
	token, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token claims")
	}

	userID := claims["user_id"].(string)
	return userID, nil
}
