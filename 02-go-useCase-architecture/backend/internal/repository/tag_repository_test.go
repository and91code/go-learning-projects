//go:build unit

package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/seu-usuario/taskflow-backend/internal/types"
	"github.com/stretchr/testify/assert"
)

var tagColumns = []string{"id", "list_id", "name", "color"}

// Cenário: Create insere a tag somente se a lista pertencer ao userID e preenche o ID gerado.
func TestTagRepository_Create_Success(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("INSERT INTO tags .* FROM lists l").
		WithArgs("urgente", "#ff0000", int64(8), int64(7)).
		WillReturnResult(sqlmock.NewResult(5, 1))
	tag := &types.Tag{ListID: 8, Name: "urgente", Color: "#ff0000"}

	err := NewTagRepository(db).Create(context.Background(), 7, tag)

	assert.NoError(t, err)
	assert.Equal(t, int64(5), tag.ID)
}

// Cenário: Create sem linhas inseridas (lista de outro usuário) devolve ErrListNotFound.
func TestTagRepository_Create_ListNotFound(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("INSERT INTO tags").WillReturnResult(sqlmock.NewResult(0, 0))

	err := NewTagRepository(db).Create(context.Background(), 7, &types.Tag{ListID: 8})

	assert.ErrorIs(t, err, types.ErrListNotFound)
}

// Cenário: Create propaga erro do banco.
func TestTagRepository_Create_Error(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("INSERT INTO tags").WillReturnError(errDatabase)

	err := NewTagRepository(db).Create(context.Background(), 7, &types.Tag{})

	assert.ErrorIs(t, err, errDatabase)
}

// Cenário: Update filtra por tagID e userID recebidos.
func TestTagRepository_Update_Success(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("UPDATE tags t").
		WithArgs(int64(9), "novo", "#00ff00", int64(5), int64(7)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := NewTagRepository(db).Update(context.Background(), 5, 7, &types.Tag{ListID: 9, Name: "novo", Color: "#00ff00"})

	assert.NoError(t, err)
}

// Cenário: Update sem linhas afetadas mas com tag e lista existentes (valores iguais) não é erro.
func TestTagRepository_Update_NoChanges(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("UPDATE tags t").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT t.id .* FROM tags t").
		WillReturnRows(sqlmock.NewRows(tagColumns).AddRow(5, 9, "novo", "#fff"))
	mock.ExpectQuery("SELECT id FROM lists").
		WithArgs(9, 7).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(9))

	err := NewTagRepository(db).Update(context.Background(), 5, 7, &types.Tag{ListID: 9, Name: "novo"})

	assert.NoError(t, err)
}

// Cenário: Update sem linhas afetadas e tag inexistente devolve ErrTagNotFound.
func TestTagRepository_Update_TagNotFound(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("UPDATE tags t").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT t.id .* FROM tags t").WillReturnError(sql.ErrNoRows)

	err := NewTagRepository(db).Update(context.Background(), 5, 7, &types.Tag{ListID: 9})

	assert.ErrorIs(t, err, types.ErrTagNotFound)
}

// Cenário: Update sem linhas afetadas, tag existente e nova lista de outro usuário devolve ErrListNotFound.
func TestTagRepository_Update_ListNotFound(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("UPDATE tags t").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT t.id .* FROM tags t").
		WillReturnRows(sqlmock.NewRows(tagColumns).AddRow(5, 8, "a", "#fff"))
	mock.ExpectQuery("SELECT id FROM lists").WillReturnError(sql.ErrNoRows)

	err := NewTagRepository(db).Update(context.Background(), 5, 7, &types.Tag{ListID: 9})

	assert.ErrorIs(t, err, types.ErrListNotFound)
}

// Cenário: Update propaga erro inesperado ao verificar a lista.
func TestTagRepository_Update_CheckListError(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("UPDATE tags t").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT t.id .* FROM tags t").
		WillReturnRows(sqlmock.NewRows(tagColumns).AddRow(5, 8, "a", "#fff"))
	mock.ExpectQuery("SELECT id FROM lists").WillReturnError(errDatabase)

	err := NewTagRepository(db).Update(context.Background(), 5, 7, &types.Tag{ListID: 9})

	assert.ErrorIs(t, err, errDatabase)
}

// Cenário: Update propaga erro do banco.
func TestTagRepository_Update_Error(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("UPDATE tags t").WillReturnError(errDatabase)

	err := NewTagRepository(db).Update(context.Background(), 5, 7, &types.Tag{})

	assert.ErrorIs(t, err, errDatabase)
}

