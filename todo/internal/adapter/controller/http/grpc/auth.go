package grpc

import (
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"todo/gen/sso"
	"todo/internal/adapter/sso/grpc"
)

type AuthHandler struct {
	client *grpc.Client
	appID  int
	log    *slog.Logger
}

func NewAuthHandler(client *grpc.Client, log *slog.Logger, appID int) *AuthHandler {
	return &AuthHandler{
		client: client,
		appID:  appID,
		log:    log,
	}
}

func (h *AuthHandler) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "grpc.SignUp"
		log := h.log.With(slog.String("op", op))

		type signInRequest struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		var in signInRequest
		if err := c.ShouldBindJSON(&in); err != nil {
			log.Debug("Bad request")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Register Request is not valid"})
			return
		}
		req := sso.RegisterRequest{
			Username: in.Username,
			Password: in.Password,
		}
		id, err := h.client.Register(context.Background(), &req)
		if err != nil {
			log.Error("Failed to register user", slog.Any("error", err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Debug("Successfully registered user", slog.Any("id", id))
		c.JSON(http.StatusOK, gin.H{
			"id": id,
		})
	}
}

func (h *AuthHandler) SignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "grpc.SignIn"
		log := h.log.With(slog.String("op", op))

		type signUpRequest struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		var in signUpRequest
		if err := c.ShouldBindJSON(&in); err != nil {
			log.Debug("Bad request")
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
			log.Error("Failed to login", slog.Any("error", err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Debug("Successfully logged in")
		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}
