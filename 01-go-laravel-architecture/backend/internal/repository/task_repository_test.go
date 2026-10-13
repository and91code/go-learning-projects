//go:build unit

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/seu-usuario/taskflow-backend/internal/types"
	"github.com/stretchr/testify/assert"
)

var taskColumns = []string{"id", "list_id", "title", "description", "status", "priority", "due_date", "created_at", "updated_at", "tag_id", "tag_name", "tag_color"}

// Cenário: Create insere a tarefa somente se a lista pertencer ao userID e preenche o ID gerado.
func TestTaskRepository_Create_Success(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO tasks .* FROM lists l").
		WithArgs("Ler", "Livro", types.StatusPending, types.PriorityHigh, nil, int64(8), int64(7)).
		WillReturnResult(sqlmock.NewResult(42, 1))
	mock.ExpectCommit()
	task := &types.Task{ListID: 8, Title: "Ler", Description: "Livro", Status: types.StatusPending, Priority: types.PriorityHigh}

	err := NewTaskRepository(db).Create(context.Background(), 7, task)

	assert.NoError(t, err)
	assert.Equal(t, int64(42), task.ID)
}

func TestTaskRepository_Create_WithTags(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO tasks .* FROM lists l").
		WillReturnResult(sqlmock.NewResult(42, 1))
	mock.ExpectQuery("SELECT id, list_id, name, color FROM tags").
		WithArgs(int64(8), int64(5)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "list_id", "name", "color"}).
			AddRow(5, 8, "urgent", "#ff0000"))
	mock.ExpectExec("INSERT INTO task_tags").
		WithArgs(int64(42), int64(5)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	task := &types.Task{
		ListID: 8, Title: "Ler", Status: types.StatusPending, Priority: types.PriorityHigh,
		Tags: []types.Tag{{ID: 5}},
	}

	err := NewTaskRepository(db).Create(context.Background(), 7, task)

	assert.NoError(t, err)
	assert.Equal(t, []types.Tag{{ID: 5, ListID: 8, Name: "urgent", Color: "#ff0000"}}, task.Tags)
}

func TestTaskRepository_Create_InvalidTagRollsBack(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO tasks .* FROM lists l").
		WillReturnResult(sqlmock.NewResult(42, 1))
	mock.ExpectQuery("SELECT id, list_id, name, color FROM tags").
		WithArgs(int64(8), int64(5)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "list_id", "name", "color"}))
	mock.ExpectRollback()

	err := NewTaskRepository(db).Create(context.Background(), 7, &types.Task{
		ListID: 8, Title: "Ler", Tags: []types.Tag{{ID: 5}},
	})

	assert.ErrorIs(t, err, types.ErrTagNotFound)
}

// Cenário: Create sem linhas inseridas (lista de outro usuário) devolve ErrListNotFound.
func TestTaskRepository_Create_ListNotFound(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO tasks").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	err := NewTaskRepository(db).Create(context.Background(), 7, &types.Task{ListID: 8})

	assert.ErrorIs(t, err, types.ErrListNotFound)
}

// Cenário: Create propaga erro do banco.
func TestTaskRepository_Create_Error(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO tasks").WillReturnError(errDatabase)
	mock.ExpectRollback()

	err := NewTaskRepository(db).Create(context.Background(), 7, &types.Task{})

	assert.ErrorIs(t, err, errDatabase)
}

