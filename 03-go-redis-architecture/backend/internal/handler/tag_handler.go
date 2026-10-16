package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seu-usuario/taskflow-backend/internal/dto"
	"github.com/seu-usuario/taskflow-backend/internal/middleware"
	"github.com/seu-usuario/taskflow-backend/internal/queue"
	"github.com/seu-usuario/taskflow-backend/internal/types"
)

type TagHandler struct {
	create interface {
		Execute(context.Context, int64, types.Tag) (types.Tag, error)
	}
	update interface {
		Execute(context.Context, int64, int64, types.Tag) (types.Tag, error)
	}
	delete interface {
		Execute(context.Context, int64, int64) error
	}
	byList interface {
		Execute(context.Context, int64, int64) ([]types.Tag, error)
	}
	byID interface {
		Execute(context.Context, int64, int64) (types.Tag, error)
	}
	jobs queue.Producer
}

func (h *TagHandler) SetQueue(q queue.Producer) { h.jobs = q }

func NewTagHandler(create interface {
	Execute(context.Context, int64, types.Tag) (types.Tag, error)
}, update interface {
	Execute(context.Context, int64, int64, types.Tag) (types.Tag, error)
}, delete interface {
	Execute(context.Context, int64, int64) error
}, byList interface {
	Execute(context.Context, int64, int64) ([]types.Tag, error)
}, byID interface {
	Execute(context.Context, int64, int64) (types.Tag, error)
}) *TagHandler {
	return &TagHandler{create: create, update: update, delete: delete, byList: byList, byID: byID}
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

	if h.jobs != nil {
		domain := req.ToDomain()
		if err := queue.ValidateTag("TAG_CREATE", userID, 0, domain); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		payload, err := json.Marshal(domain)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to encode tag"})
			return
		}
		if err := h.jobs.Enqueue(c.Request.Context(), queue.Job{Action: "TAG_CREATE", UserID: userID, Payload: payload}); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "unable to queue tag"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"status": "accepted"})
		return
	}
	tag, err := h.create.Execute(c.Request.Context(), userID, req.ToDomain())
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

	tags, err := h.byList.Execute(c.Request.Context(), listID, userID)
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

	tag, err := h.byID.Execute(c.Request.Context(), tagID, userID)
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

	if h.jobs != nil {
		domain := req.ToDomain()
		if err := queue.ValidateTag("TAG_UPDATE", userID, tagID, domain); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		payload, err := json.Marshal(domain)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to encode tag"})
			return
		}
		if err := h.jobs.Enqueue(c.Request.Context(), queue.Job{Action: "TAG_UPDATE", UserID: userID, EntityID: tagID, Payload: payload}); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "unable to queue tag"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"status": "accepted"})
		return
	}
	tag, err := h.update.Execute(c.Request.Context(), tagID, userID, req.ToDomain())
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

	if h.jobs != nil {
		if err := h.jobs.Enqueue(c.Request.Context(), queue.Job{Action: "TAG_DELETE", UserID: userID, EntityID: tagID}); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "unable to queue tag"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"status": "accepted"})
		return
	}
	if err := h.delete.Execute(c.Request.Context(), tagID, userID); err != nil {
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
