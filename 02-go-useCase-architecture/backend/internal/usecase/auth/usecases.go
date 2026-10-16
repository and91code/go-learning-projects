package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/seu-usuario/taskflow-backend/internal/types"
	"golang.org/x/crypto/bcrypt"
)

type Register interface {
	Execute(context.Context, types.RegisterInput) (types.AuthResult, error)
}
type Login interface {
	Execute(context.Context, types.LoginInput) (types.AuthResult, error)
}
type GetUser interface {
	Execute(context.Context, int64) (types.User, error)
}
type UpdateProfile interface {
	Execute(context.Context, int64, types.UpdateProfileInput) (types.User, error)
}
type UpdatePassword interface {
	Execute(context.Context, int64, types.UpdatePasswordInput) error
}

type userUseCases struct{ repository types.UserRepository }

func New(repository types.UserRepository) (Register, Login, GetUser, UpdateProfile, UpdatePassword) {
	u := &userUseCases{repository}
	return register{u}, login{u}, getUser{u}, updateProfile{u}, updatePassword{u}
}

type register struct{ *userUseCases }

func (u register) Execute(ctx context.Context, in types.RegisterInput) (types.AuthResult, error) {
	_, err := u.repository.FindByEmail(ctx, in.Email)
	if err == nil {
		return types.AuthResult{}, types.ErrEmailExists
	}
	if !errors.Is(err, types.ErrUserNotFound) {
		return types.AuthResult{}, fmt.Errorf("check existing user: %w", err)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return types.AuthResult{}, fmt.Errorf("falha ao criptografar senha: %w", err)
	}
	user := &types.User{Name: in.Name, Email: in.Email, PasswordHash: string(hash), CreatedAt: time.Now()}
	if err := u.repository.Create(ctx, user); err != nil {
		return types.AuthResult{}, err
	}
	token, err := generateJWT(user.ID)
	if err != nil {
		return types.AuthResult{}, err
	}
	return types.AuthResult{Token: token, User: *user}, nil
}

type login struct{ *userUseCases }

func (u login) Execute(ctx context.Context, in types.LoginInput) (types.AuthResult, error) {
	user, err := u.repository.FindByEmail(ctx, in.Email)
	if err != nil {
		if errors.Is(err, types.ErrUserNotFound) {
			return types.AuthResult{}, types.ErrInvalidCredentials
		}
		return types.AuthResult{}, fmt.Errorf("find user for login: %w", err)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)) != nil {
		return types.AuthResult{}, types.ErrInvalidCredentials
	}
	token, err := generateJWT(user.ID)
	if err != nil {
		return types.AuthResult{}, err
	}
	return types.AuthResult{Token: token, User: *user}, nil
}

type getUser struct{ *userUseCases }

func (u getUser) Execute(ctx context.Context, id int64) (types.User, error) {
	user, err := u.repository.FindByID(ctx, id)
	if err != nil {
		return types.User{}, err
	}
	return *user, nil
}

type updateProfile struct{ *userUseCases }

func (u updateProfile) Execute(ctx context.Context, id int64, in types.UpdateProfileInput) (types.User, error) {
	existing, err := u.repository.FindByEmail(ctx, in.Email)
	switch {
	case err == nil && existing.ID != id:
		return types.User{}, types.ErrEmailExists
	case err != nil && !errors.Is(err, types.ErrUserNotFound):
		return types.User{}, fmt.Errorf("check existing email for profile update: %w", err)
	}
	if err := u.repository.UpdateProfile(ctx, id, in.Name, in.Email); err != nil {
		return types.User{}, err
	}
	user, err := u.repository.FindByID(ctx, id)
	if err != nil {
		return types.User{}, err
	}
	return *user, nil
}

type updatePassword struct{ *userUseCases }

func (u updatePassword) Execute(ctx context.Context, id int64, in types.UpdatePasswordInput) error {
	user, err := u.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.CurrentPassword)) != nil {
		return types.ErrInvalidCurrentPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("falha ao criptografar senha: %w", err)
	}
	return u.repository.UpdatePassword(ctx, id, string(hash))
}
func generateJWT(userID int64) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) == 0 {
		return "", errors.New("JWT_SECRET não configurado")
	}
	claims := jwt.MapClaims{"sub": userID, "exp": time.Now().Add(24 * time.Hour).Unix(), "iat": time.Now().Unix()}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}
