package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"strings"
	"todo/internal/adapter/sso/grpc"
)

type auth struct {
	secret []byte
	client *grpc.Client
}

func AuthMiddleware(secret []byte, client *grpc.Client) gin.HandlerFunc {
	slog.Debug("AuthMiddleware")
	auth := &auth{
		secret: secret,
		client: client,
	}
	return auth.userIdentityMiddleware
}

func (a *auth) userIdentityMiddleware(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "No Authorization header")
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "Invalid Authorization header")
		return
	}

	// JWT Token parse
	userID, err := a.parseUserIDFromToken(headerParts[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err.Error())
		return
	}
	c.Set("userID", userID.String())
}

func (a *auth) parseUserIDFromToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.secret, nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			return uuid.Nil, errors.New("invalid user id claim")
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return uuid.Nil, errors.New("invalid user uuid")
		}
		return userID, nil
	}
	return uuid.Nil, errors.New("invalid token")
}