// Cenário: Update filtra por taskID e userID recebidos.
func TestTaskRepository_Update_Success(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE tasks t").
		WithArgs("Novo", "Desc", types.StatusCompleted, types.PriorityLow, nil, int64(42), int64(7)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	task := &types.Task{Title: "Novo", Description: "Desc", Status: types.StatusCompleted, Priority: types.PriorityLow}

	err := NewTaskRepository(db).Update(context.Background(), 42, 7, task)

	assert.NoError(t, err)
}

// Cenário: Update sem linhas afetadas e tarefa inexistente devolve ErrTaskNotFound.
func TestTaskRepository_Update_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE tasks t").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery("SELECT t.id").
		WithArgs(int64(42), int64(7)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectRollback()

	err := NewTaskRepository(db).Update(context.Background(), 42, 7, &types.Task{})

	assert.ErrorIs(t, err, types.ErrTaskNotFound)
}

// Cenário: Update propaga erro do banco.
func TestTaskRepository_Update_Error(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE tasks t").WillReturnError(errDatabase)
	mock.ExpectRollback()

	err := NewTaskRepository(db).Update(context.Background(), 42, 7, &types.Task{})

	assert.ErrorIs(t, err, errDatabase)
}

func TestTaskRepository_Update_ReplacesTags(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE tasks t").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT id, list_id, name, color FROM tags").
		WithArgs(int64(8), int64(5)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "list_id", "name", "color"}).
			AddRow(5, 8, "urgent", "#ff0000"))
	mock.ExpectExec("DELETE FROM task_tags WHERE task_id = ?").
		WithArgs(int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO task_tags").
		WithArgs(int64(42), int64(5)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	task := &types.Task{
		ID: 42, ListID: 8, Tags: []types.Tag{{ID: 5}},
	}

	err := NewTaskRepository(db).Update(context.Background(), 42, 7, task)

	assert.NoError(t, err)
	assert.Equal(t, []types.Tag{{ID: 5, ListID: 8, Name: "urgent", Color: "#ff0000"}}, task.Tags)
}

func TestTaskRepository_Update_InvalidTagRollsBack(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE tasks t").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery("SELECT id, list_id, name, color FROM tags").
		WithArgs(int64(8), int64(5)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "list_id", "name", "color"}))
	mock.ExpectRollback()

	err := NewTaskRepository(db).Update(context.Background(), 42, 7, &types.Task{
		ID: 42, ListID: 8, Tags: []types.Tag{{ID: 5}},
	})

	assert.ErrorIs(t, err, types.ErrTagNotFound)
}

func TestTaskRepository_Update_EmptyTagsClearsLinks(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE tasks t").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("DELETE FROM task_tags WHERE task_id = ?").
		WithArgs(int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	task := &types.Task{ID: 42, ListID: 8, Tags: []types.Tag{}}

	err := NewTaskRepository(db).Update(context.Background(), 42, 7, task)

	assert.NoError(t, err)
	assert.Empty(t, task.Tags)
}

// Cenário: Delete remove a tarefa do usuário.
func TestTaskRepository_Delete_Success(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("DELETE t FROM tasks t").
		WithArgs(42, 7).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := NewTaskRepository(db).Delete(context.Background(), 42, 7)

	assert.NoError(t, err)
}

// Cenário: Delete sem linhas afetadas devolve ErrTaskNotFound.
func TestTaskRepository_Delete_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("DELETE t FROM tasks t").WillReturnResult(sqlmock.NewResult(0, 0))

	err := NewTaskRepository(db).Delete(context.Background(), 42, 7)

	assert.ErrorIs(t, err, types.ErrTaskNotFound)
}

// Cenário: Delete propaga erro do banco.
func TestTaskRepository_Delete_Error(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectExec("DELETE t FROM tasks t").WillReturnError(errDatabase)

	err := NewTaskRepository(db).Delete(context.Background(), 42, 7)

	assert.ErrorIs(t, err, errDatabase)
}

// Cenário: FindByIDAndUserID mapeia descrição e prazo preenchidos.
func TestTaskRepository_FindByIDAndUserID_Success(t *testing.T) {
	dueDate := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM tasks t").
		WithArgs(42, 7).
		WillReturnRows(sqlmock.NewRows(taskColumns).
			AddRow(42, 8, "Ler", "Livro", "pending", "high", dueDate, time.Now(), time.Now(), 5, "urgente", "#ff0000"))

	task, err := NewTaskRepository(db).FindByIDAndUserID(context.Background(), 42, 7)

	assert.NoError(t, err)
	assert.Equal(t, "Livro", task.Description)
	assert.Equal(t, types.PriorityHigh, task.Priority)
	assert.NotNil(t, task.DueDate)
	assert.Equal(t, []types.Tag{{ID: 5, ListID: 8, Name: "urgente", Color: "#ff0000"}}, task.Tags)
}

