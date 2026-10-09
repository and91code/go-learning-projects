package service

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

type authService struct {
	userRepo types.UserRepository
}

var _ types.AuthService = (*authService)(nil)
var _ types.UserService = (*authService)(nil)

func NewAuthService(repository types.UserRepository) types.AuthService {
	return &authService{userRepo: repository}
}

func (s *authService) GetUserByID(ctx context.Context, id int64) (types.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return types.User{}, err
	}
	return *user, nil
}

func (s *authService) Register(ctx context.Context, input types.RegisterInput) (types.AuthResult, error) {
	_, err := s.userRepo.FindByEmail(ctx, input.Email)
	switch {
	case err == nil:
		return types.AuthResult{}, types.ErrEmailExists
	case !errors.Is(err, types.ErrUserNotFound):
		return types.AuthResult{}, fmt.Errorf("check existing user: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return types.AuthResult{}, fmt.Errorf("falha ao criptografar senha: %w", err)
	}

	user := &types.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return types.AuthResult{}, err
	}

	token, err := generateJWT(user.ID)
	if err != nil {
		return types.AuthResult{}, err
	}
	return types.AuthResult{Token: token, User: *user}, nil
}

func (s *authService) Login(ctx context.Context, input types.LoginInput) (types.AuthResult, error) {
	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, types.ErrUserNotFound) {
			return types.AuthResult{}, types.ErrInvalidCredentials
		}
		return types.AuthResult{}, fmt.Errorf("find user for login: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return types.AuthResult{}, types.ErrInvalidCredentials
	}

	token, err := generateJWT(user.ID)
	if err != nil {
		return types.AuthResult{}, err
	}
	return types.AuthResult{Token: token, User: *user}, nil
}

func generateJWT(userID int64) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) == 0 {
		return "", errors.New("JWT_SECRET não configurado")
	}

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
