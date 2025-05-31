package grpc

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"todo/internal/adapter/controller/http/auth"
	"todo/internal/adapter/repository"
	"todo/internal/domain/entity"
	"todo/internal/domain/usecase"
)

type TaskHandler struct {
	uc  *usecase.Task
	log *slog.Logger
}

func NewTaskHandler(us *usecase.Task, log *slog.Logger) *TaskHandler {
	return &TaskHandler{
		uc:  us,
		log: log,
	}
}

func (h *TaskHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.log.With(slog.String("operation", "handler.Task.Create"))

		userID, err := auth.ExtractUserID(c)
		if err != nil {
			log.Debug("Unauthorized user")
			return
		}
		log = log.With("userID", userID)

		listID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			log.Debug("Bad request")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var req entity.CreateTaskIn
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Debug("Bad request")
			c.JSON(http.StatusBadRequest, gin.H{"Invalid task": err.Error()})
			return
		}

		taskID, err := h.uc.Create(c.Request.Context(), userID, listID, req)
		if err != nil {
			log.Error("Internal error", slog.Any("error", err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Debug("Task created", slog.Any("taskID", taskID))
		c.JSON(http.StatusOK, gin.H{"taskID": taskID})
	}
}

func (h *TaskHandler) UpdateByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.log.With(slog.String("operation", "handler.Task.Create"))

		userID, err := auth.ExtractUserID(c)
		if err != nil {
			log.Debug("Unauthorized user")
			return
		}
		taskID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Task ID is not valid"})
			return
		}
		var in entity.UpdateTaskIn
		if err := c.ShouldBindJSON(&in); err != nil {
			log.Debug("Bad request")
			c.JSON(http.StatusBadRequest, gin.H{"Invalid task format": err.Error()})
		}

		if err := h.uc.Update(c.Request.Context(), userID, taskID, in); err != nil {
			if errors.Is(err, repository.ErrTaskNotFound) {
				log.Debug("Task not found")
				c.AbortWithStatusJSON(http.StatusNotFound, "Task not found")
				return
			}
			log.Error("Internal error", slog.Any("error", err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Debug("Task updated", slog.Any("taskID", taskID))
		c.JSON(http.StatusOK, gin.H{"task updated": taskID})
	}
}

func (h *TaskHandler) GetByListID() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.log.With(slog.String("operation", "handler.GetByListID"))

		userID, err := auth.ExtractUserID(c)
		if err != nil {
			log.Debug("Unauthorized user")
			return
		}
		listID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			log.Debug("Bad request")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "List ID is not valid"})
			return
		}

		tasks, err := h.uc.GetByListID(c.Request.Context(), userID, listID)
		if err != nil {
			if errors.Is(err, repository.ErrListNotFound) {
				log.Debug("List not found")
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "List not found"})
				return
			}
			log.Error("Internal error", slog.Any("error", err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Debug("Tasks found", slog.Any("amount", len(tasks)))
		if len(tasks) == 0 {
			c.JSON(http.StatusOK, tasks)
		} else {
			c.JSON(http.StatusOK, "List is empty")
		}
	}
}

func (h *TaskHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.log.With(slog.String("operation", "handler.GetByID"))

		userID, err := auth.ExtractUserID(c)
		if err != nil {
			log.Debug("Unauthorized user")
			return
		}

		taskID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			log.Debug("Bad request")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Task ID is not valid"})
			return
		}
		task, err := h.uc.Get(c, userID, taskID)
		if err != nil {
			if errors.Is(err, repository.ErrTaskNotFound) {
				log.Debug("Task not found")
				c.AbortWithStatusJSON(http.StatusNotFound, "Task not found")
				return
			}
			log.Error("Internal error", slog.Any("error", err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Debug("Task found", slog.Any("task", task.ID))
		c.JSON(http.StatusOK, gin.H{"task": task})
	}
}

func (h *TaskHandler) DeleteByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.log.With(slog.String("operation", "handler.DeleteByID"))

		userID, err := auth.ExtractUserID(c)
		if err != nil {
			log.Debug("Unauthorized user")
			return
		}

		taskID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			log.Debug("Bad request")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Task ID is not valid"})
			return
		}

		if err := h.uc.Delete(c.Request.Context(), userID, taskID); err != nil {
			if errors.Is(err, repository.ErrTaskNotFound) {
				log.Debug("Task not found")
				c.AbortWithStatusJSON(http.StatusNotFound, "Task not found")
				return
			}
			log.Error("Internal error", slog.Any("error", err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Debug("Task deleted", slog.Any("taskID", taskID))
		c.JSON(http.StatusOK, "Task deleted")
	}
}
