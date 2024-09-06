package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) createList(c *gin.Context) {
	id, ok := c.Get("userID")
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "No UserID found")
		return
	}
	var input todo
}

func (h *Handler) getAllLists(c *gin.Context) {
}

func (h *Handler) getList(c *gin.Context) {

}

func (h *Handler) updateList(c *gin.Context) {

}

func (h *Handler) deleteList(c *gin.Context) {

}
