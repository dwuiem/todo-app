package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"todo-app/internal/model"
)

func (h *Handler) signUp(c *gin.Context) {
	var input model.User
	if err := c.BindJSON(&input); err != nil {
		newErrorResponce(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

type signInInput struct {
}

func (h *Handler) signIn(c *gin.Context) {
	var input model.User
	if err := c.BindJSON(&input); err != nil {
		newErrorResponce(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}
