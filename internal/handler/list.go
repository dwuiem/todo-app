package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"todo-app/internal/model"
)

func (h *Handler) createList(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}

	var input model.List
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.List.Create(userID, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (h *Handler) getAllLists(c *gin.Context) {
}

func (h *Handler) getList(c *gin.Context) {

}

func (h *Handler) updateList(c *gin.Context) {

}

func (h *Handler) deleteList(c *gin.Context) {

}
