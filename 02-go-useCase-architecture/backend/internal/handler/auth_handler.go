package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seu-usuario/taskflow-backend/internal/dto"
	"github.com/seu-usuario/taskflow-backend/internal/middleware"
	"github.com/seu-usuario/taskflow-backend/internal/types"
)

type AuthHandler struct {
	register interface {
		Execute(context.Context, types.RegisterInput) (types.AuthResult, error)
	}
	login interface {
		Execute(context.Context, types.LoginInput) (types.AuthResult, error)
	}
	getUser interface {
		Execute(context.Context, int64) (types.User, error)
	}
	updateProfile interface {
		Execute(context.Context, int64, types.UpdateProfileInput) (types.User, error)
	}
	updatePassword interface {
		Execute(context.Context, int64, types.UpdatePasswordInput) error
	}
}

func NewAuthHandler(register interface {
	Execute(context.Context, types.RegisterInput) (types.AuthResult, error)
}, login interface {
	Execute(context.Context, types.LoginInput) (types.AuthResult, error)
}, getUser interface {
	Execute(context.Context, int64) (types.User, error)
}, updateProfile interface {
	Execute(context.Context, int64, types.UpdateProfileInput) (types.User, error)
}, updatePassword interface {
	Execute(context.Context, int64, types.UpdatePasswordInput) error
}) *AuthHandler {
	return &AuthHandler{register: register, login: login, getUser: getUser, updateProfile: updateProfile, updatePassword: updatePassword}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	result, err := h.register.Execute(c.Request.Context(), req.ToDomain())
	if err != nil {
		if errors.Is(err, types.ErrEmailExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.AuthResponseFromDomain(result))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	result, err := h.login.Execute(c.Request.Context(), req.ToDomain())
	if err != nil {
		if errors.Is(err, types.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponseFromDomain(result))
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	user, err := h.updateProfile.Execute(c.Request.Context(), userID, req.ToDomain())
	if err != nil {
		if errors.Is(err, types.ErrEmailExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, types.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.UserResponseFromDomain(user))
}

func (h *AuthHandler) UpdatePassword(c *gin.Context) {
	var req dto.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Dados inválidos", "details": err.Error()})
		return
	}

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	err = h.updatePassword.Execute(c.Request.Context(), userID, req.ToDomain())
	switch {
	case errors.Is(err, types.ErrInvalidCurrentPassword):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case errors.Is(err, types.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusOK, gin.H{"message": "Senha atualizada com sucesso"})
	}
}
