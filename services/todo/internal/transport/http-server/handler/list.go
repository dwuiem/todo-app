package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"todo/internal/domain/model"
)

type ListService interface {
	Create(list model.List) (int64, error)
	Update(list model.List) error
	GetAllByUserID(userID int64) ([]model.List, error)
	GetByID(userID int64, listID int64) (model.List, error)
	DeleteByID(listID int64) error
}

type listRequest struct {
	Title string `json:"title"`
}

type listResponse struct {
	Title string `json:"title"`
}

type listItemResponse struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type listsResponse struct {
	list []listItemResponse
}

func (h *Handler) getAllLists(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	// TODO: Handle Error
	lists, err := h.service.List.GetAllByUserID(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	listsResponse := listsResponse{
		list: make([]listItemResponse, len(lists)),
	}
	for i, list := range lists {
		listsResponse.list[i] = listItemResponse{
			ID:    list.ID,
			Title: list.Title,
		}
	}
	if len(lists) == 0 {
		c.JSON(http.StatusOK, "You don't have any lists")
	} else {
		c.JSON(http.StatusOK, listsResponse)
	}
}

func (h *Handler) getListByID(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	listID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid list ID")
		return
	}

	// TODO: Handle error
	list, err := h.service.List.GetByID(userID, listID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := listResponse{
		Title: list.Title,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) updateList(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	listID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid list ID")
		return
	}

	if _, err := h.service.List.GetByID(userID, listID); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, "list not found")
		return
	}

	var req listRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"invalid update request": err.Error()})
		return
	}

	list := model.List{
		ID:    listID,
		Title: req.Title,
	}

	// TODO: Handle error
	if err := h.service.List.Update(list); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"updated list": list})
}

func (h *Handler) deleteList(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	listID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "invalid list ID")
		return
	}

	if _, err := h.service.List.GetByID(userID, listID); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, "list not found")
		return
	}

	if err := h.service.List.DeleteByID(listID); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deleted list": listID})
}
