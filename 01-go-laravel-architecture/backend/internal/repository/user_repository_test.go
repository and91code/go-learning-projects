//go:build unit

package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/seu-usuario/taskflow-backend/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newMockDB cria um *sql.DB falso e garante que todas as expectativas foram cumpridas ao final do teste.
func newMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, mock.ExpectationsWereMet())
		db.Close()
	})
	return db, mock
}

var errDatabase = errors.New("database unavailable")

// Cenário: Create grava o usuário e preenche o ID gerado.
func TestUserRepository_Create_Success(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("INSERT INTO users").
		WithArgs("Ana", "ana@example.com", "hash").
		WillReturnResult(sqlmock.NewResult(1, 1))
	user := &types.User{Name: "Ana", Email: "ana@example.com", PasswordHash: "hash"}

	err := NewUserRepository(db).Create(context.Background(), user)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.ID)
}

// Cenário: Create propaga erro do banco.
func TestUserRepository_Create_Error(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("INSERT INTO users").WillReturnError(errDatabase)

	err := NewUserRepository(db).Create(context.Background(), &types.User{})

	assert.ErrorIs(t, err, errDatabase)
}

// Cenário: FindByEmail devolve o usuário encontrado.
func TestUserRepository_FindByEmail_Success(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM users WHERE email").
		WithArgs("ana@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password_hash", "created_at"}).
			AddRow(1, "Ana", "ana@example.com", "hash", time.Now()))

	user, err := NewUserRepository(db).FindByEmail(context.Background(), "ana@example.com")

	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.ID)
}

// Cenário: FindByEmail converte "sem linhas" em ErrUserNotFound.
func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM users WHERE email").WillReturnError(sql.ErrNoRows)

	user, err := NewUserRepository(db).FindByEmail(context.Background(), "x@example.com")

	assert.ErrorIs(t, err, types.ErrUserNotFound)
	assert.Nil(t, user)
}

// Cenário: FindByEmail propaga erro inesperado do banco.
func TestUserRepository_FindByEmail_Error(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM users WHERE email").WillReturnError(errDatabase)

	_, err := NewUserRepository(db).FindByEmail(context.Background(), "x@example.com")

	assert.ErrorIs(t, err, errDatabase)
}

// Cenário: FindByID devolve o usuário encontrado.
func TestUserRepository_FindByID_Success(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM users WHERE id").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password_hash", "created_at"}).
			AddRow(1, "Ana", "ana@example.com", "hash", time.Now()))

	user, err := NewUserRepository(db).FindByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, "Ana", user.Name)
}

// Cenário: FindByID converte "sem linhas" em ErrUserNotFound.
func TestUserRepository_FindByID_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM users WHERE id").WillReturnError(sql.ErrNoRows)

	_, err := NewUserRepository(db).FindByID(context.Background(), 1)

	assert.ErrorIs(t, err, types.ErrUserNotFound)
}

// Cenário: FindByID propaga erro inesperado do banco.
func TestUserRepository_FindByID_Error(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM users WHERE id").WillReturnError(errDatabase)

	_, err := NewUserRepository(db).FindByID(context.Background(), 1)

	assert.ErrorIs(t, err, errDatabase)
}
