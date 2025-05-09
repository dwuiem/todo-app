package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ExtractUserID(c *gin.Context) (uuid.UUID, error) {
	return uuid.Parse(c.Query("userID"))
}
