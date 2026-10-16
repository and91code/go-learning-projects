package dto

import (
	"testing"
	"time"

	"github.com/seu-usuario/taskflow-backend/internal/types"
	"github.com/stretchr/testify/require"
)

// Cenário: RegisterRequest.ToDomain copia nome, e-mail e senha para o input de domínio.
func TestRegisterRequest_ToDomain(t *testing.T) {
	req := RegisterRequest{Name: "Ana", Email: "ana@example.com", Password: "secret1"}

	input := req.ToDomain()

	require.Equal(t, types.RegisterInput{Name: "Ana", Email: "ana@example.com", Password: "secret1"}, input)
}

// Cenário: LoginRequest.ToDomain copia e-mail e senha para o input de domínio.
func TestLoginRequest_ToDomain(t *testing.T) {
	req := LoginRequest{Email: "ana@example.com", Password: "secret1"}

	input := req.ToDomain()

	require.Equal(t, types.LoginInput{Email: "ana@example.com", Password: "secret1"}, input)
}

// Cenário: AuthResponseFromDomain devolve token e usuário com data em RFC3339 e sem expor o hash da senha.
func TestAuthResponseFromDomain(t *testing.T) {
	createdAt := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)
	result := types.AuthResult{
		Token: "jwt-token",
		User:  types.User{ID: 1, Name: "Ana", Email: "ana@example.com", PasswordHash: "hash", CreatedAt: createdAt},
	}

	response := AuthResponseFromDomain(result)

	require.Equal(t, AuthResponse{
		Token: "jwt-token",
		User:  UserResponse{ID: 1, Name: "Ana", Email: "ana@example.com", CreatedAt: "2026-08-01T10:00:00Z"},
	}, response)
}
