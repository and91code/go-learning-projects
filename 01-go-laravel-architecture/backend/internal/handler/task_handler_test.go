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
	"github.com/seu-usuario/taskflow-backend/internal/dto"
	"github.com/seu-usuario/taskflow-backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) CreateTask(ctx context.Context, request dto.CreateTaskRequest) (dto.TaskResponse, error) {
	args := m.Called(ctx, request)
	result, _ := args.Get(0).(dto.TaskResponse)
	return result, args.Error(1)
}

func setupTaskHandler(serviceMock *MockTaskService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/tasks", NewTaskHandler(serviceMock).Store)
	return router
}

func TestStore_Success(t *testing.T) {
	serviceMock := new(MockTaskService)
	expected := dto.TaskResponse{ID: 42, Title: "Read a book", Priority: "high", Status: "pending"}
	serviceMock.On("CreateTask", mock.Anything, dto.CreateTaskRequest{
		Title: "Read a book", Priority: "high",
	}).Return(expected, nil).Once()
	router := setupTaskHandler(serviceMock)

	request := httptest.NewRequest(http.MethodPost, "/tasks",
		strings.NewReader(`{"title":"Read a book","priority":"high"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusCreated, response.Code)
	assert.JSONEq(t, `{"id":42,"title":"Read a book","priority":"high","status":"pending"}`, response.Body.String())
	serviceMock.AssertExpectations(t)
}

func TestStore_BadRequest(t *testing.T) {
	serviceMock := new(MockTaskService)
	router := setupTaskHandler(serviceMock)

	request := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"title":`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	serviceMock.AssertNotCalled(t, "CreateTask", mock.Anything, mock.Anything)
}

func TestStore_UnprocessableEntity(t *testing.T) {
	serviceMock := new(MockTaskService)
	serviceMock.On("CreateTask", mock.Anything, dto.CreateTaskRequest{
		Title: "Read a book", Priority: "high",
	}).Return(dto.TaskResponse{}, service.ErrInvalidTask).Once()
	router := setupTaskHandler(serviceMock)

	request := httptest.NewRequest(http.MethodPost, "/tasks",
		strings.NewReader(`{"title":"Read a book","priority":"high"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	serviceMock.AssertExpectations(t)
}

func TestStore_InternalServerError(t *testing.T) {
	serviceMock := new(MockTaskService)
	serviceMock.On("CreateTask", mock.Anything, dto.CreateTaskRequest{
		Title: "Read a book", Priority: "high",
	}).Return(dto.TaskResponse{}, errors.New("unexpected failure")).Once()
	router := setupTaskHandler(serviceMock)

	request := httptest.NewRequest(http.MethodPost, "/tasks",
		strings.NewReader(`{"title":"Read a book","priority":"high"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
	serviceMock.AssertExpectations(t)
}
