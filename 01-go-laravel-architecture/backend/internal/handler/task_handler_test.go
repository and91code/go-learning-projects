//go:build unit

package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/seu-usuario/taskflow-backend/internal/middleware"
	"github.com/seu-usuario/taskflow-backend/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) CreateTask(ctx context.Context, userID int64, task types.Task) (types.Task, error) {
	args := m.Called(ctx, userID, task)
	result, _ := args.Get(0).(types.Task)
	return result, args.Error(1)
}

func (m *MockTaskService) GetTasksByList(ctx context.Context, listID, userID int64) ([]types.Task, error) {
	args := m.Called(ctx, listID, userID)
	result, _ := args.Get(0).([]types.Task)
	return result, args.Error(1)
}

func (m *MockTaskService) GetTaskByID(ctx context.Context, taskID, userID int64) (types.Task, error) {
	args := m.Called(ctx, taskID, userID)
	result, _ := args.Get(0).(types.Task)
	return result, args.Error(1)
}

func (m *MockTaskService) UpdateTask(ctx context.Context, taskID, userID int64, task types.Task) (types.Task, error) {
	args := m.Called(ctx, taskID, userID, task)
	result, _ := args.Get(0).(types.Task)
	return result, args.Error(1)
}

func (m *MockTaskService) DeleteTask(ctx context.Context, taskID, userID int64) error {
	return m.Called(ctx, taskID, userID).Error(0)
}

func setupTaskHandler(serviceMock *MockTaskService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UserIDKey, int64(7))
		c.Next()
	})
	handler := NewTaskHandler(serviceMock)
	router.GET("/tasks", handler.IndexByList)
	router.POST("/tasks", handler.Store)
	router.GET("/tasks/:id", handler.Show)
	router.PUT("/tasks/:id", handler.Update)
	router.DELETE("/tasks/:id", handler.Delete)
	return router
}

func TestStore_Success(t *testing.T) {
	serviceMock := new(MockTaskService)
	input := types.Task{ListID: 8, Title: "Read a book", Priority: types.PriorityHigh}
	expected := input
	expected.ID = 42
	expected.Status = types.StatusPending
	serviceMock.On("CreateTask", mock.Anything, int64(7), input).Return(expected, nil).Once()
	router := setupTaskHandler(serviceMock)

	request := httptest.NewRequest(http.MethodPost, "/tasks",
		strings.NewReader(`{"list_id":8,"title":"Read a book","priority":"high"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusCreated, response.Code)
	assert.JSONEq(t, `{"id":42,"list_id":8,"title":"Read a book","priority":"high","status":"pending","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`, response.Body.String())
	serviceMock.AssertExpectations(t)
}

func TestStore_BadRequest(t *testing.T) {
	serviceMock := new(MockTaskService)
	router := setupTaskHandler(serviceMock)

	request := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"title":`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	serviceMock.AssertNotCalled(t, "CreateTask", mock.Anything, mock.Anything, mock.Anything)
}

func TestStore_ListNotFound(t *testing.T) {
	serviceMock := new(MockTaskService)
	serviceMock.On("CreateTask", mock.Anything, int64(7), mock.Anything).
		Return(types.Task{}, types.ErrListNotFound).Once()
	router := setupTaskHandler(serviceMock)

	request := httptest.NewRequest(http.MethodPost, "/tasks",
		strings.NewReader(`{"list_id":8,"title":"Read a book","priority":"high"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNotFound, response.Code)
	serviceMock.AssertExpectations(t)
}

func TestIndexByList_MapsTypesToResponse(t *testing.T) {
	serviceMock := new(MockTaskService)
	serviceMock.On("GetTasksByList", mock.Anything, int64(8), int64(7)).
		Return([]types.Task{{ID: 42, ListID: 8, Title: "Read a book", Status: types.StatusPending, Priority: types.PriorityHigh}}, nil).Once()
	router := setupTaskHandler(serviceMock)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/tasks?list_id=8", nil))

	assert.Equal(t, http.StatusOK, response.Code)
	assert.JSONEq(t, `[{"id":42,"list_id":8,"title":"Read a book","status":"pending","priority":"high","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}]`, response.Body.String())
	serviceMock.AssertExpectations(t)
}

func TestShow_NotFound(t *testing.T) {
	serviceMock := new(MockTaskService)
	serviceMock.On("GetTaskByID", mock.Anything, int64(42), int64(7)).
		Return(types.Task{}, types.ErrTaskNotFound).Once()
	router := setupTaskHandler(serviceMock)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/tasks/42", nil))

	assert.Equal(t, http.StatusNotFound, response.Code)
	serviceMock.AssertExpectations(t)
}

func TestUpdate_InternalServerError(t *testing.T) {
	serviceMock := new(MockTaskService)
	serviceMock.On("UpdateTask", mock.Anything, int64(42), int64(7), mock.Anything).
		Return(types.Task{}, errors.New("unexpected failure")).Once()
	router := setupTaskHandler(serviceMock)

	request := httptest.NewRequest(http.MethodPut, "/tasks/42",
		strings.NewReader(`{"title":"Read a book","status":"pending","priority":"high"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	serviceMock.AssertExpectations(t)
}

func TestDelete_Success(t *testing.T) {
	serviceMock := new(MockTaskService)
	serviceMock.On("DeleteTask", mock.Anything, int64(42), int64(7)).Return(nil).Once()
	router := setupTaskHandler(serviceMock)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodDelete, "/tasks/42", nil))

	assert.Equal(t, http.StatusNoContent, response.Code)
	serviceMock.AssertExpectations(t)
}
