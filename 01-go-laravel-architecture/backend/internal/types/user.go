package types

import (
	"context"
	"time"
)

type User struct {
	ID           int64
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResult struct {
	Token string
	User  User
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id int64) (*User, error)
}

type UserService interface {
	GetUserByID(ctx context.Context, id int64) (User, error)
}

type AuthService interface {
	UserService
	Register(ctx context.Context, input RegisterInput) (AuthResult, error)
	Login(ctx context.Context, input LoginInput) (AuthResult, error)
}
