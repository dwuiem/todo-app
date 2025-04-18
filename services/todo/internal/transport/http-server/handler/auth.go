package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"todo/gen/sso"
)

type SignRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signUp(c *gin.Context) {
	var in SignRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Register Request is not valid"})
		return
	}
	req := sso.RegisterRequest{
		Username: in.Username,
		Password: in.Password,
	}
	id, err := h.client.Register(context.Background(), &req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (h *Handler) signIn(c *gin.Context) {
	var in SignRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Login Request is not valid"})
		return
	}
	req := sso.LoginRequest{
		Username: in.Username,
		Password: in.Password,
		AppId:    int32(h.id),
	}
	token, err := h.client.Login(context.Background(), &req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
