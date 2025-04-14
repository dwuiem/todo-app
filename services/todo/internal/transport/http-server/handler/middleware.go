package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "No Authorization header")
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid Authorization header")
		return
	}
	// JWT Token parse
	userID, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set("userID", userID)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get("userID")
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, "No user id found")
		return -1, errors.New("No user id found")
	}
	res, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid type of user id")
		return -1, errors.New("Invalid type of user id")
	}
	return res, nil
}
