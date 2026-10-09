//go:build unit

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/seu-usuario/taskflow-backend/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task *types.Task) error {
	return m.Called(ctx, task).Error(0)
}

func (m *MockTaskRepository) FindByIDAndUserID(ctx context.Context, taskID, userID int64) (*types.Task, error) {
	args := m.Called(ctx, taskID, userID)
	task, _ := args.Get(0).(*types.Task)
	return task, args.Error(1)
}

func (m *MockTaskRepository) FindByListIDAndUserID(ctx context.Context, listID, userID int64) ([]types.Task, error) {
	args := m.Called(ctx, listID, userID)
	tasks, _ := args.Get(0).([]types.Task)
	return tasks, args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *types.Task, userID int64) error {
	return m.Called(ctx, task, userID).Error(0)
}

func (m *MockTaskRepository) Delete(ctx context.Context, taskID, userID int64) error {
	return m.Called(ctx, taskID, userID).Error(0)
}

type MockListRepository struct {
	mock.Mock
}

func (m *MockListRepository) Create(ctx context.Context, list *types.List) error {
	return m.Called(ctx, list).Error(0)
}

func (m *MockListRepository) FindByIDAndUserID(ctx context.Context, id, userID int64) (*types.List, error) {
	args := m.Called(ctx, id, userID)
	list, _ := args.Get(0).(*types.List)
	return list, args.Error(1)
}

func (m *MockListRepository) FindByUserID(ctx context.Context, userID int64) ([]types.List, error) {
	args := m.Called(ctx, userID)
	lists, _ := args.Get(0).([]types.List)
	return lists, args.Error(1)
}

func (m *MockListRepository) Update(ctx context.Context, list *types.List) error {
	return m.Called(ctx, list).Error(0)
}

func (m *MockListRepository) Delete(ctx context.Context, id, userID int64) error {
	return m.Called(ctx, id, userID).Error(0)
}

func newTaskService(taskRepo *MockTaskRepository, listRepo *MockListRepository) types.TaskService {
	return NewTaskService(taskRepo, listRepo)
}

func TestCreateTask_Success(t *testing.T) {
	taskRepo := new(MockTaskRepository)
	listRepo := new(MockListRepository)
	listRepo.On("FindByIDAndUserID", mock.Anything, int64(8), int64(7)).
		Return(&types.List{ID: 8, UserID: 7}, nil).Once()
	taskRepo.On("Create", mock.Anything, mock.MatchedBy(func(task *types.Task) bool {
		return task.ListID == 8 && task.Title == "Read a book" &&
			task.Priority == types.PriorityHigh && task.Status == types.StatusPending
	})).Run(func(args mock.Arguments) {
		args.Get(1).(*types.Task).ID = 42
	}).Return(nil).Once()
	svc := newTaskService(taskRepo, listRepo)

	result, err := svc.CreateTask(context.Background(), 7, types.Task{
		ListID: 8, Title: "Read a book", Priority: types.PriorityHigh,
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(42), result.ID)
	assert.Equal(t, int64(8), result.ListID)
	assert.Equal(t, "Read a book", result.Title)
	assert.Equal(t, types.PriorityHigh, result.Priority)
	assert.Equal(t, types.StatusPending, result.Status)
	taskRepo.AssertExpectations(t)
	listRepo.AssertExpectations(t)
}

func TestCreateTask_InvalidPriority(t *testing.T) {
	taskRepo := new(MockTaskRepository)
	listRepo := new(MockListRepository)
	svc := newTaskService(taskRepo, listRepo)

	result, err := svc.CreateTask(context.Background(), 7, types.Task{
		ListID: 8, Title: "Read a book", Priority: types.TaskPriority("urgent"),
	})

	assert.ErrorIs(t, err, types.ErrInvalidTask)
	assert.Empty(t, result)
	listRepo.AssertNotCalled(t, "FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything)
	taskRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestCreateTask_EmptyTitle(t *testing.T) {
	taskRepo := new(MockTaskRepository)
	listRepo := new(MockListRepository)
	svc := newTaskService(taskRepo, listRepo)

	result, err := svc.CreateTask(context.Background(), 7, types.Task{
		ListID: 8, Title: "   ", Priority: types.PriorityMedium,
	})

	assert.ErrorIs(t, err, types.ErrInvalidTask)
	assert.Empty(t, result)
	taskRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestCreateTask_RepositoryError(t *testing.T) {
	repositoryErr := errors.New("database unavailable")
	taskRepo := new(MockTaskRepository)
	listRepo := new(MockListRepository)
	listRepo.On("FindByIDAndUserID", mock.Anything, int64(8), int64(7)).
		Return(&types.List{ID: 8, UserID: 7}, nil).Once()
	taskRepo.On("Create", mock.Anything, mock.Anything).Return(repositoryErr).Once()
	svc := newTaskService(taskRepo, listRepo)

	result, err := svc.CreateTask(context.Background(), 7, types.Task{
		ListID: 8, Title: "Read a book", Priority: types.PriorityLow,
	})

	assert.ErrorIs(t, err, repositoryErr)
	assert.Empty(t, result)
	taskRepo.AssertExpectations(t)
	listRepo.AssertExpectations(t)
}
