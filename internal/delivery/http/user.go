package http

import (
	domain "chat-app/internal/domain/interfaces"
	"chat-app/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

type UserHandler struct {
	userUseCase domain.UserUseCase
}

func NewUserHandler(userUseCase domain.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// Register godoc
// @Tags Auth
// @Summary Register
// @Description Register user
// @Param userDto body dto.UserCreate true "User credentials"
// @Router /api/auth/register [post]
func (handler *UserHandler) Register(c *gin.Context) {
	var userDto dto.UserCreate
	if err := c.ShouldBindJSON(&userDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(userDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	context.WithTimeout(ctx, 5*time.Second)

	err := handler.userUseCase.CreateUser(ctx, &userDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully!"})
}
