//go:build unit

package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/seu-usuario/taskflow-backend/internal/types"
	"github.com/stretchr/testify/assert"
)

var listColumns = []string{"id", "user_id", "title", "description", "created_at"}

// Cenário: Create usa o userID informado (e não o da struct) e preenche o ID gerado.
func TestListRepository_Create_Success(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("INSERT INTO lists").
		WithArgs(int64(7), "Casa", "Desc").
		WillReturnResult(sqlmock.NewResult(3, 1))
	list := &types.List{Title: "Casa", Description: "Desc"}

	err := NewListRepository(db).Create(context.Background(), 7, list)

	assert.NoError(t, err)
	assert.Equal(t, int64(3), list.ID)
}

// Cenário: Create propaga erro do banco.
func TestListRepository_Create_Error(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("INSERT INTO lists").WillReturnError(errDatabase)

	err := NewListRepository(db).Create(context.Background(), 7, &types.List{})

	assert.ErrorIs(t, err, errDatabase)
}

// Cenário: FindByIDAndUserID devolve a lista e trata descrição NULL como vazia.
func TestListRepository_FindByIDAndUserID_Success(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM lists WHERE id = \\? AND user_id = \\?").
		WithArgs(3, 7).
		WillReturnRows(sqlmock.NewRows(listColumns).AddRow(3, 7, "Casa", nil, time.Now()))

	list, err := NewListRepository(db).FindByIDAndUserID(context.Background(), 3, 7)

	assert.NoError(t, err)
	assert.Equal(t, "Casa", list.Title)
	assert.Empty(t, list.Description)
}

// Cenário: FindByIDAndUserID converte "sem linhas" em ErrListNotFound.
func TestListRepository_FindByIDAndUserID_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM lists").WillReturnError(sql.ErrNoRows)

	list, err := NewListRepository(db).FindByIDAndUserID(context.Background(), 3, 7)

	assert.ErrorIs(t, err, types.ErrListNotFound)
	assert.Nil(t, list)
}

// Cenário: FindByIDAndUserID propaga erro inesperado do banco.
func TestListRepository_FindByIDAndUserID_Error(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM lists").WillReturnError(errDatabase)

	_, err := NewListRepository(db).FindByIDAndUserID(context.Background(), 3, 7)

	assert.ErrorIs(t, err, errDatabase)
}

// Cenário: FindByUserID devolve todas as listas do usuário.
func TestListRepository_FindByUserID_Success(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM lists WHERE user_id").
		WithArgs(7).
		WillReturnRows(sqlmock.NewRows(listColumns).
			AddRow(1, 7, "A", "d", time.Now()).
			AddRow(2, 7, "B", nil, time.Now()))

	lists, err := NewListRepository(db).FindByUserID(context.Background(), 7)

	assert.NoError(t, err)
	assert.Len(t, lists, 2)
	assert.Equal(t, "d", lists[0].Description)
}

// Cenário: FindByUserID propaga erro da consulta.
func TestListRepository_FindByUserID_QueryError(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM lists WHERE user_id").WillReturnError(errDatabase)

	_, err := NewListRepository(db).FindByUserID(context.Background(), 7)

	assert.ErrorIs(t, err, errDatabase)
}

// Cenário: FindByUserID propaga erro ao ler uma linha (tipo incompatível).
func TestListRepository_FindByUserID_ScanError(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM lists WHERE user_id").
		WillReturnRows(sqlmock.NewRows(listColumns).AddRow("não-numérico", 7, "A", "d", time.Now()))

	_, err := NewListRepository(db).FindByUserID(context.Background(), 7)

	assert.Error(t, err)
}

// Cenário: Update filtra por id e userID recebidos.
func TestListRepository_Update_Success(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("UPDATE lists SET").
		WithArgs("Novo", "Desc", int64(3), int64(7)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := NewListRepository(db).Update(context.Background(), 3, 7, &types.List{Title: "Novo", Description: "Desc"})

	assert.NoError(t, err)
}

// Cenário: Update sem linhas afetadas mas com lista existente (valores iguais no MySQL) não é erro.
func TestListRepository_Update_NoChanges(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("UPDATE lists SET").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT .* FROM lists WHERE id").
		WillReturnRows(sqlmock.NewRows(listColumns).AddRow(3, 7, "Novo", "Desc", time.Now()))

	err := NewListRepository(db).Update(context.Background(), 3, 7, &types.List{Title: "Novo"})

	assert.NoError(t, err)
}

// Cenário: Update sem linhas afetadas e lista inexistente devolve ErrListNotFound.
func TestListRepository_Update_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("UPDATE lists SET").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT .* FROM lists WHERE id").WillReturnError(sql.ErrNoRows)

	err := NewListRepository(db).Update(context.Background(), 3, 7, &types.List{})

	assert.ErrorIs(t, err, types.ErrListNotFound)
}

// Cenário: Update propaga erro do banco.
func TestListRepository_Update_Error(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("UPDATE lists SET").WillReturnError(errDatabase)

	err := NewListRepository(db).Update(context.Background(), 3, 7, &types.List{})

	assert.ErrorIs(t, err, errDatabase)
}

// Cenário: Delete remove a lista do usuário.
func TestListRepository_Delete_Success(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("DELETE FROM lists").
		WithArgs(3, 7).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := NewListRepository(db).Delete(context.Background(), 3, 7)

	assert.NoError(t, err)
}

// Cenário: Delete sem linhas afetadas devolve ErrListNotFound.
func TestListRepository_Delete_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("DELETE FROM lists").WillReturnResult(sqlmock.NewResult(0, 0))

	err := NewListRepository(db).Delete(context.Background(), 3, 7)

	assert.ErrorIs(t, err, types.ErrListNotFound)
}

// Cenário: Delete propaga erro do banco.
func TestListRepository_Delete_Error(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("DELETE FROM lists").WillReturnError(errDatabase)

	err := NewListRepository(db).Delete(context.Background(), 3, 7)

	assert.ErrorIs(t, err, errDatabase)
}
