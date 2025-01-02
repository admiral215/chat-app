package helpers

import (
	"github.com/gin-gonic/gin"
	httpPkg "net/http"
)

func GetUserFromContext(c *gin.Context) string {
	userId, ok := c.Get("userId")
	if !ok {
		c.JSON(httpPkg.StatusUnauthorized, gin.H{"error": "user id not found in context"})
		return ""
	}

	userIdStr, ok := userId.(string)
	if !ok {
		c.JSON(httpPkg.StatusBadRequest, gin.H{"error": "invalid user id format"})
		return ""
	}
	return userIdStr
}
