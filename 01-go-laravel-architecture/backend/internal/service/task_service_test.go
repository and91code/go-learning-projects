//go:build unit

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/seu-usuario/taskflow-backend/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) CreateTask(ctx context.Context, title, priority string) (int64, error) {
	args := m.Called(ctx, title, priority)
	return args.Get(0).(int64), args.Error(1)
}

func TestCreateTask_Success(t *testing.T) {
	repo := new(MockTaskRepository)
	repo.On("CreateTask", mock.Anything, "Read a book", "high").Return(int64(42), nil).Once()
	svc := NewTaskService(repo)

	result, err := svc.CreateTask(context.Background(), dto.CreateTaskRequest{
		Title:    "Read a book",
		Priority: "high",
	})

	assert.NoError(t, err)
	assert.Equal(t, dto.TaskResponse{
		ID:       42,
		Title:    "Read a book",
		Priority: "high",
		Status:   "pending",
	}, result)
	repo.AssertExpectations(t)
}

func TestCreateTask_InvalidPriority(t *testing.T) {
	repo := new(MockTaskRepository)
	svc := NewTaskService(repo)

	result, err := svc.CreateTask(context.Background(), dto.CreateTaskRequest{
		Title:    "Read a book",
		Priority: "urgent",
	})

	assert.ErrorIs(t, err, ErrInvalidTask)
	assert.Empty(t, result)
	repo.AssertNotCalled(t, "CreateTask", mock.Anything, mock.Anything, mock.Anything)
}

func TestCreateTask_EmptyTitle(t *testing.T) {
	repo := new(MockTaskRepository)
	svc := NewTaskService(repo)

	result, err := svc.CreateTask(context.Background(), dto.CreateTaskRequest{
		Title:    "   ",
		Priority: "medium",
	})

	assert.ErrorIs(t, err, ErrInvalidTask)
	assert.Empty(t, result)
	repo.AssertNotCalled(t, "CreateTask", mock.Anything, mock.Anything, mock.Anything)
}

func TestCreateTask_RepositoryError(t *testing.T) {
	repositoryErr := errors.New("database unavailable")
	repo := new(MockTaskRepository)
	repo.On("CreateTask", mock.Anything, "Read a book", "low").Return(int64(0), repositoryErr).Once()
	svc := NewTaskService(repo)

	result, err := svc.CreateTask(context.Background(), dto.CreateTaskRequest{
		Title:    "Read a book",
		Priority: "low",
	})

	assert.ErrorIs(t, err, repositoryErr)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}
