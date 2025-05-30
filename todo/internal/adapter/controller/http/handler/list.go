package handler

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

type ListHandler struct {
	uc  *usecase.List
	log *slog.Logger
}

func NewListHandler(uc *usecase.List, log *slog.Logger) *ListHandler {
	return &ListHandler{
		uc:  uc,
		log: log,
	}
}

func (h *ListHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.log.With(slog.String("operation", "handler.List.Create"))

		userID, err := auth.ExtractUserID(c)
		if err != nil {
			log.Debug("Unauthorized user")
			return
		}
		log = log.With("userID", userID)

		var req entity.CreateListIn
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Debug("Bad request")
			c.JSON(http.StatusBadRequest, gin.H{"Invalid input": err.Error()})
			return
		}

		listID, err := h.uc.Create(c.Request.Context(), userID, req)
		if err != nil {
			if errors.Is(err, usecase.ErrListTitleNotValid) {
				log.Debug("List title not valid")
				c.JSON(http.StatusBadRequest, gin.H{"Invalid list title": err.Error()})
				return
			}
			log.Error("Internal error", slog.Any("error", err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Debug("Created list", slog.Any("listID", listID))
		c.JSON(http.StatusOK, gin.H{"listID": listID})
	}
}

func (h *ListHandler) UpdateByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.log.With(slog.String("operation", "handler.List.UpdateByID"))

		userID, err := auth.ExtractUserID(c)
		if err != nil {
			log.Debug("Unauthorized user")
			return
		}
		log = log.With("userID", userID)

		listID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			log.Debug("Bad request")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "List ID param is not valid"})
			return
		}

		var in entity.UpdateListIn
		if err := c.ShouldBindJSON(&in); err != nil {
			log.Debug("Bad request")
			c.JSON(http.StatusBadRequest, gin.H{"Invalid list format": err.Error()})
			return
		}

		if err := h.uc.Update(c.Request.Context(), userID, in); err != nil {
			if errors.Is(err, repository.ErrListNotFound) {
				log.Debug("List not found")
				c.AbortWithStatusJSON(http.StatusNotFound, "List not found")
				return
			}
			log.Error("Internal error", slog.Any("error", err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Debug("Updated list", slog.Any("listID", listID))
		c.JSON(http.StatusOK, gin.H{"list updated": listID})
	}
}

func (h *ListHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.log.With(slog.String("operation", "handler.GetByID"))

		userID, err := auth.ExtractUserID(c)
		if err != nil {
			log.Debug("Unauthorized user")
			return
		}
		log = log.With("userID", userID)

		listID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			log.Debug("Bad request")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "List ID is not valid"})
			return
		}

		list, err := h.uc.Get(c, userID, listID)
		if err != nil {
			if errors.Is(err, repository.ErrListNotFound) {
				log.Debug("List not found")
				c.AbortWithStatusJSON(http.StatusNotFound, "List not found")
				return
			}
			log.Error("Internal error", slog.Any("error", err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Debug("Got list by ID", slog.Any("listID", listID))
		c.JSON(http.StatusOK, gin.H{"list": list})
	}
}

func (h *ListHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.log.With(slog.String("operation", "handler.GetAll"))
		userID, err := auth.ExtractUserID(c)
		if err != nil {
			log.Debug("Unauthorized user")
			return
		}
		log = log.With("userID", userID)

		lists, err := h.uc.GetAll(c.Request.Context(), userID)
		if err != nil {
			log.Error("Internal error", slog.Any("error", err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Debug("Got lists")
		if len(lists) == 0 {
			c.JSON(http.StatusOK, "No lists found")
		} else {
			c.JSON(http.StatusOK, lists)
		}
	}
}

func (h *ListHandler) DeleteByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := h.log.With(slog.String("operation", "handler.DeleteByID"))
		userID, err := auth.ExtractUserID(c)
		if err != nil {
			log.Debug("Unauthorized user")
			return
		}
		log = log.With("userID", userID)

		listID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			log.Debug("Bad request")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "List ID is not valid"})
			return
		}
		if err := h.uc.Delete(c.Request.Context(), userID, listID); err != nil {
			if errors.Is(err, repository.ErrListNotFound) {
				log.Debug("List not found")
				c.AbortWithStatusJSON(http.StatusNotFound, "List not found")
				return
			}
			log.Error("Internal error", slog.Any("error", err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Debug("Deleted list", slog.Any("listID", listID))
		c.JSON(http.StatusOK, "List deleted")
	}
}
