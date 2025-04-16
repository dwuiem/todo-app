package handler

import "github.com/gin-gonic/gin"

type StatusResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, ErrorResponse{message})
}
