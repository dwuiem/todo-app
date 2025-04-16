package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"todo/internal/domain/model"
)

type TaskService interface {
	Create(task model.Task) (int64, error)
	Update(task model.Task) (int64, error)
	GetAllByListID(listID int64) ([]model.Task, error)
	GetByID(taskID int64) (model.Task, error)
	Delete(taskID int64) error
}

type taskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func (h *Handler) createTask(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}
	listID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid List ID")
		return
	}

	if _, err := h.service.List.GetByID(userID, listID); err != nil {
		newErrorResponse(c, http.StatusNotFound, "List not Found")
	}

	var req taskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid task": err.Error()})
		return
	}

	task := model.Task{
		Title:       req.Title,
		Description: req.Description,
		Completed:   req.Completed,
		ListID:      int64(listID),
	}

	// TODO: Handle error
	taskID, err := h.service.Task.Create(task)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"taskID": taskID})
}

// TODO
func (h *Handler) updateTask(c *gin.Context) {}

func (h *Handler) getAllTasksByList(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}
	listID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid List ID")
		return
	}
	if _, err := h.service.List.GetByID(userID, listID); err != nil {
		newErrorResponse(c, http.StatusNotFound, "List not Found")
	}

	// TODO: Handle error
	tasks, err := h.service.Task.GetAllByListID(listID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// TODO
func (h *Handler) getTaskByID(c *gin.Context) {}

// TODO
func (h *Handler) deleteTask(c *gin.Context) {}
