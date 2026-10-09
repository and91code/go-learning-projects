package dto

import (
	"time"

	"github.com/seu-usuario/taskflow-backend/internal/types"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func (req RegisterRequest) ToDomain() types.RegisterInput {
	return types.RegisterInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
}

func (req LoginRequest) ToDomain() types.LoginInput {
	return types.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}
}

func AuthResponseFromDomain(result types.AuthResult) AuthResponse {
	return AuthResponse{
		Token: result.Token,
		User:  UserResponseFromDomain(result.User),
	}
}

func UserResponseFromDomain(user types.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
}
