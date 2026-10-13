package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seu-usuario/taskflow-backend/internal/dto"
	"github.com/seu-usuario/taskflow-backend/internal/middleware"
	"github.com/seu-usuario/taskflow-backend/internal/types"
)

type TagHandler struct {
	tagService types.TagService
}

func NewTagHandler(s types.TagService) *TagHandler {
	return &TagHandler{tagService: s}
}

func (h *TagHandler) Store(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req dto.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	tag, err := h.tagService.CreateTag(c.Request.Context(), userID, req.ToDomain())
	if err != nil {
		writeTagServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.TagDetailResponseFromDomain(tag))
}

func (h *TagHandler) IndexByList(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	listID, err := strconv.ParseInt(c.Query("list_id"), 10, 64)
	if err != nil || listID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de lista inválido"})
		return
	}

	tags, err := h.tagService.GetTagsByList(c.Request.Context(), listID, userID)
	if err != nil {
		writeTagServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.TagResponsesFromDomain(tags))
}

func (h *TagHandler) Show(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tagID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || tagID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de tag inválido"})
		return
	}

	tag, err := h.tagService.GetTagById(c.Request.Context(), tagID, userID)
	if err != nil {
		writeTagServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.TagDetailResponseFromDomain(tag))
}

func (h *TagHandler) Update(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tagID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || tagID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de tag inválido"})
		return
	}

	var req dto.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	tag, err := h.tagService.UpdateTag(c.Request.Context(), tagID, userID, req.ToDomain())
	if err != nil {
		writeTagServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.TagDetailResponseFromDomain(tag))
}

func (h *TagHandler) Delete(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tagID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || tagID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de tag inválido"})
		return
	}

	if err := h.tagService.DeleteTag(c.Request.Context(), tagID, userID); err != nil {
		writeTagServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func writeTagServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, types.ErrTagNotFound), errors.Is(err, types.ErrListNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, types.ErrInvalidTag):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
