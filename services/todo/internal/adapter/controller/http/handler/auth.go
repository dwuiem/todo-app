package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"todo/gen/sso"
	"todo/internal/adapter/sso/grpc"
)

type AuthHandler struct {
	client *grpc.Client
	appID  int
}

func NewAuthHandler(client *grpc.Client, appID int) *AuthHandler {
	return &AuthHandler{
		client: client,
		appID:  appID,
	}
}

func (h *AuthHandler) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		type signInRequest struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		var in signInRequest
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
}

func (h *AuthHandler) SignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		type signUpRequest struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		var in signUpRequest
		if err := c.ShouldBindJSON(&in); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Login Request is not valid"})
			return
		}
		req := sso.LoginRequest{
			Username: in.Username,
			Password: in.Password,
			AppId:    int32(h.appID),
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
}
