package http

import (
	domain "chat-app/internal/domain/interfaces"
	"chat-app/internal/dto"
	"chat-app/pkg/helpers"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"time"
)

type GroupHandler struct {
	groupUseCase domain.GroupUseCase
}

func NewGroupHandler(groupUseCase domain.GroupUseCase) *GroupHandler {
	return &GroupHandler{
		groupUseCase: groupUseCase,
	}
}

// CreateGroup godoc
// @Tags Group
// @Summary Create a group chat
// @Description Create a group chat
// @Security BearerAuth
// @Param userDto body dto.GroupCreate true "group insert"
// @Router /api/groups [post]
func (g *GroupHandler) CreateGroup(c *gin.Context) {
	var groupCreate dto.GroupCreate
	if err := c.ShouldBindJSON(&groupCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(groupCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	context.WithTimeout(ctx, 5*time.Second)

	userId := helpers.GetUserFromContext(c)
	if userId == "" {
		c.JSON(500, gin.H{
			"error": "user id is empty",
		})
		return
	}

	if err := g.groupUseCase.CreateGroup(ctx, userId, &groupCreate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Group created successfully!"})
}
