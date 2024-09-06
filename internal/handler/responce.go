package handler

import "github.com/gin-gonic/gin"

type Error struct {
	Message string `json:"message"`
}

func newErrorResponce(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, Error{message})
}
