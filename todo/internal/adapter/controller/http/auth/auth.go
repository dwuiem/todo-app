package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ExtractUserID(c *gin.Context) (uuid.UUID, error) {
	id, exists := c.Get("userID")
	if !exists {
		return uuid.Nil, errors.New("no user ID in context")
	}
	return uuid.Parse(id.(string))
}
