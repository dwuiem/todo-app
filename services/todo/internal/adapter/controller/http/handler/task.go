package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"todo/internal/adapter/controller/http/auth"
	"todo/internal/adapter/repository"
	"todo/internal/domain/entity"
	"todo/internal/domain/usecase"
)

type TaskHandler struct {
	uc *usecase.Task
}

func NewTaskHandler(us *usecase.Task) *TaskHandler {
	return &TaskHandler{uc: us}
}

func (h *TaskHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := auth.ExtractUserID(c)
		if err != nil {
			return
		}
		listID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var req entity.CreateTaskIn
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Invalid task": err.Error()})
			return
		}

		taskID, err := h.uc.Create(c.Request.Context(), userID, listID, req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"taskID": taskID})
	}
}

func (h *TaskHandler) UpdateByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := auth.ExtractUserID(c)
		if err != nil {
			return
		}
		taskID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Task ID is not valid"})
			return
		}
		var in entity.UpdateTaskIn
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Invalid task format": err.Error()})
		}

		if err := h.uc.Update(c.Request.Context(), userID, taskID, in); err != nil {
			if errors.Is(err, repository.ErrTaskNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, "Task not found")
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"uc updated": taskID})
	}
}

func (h *TaskHandler) GetByListID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := auth.ExtractUserID(c)
		if err != nil {
			return
		}
		listID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "List ID is not valid"})
			return
		}

		tasks, err := h.uc.GetByListID(c.Request.Context(), userID, listID)
		if err != nil {
			if errors.Is(err, repository.ErrListNotFound) {
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
}

func (h *TaskHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := auth.ExtractUserID(c)
		if err != nil {
			return
		}
		taskID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Task ID is not valid"})
			return
		}
		task, err := h.uc.Get(c, userID, taskID)
		if err != nil {
			if errors.Is(err, repository.ErrTaskNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, "Task not found")
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"uc": task})
	}
}

func (h *TaskHandler) DeleteByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := auth.ExtractUserID(c)
		if err != nil {
			return
		}
		taskID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Task ID is not valid"})
			return
		}
		if err := h.uc.Delete(c.Request.Context(), userID, taskID); err != nil {
			if errors.Is(err, repository.ErrTaskNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, "Task not found")
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, "Task deleted")
	}
}