// Cenário: FindByIDAndUserID trata descrição e prazo NULL como vazio/nil.
func TestTaskRepository_FindByIDAndUserID_NullFields(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM tasks t").
		WillReturnRows(sqlmock.NewRows(taskColumns).
			AddRow(42, 8, "Ler", nil, "pending", "high", nil, time.Now(), time.Now(), nil, nil, nil))

	task, err := NewTaskRepository(db).FindByIDAndUserID(context.Background(), 42, 7)

	assert.NoError(t, err)
	assert.Empty(t, task.Description)
	assert.Nil(t, task.DueDate)
}

// Cenário: FindByIDAndUserID converte "sem linhas" em ErrTaskNotFound.
func TestTaskRepository_FindByIDAndUserID_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM tasks t").
		WillReturnRows(sqlmock.NewRows(taskColumns))

	task, err := NewTaskRepository(db).FindByIDAndUserID(context.Background(), 42, 7)

	assert.ErrorIs(t, err, types.ErrTaskNotFound)
	assert.Nil(t, task)
}

// Cenário: FindByIDAndUserID propaga erro inesperado do banco.
func TestTaskRepository_FindByIDAndUserID_Error(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM tasks t").WillReturnError(errDatabase)

	_, err := NewTaskRepository(db).FindByIDAndUserID(context.Background(), 42, 7)

	assert.ErrorIs(t, err, errDatabase)
}

// Cenário: FindByListIDAndUserID devolve todas as tarefas da lista.
func TestTaskRepository_FindByListIDAndUserID_Success(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM tasks t").
		WithArgs(8, 7).
		WillReturnRows(sqlmock.NewRows(taskColumns).
			AddRow(1, 8, "A", "d", "pending", "low", nil, time.Now(), time.Now(), 5, "urgent", "#ff0000").
			AddRow(1, 8, "A", "d", "pending", "low", nil, time.Now(), time.Now(), 6, "home", "#00ff00").
			AddRow(2, 8, "B", nil, "completed", "high", time.Now(), time.Now(), time.Now(), nil, nil, nil))

	tasks, err := NewTaskRepository(db).FindByListIDAndUserID(context.Background(), 8, 7)

	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	assert.NotNil(t, tasks[1].DueDate)
	assert.Equal(t, []types.Tag{
		{ID: 5, ListID: 8, Name: "urgent", Color: "#ff0000"},
		{ID: 6, ListID: 8, Name: "home", Color: "#00ff00"},
	}, tasks[0].Tags)
	assert.Empty(t, tasks[1].Tags)
}

// Cenário: FindByListIDAndUserID propaga erro da consulta.
func TestTaskRepository_FindByListIDAndUserID_QueryError(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM tasks t").WillReturnError(errDatabase)

	_, err := NewTaskRepository(db).FindByListIDAndUserID(context.Background(), 8, 7)

	assert.ErrorIs(t, err, errDatabase)
}

// Cenário: FindByListIDAndUserID propaga erro ao ler uma linha (tipo incompatível).
func TestTaskRepository_FindByListIDAndUserID_ScanError(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery("SELECT .* FROM tasks t").
		WillReturnRows(sqlmock.NewRows(taskColumns).
			AddRow("não-numérico", 8, "A", "d", "pending", "low", nil, time.Now(), time.Now(), nil, nil, nil))

	_, err := NewTaskRepository(db).FindByListIDAndUserID(context.Background(), 8, 7)

	assert.Error(t, err)
}
