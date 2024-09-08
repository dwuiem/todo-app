package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"todo-app/internal/model"
)

func (h *Handler) createItem(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}
	listID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid list ID param")
		return
	}
	var input model.Item
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.Item.Create(userID, listID, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllItems(c *gin.Context) {

}

func (h *Handler) getItem(c *gin.Context) {

}

func (h *Handler) updateItem(c *gin.Context) {

}

func (h *Handler) deleteItem(c *gin.Context) {

}