// Cenário: Delete remove a tag do usuário.
func TestTagRepository_Delete_Success(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("DELETE t FROM tags t").
		WithArgs(5, 7).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := NewTagRepository(db).Delete(context.Background(), 5, 7)

	assert.NoError(t, err)
}

// Cenário: Delete sem linhas afetadas devolve ErrTagNotFound.
func TestTagRepository_Delete_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("DELETE t FROM tags t").WillReturnResult(sqlmock.NewResult(0, 0))

	err := NewTagRepository(db).Delete(context.Background(), 5, 7)

	assert.ErrorIs(t, err, types.ErrTagNotFound)
}

// Cenário: Delete propaga erro do banco.
func TestTagRepository_Delete_Error(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("DELETE t FROM tags t").WillReturnError(errDatabase)

	err := NewTagRepository(db).Delete(context.Background(), 5, 7)

	assert.ErrorIs(t, err, errDatabase)
}

// Cenário: FindByIDAndUserID devolve a tag do usuário.
func TestTagRepository_FindByIDAndUserID_Success(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT t.id .* FROM tags t").
		WithArgs(5, 7).
		WillReturnRows(sqlmock.NewRows(tagColumns).AddRow(5, 8, "urgente", "#ff0000"))

	tag, err := NewTagRepository(db).FindByIDAndUserID(context.Background(), 5, 7)

	assert.NoError(t, err)
	assert.Equal(t, types.Tag{ID: 5, ListID: 8, Name: "urgente", Color: "#ff0000"}, tag)
}

// Cenário: FindByIDAndUserID converte "sem linhas" em ErrTagNotFound.
func TestTagRepository_FindByIDAndUserID_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT t.id .* FROM tags t").WillReturnError(sql.ErrNoRows)

	tag, err := NewTagRepository(db).FindByIDAndUserID(context.Background(), 5, 7)

	assert.ErrorIs(t, err, types.ErrTagNotFound)
	assert.Empty(t, tag)
}

// Cenário: FindByIDAndUserID propaga erro inesperado do banco.
func TestTagRepository_FindByIDAndUserID_Error(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT t.id .* FROM tags t").WillReturnError(errDatabase)

	_, err := NewTagRepository(db).FindByIDAndUserID(context.Background(), 5, 7)

	assert.ErrorIs(t, err, errDatabase)
}

// Cenário: FindByListIDAndUserID devolve as tags da lista.
func TestTagRepository_FindByListIDAndUserID_Success(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT t.id .* FROM tags t").
		WithArgs(8, 7).
		WillReturnRows(sqlmock.NewRows(tagColumns).
			AddRow(1, 8, "a", "#fff").
			AddRow(2, 8, "b", "#000"))

	tags, err := NewTagRepository(db).FindByListIDAndUserID(context.Background(), 8, 7)

	assert.NoError(t, err)
	assert.Len(t, tags, 2)
	assert.Equal(t, "b", tags[1].Name)
}

// Cenário: FindByListIDAndUserID propaga erro da consulta.
func TestTagRepository_FindByListIDAndUserID_QueryError(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT t.id .* FROM tags t").WillReturnError(errDatabase)

	_, err := NewTagRepository(db).FindByListIDAndUserID(context.Background(), 8, 7)

	assert.ErrorIs(t, err, errDatabase)
}

// Cenário: FindByListIDAndUserID propaga erro ao ler uma linha (tipo incompatível).
func TestTagRepository_FindByListIDAndUserID_ScanError(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT t.id .* FROM tags t").
		WillReturnRows(sqlmock.NewRows(tagColumns).AddRow("não-numérico", 8, "a", "#fff"))

	_, err := NewTagRepository(db).FindByListIDAndUserID(context.Background(), 8, 7)

	assert.Error(t, err)
}

// Cenário: FindByListIDAndUserID propaga erro de iteração das linhas.
func TestTagRepository_FindByListIDAndUserID_RowsError(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT t.id .* FROM tags t").
		WillReturnRows(sqlmock.NewRows(tagColumns).AddRow(1, 8, "a", "#fff").RowError(0, errDatabase))

	_, err := NewTagRepository(db).FindByListIDAndUserID(context.Background(), 8, 7)

	assert.ErrorIs(t, err, errDatabase)
}
