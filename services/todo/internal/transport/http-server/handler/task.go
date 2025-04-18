package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"todo/internal/domain/model"
	"todo/internal/storage"
)

type TaskStorage interface {
	Create(c *gin.Context, task model.Task) (int64, error)
	Update(c *gin.Context, task model.Task) (int64, error)
	GetAllByListID(c *gin.Context, listID int64) ([]model.Task, error)
	GetByID(c *gin.Context, user, taskID int64) (model.Task, error)
	DeleteByID(c *gin.Context, taskID int64) error
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := h.service.List.GetByID(c, userID, listID); err != nil {
		if errors.Is(err, storage.ErrListNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, "List not found")
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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
		ListID:      listID,
	}

	taskID, err := h.service.Task.Create(c, task)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"taskID": taskID})
}

func (h *Handler) updateTask(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Task ID is not valid"})
		return
	}
	task, err := h.service.Task.GetByID(c, userID, taskID)
	if err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, "Task not found")
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var in taskRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"invalid task format": err.Error()})
	}
	newTask := model.Task{
		ID:          taskID,
		Title:       in.Title,
		Description: in.Description,
		Completed:   in.Completed,
		ListID:      task.ListID,
	}
	if _, err := h.service.Task.Update(c, newTask); err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, "Task not found")
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task updated": taskID})
}

func (h *Handler) getAllTasksByList(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}
	listID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "List ID is not valid"})
		return
	}
	if _, err := h.service.List.GetByID(c, userID, listID); err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, "Task not found")
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tasks, err := h.service.Task.GetAllByListID(c, listID)
	if err != nil {
		if errors.Is(err, storage.ErrListNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "List not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(tasks) == 0 {
		c.JSON(http.StatusOK, tasks)
	} else {
		c.JSON(http.StatusOK, "List is empty")
	}
}

func (h *Handler) getTaskByID(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Task ID is not valid"})
		return
	}
	task, err := h.service.Task.GetByID(c, userID, taskID)
	if err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, "Task not found")
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": task})
}

func (h *Handler) deleteTask(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Task ID is not valid"})
		return
	}
	if _, err := h.service.Task.GetByID(c, userID, taskID); err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, "Task not found")
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Task.DeleteByID(c, taskID); err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, "Task not found")
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "Task deleted")
}
