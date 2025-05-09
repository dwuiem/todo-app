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

type ListHandler struct {
	uc *usecase.List
}

func NewListHandler(us *usecase.List) *ListHandler {
	return &ListHandler{uc: us}
}

func (h *ListHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := auth.ExtractUserID(c)
		if err != nil {
			return
		}

		var req entity.CreateListIn
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Invalid list": err.Error()})
			return
		}

		listID, err := h.uc.Create(c.Request.Context(), userID, req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"listID": listID})
	}
}

func (h *ListHandler) UpdateByID() gin.HandlerFunc {
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

		var in entity.UpdateListIn
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Invalid list format": err.Error()})
			return
		}

		if err := h.uc.Update(c.Request.Context(), userID, in); err != nil {
			if errors.Is(err, repository.ErrListNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, "List not found")
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"list updated": listID})
	}
}

func (h *ListHandler) GetByID() gin.HandlerFunc {
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

		list, err := h.uc.Get(c, userID, listID)
		if err != nil {
			if errors.Is(err, repository.ErrListNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, "List not found")
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"list": list})
	}
}

func (h *ListHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := auth.ExtractUserID(c)
		if err != nil {
			return
		}

		lists, err := h.uc.GetAll(c.Request.Context(), userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(lists) == 0 {
			c.JSON(http.StatusOK, "No lists found")
		} else {
			c.JSON(http.StatusOK, lists)
		}
	}
}

func (h *ListHandler) DeleteByID() gin.HandlerFunc {
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
		if err := h.uc.Delete(c.Request.Context(), userID, listID); err != nil {
			if errors.Is(err, repository.ErrListNotFound) {
				c.AbortWithStatusJSON(http.StatusNotFound, "List not found")
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, "List deleted")
	}
}
