package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"todo/gen/sso"
)

func (h *Handler) signUp(c *gin.Context) {
	var input sso.RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.client.Register(context.Background(), &input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (h *Handler) signIn(c *gin.Context) {
	var input sso.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := h.client.Login(context.Background(), &input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
