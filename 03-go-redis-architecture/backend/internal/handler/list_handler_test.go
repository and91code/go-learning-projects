//go:build unit

package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/seu-usuario/taskflow-backend/internal/middleware"
	"github.com/seu-usuario/taskflow-backend/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockListService struct {
	mock.Mock
}

func (m *MockListService) CreateList(ctx context.Context, userID int64, list types.List) (types.List, error) {
	args := m.Called(ctx, userID, list)
	result, _ := args.Get(0).(types.List)
	return result, args.Error(1)
}

func (m *MockListService) GetListByID(ctx context.Context, id, userID int64) (types.List, error) {
	args := m.Called(ctx, id, userID)
	result, _ := args.Get(0).(types.List)
	return result, args.Error(1)
}

func (m *MockListService) GetUserLists(ctx context.Context, userID int64) ([]types.List, error) {
	args := m.Called(ctx, userID)
	result, _ := args.Get(0).([]types.List)
	return result, args.Error(1)
}

func (m *MockListService) UpdateList(ctx context.Context, id, userID int64, list types.List) (types.List, error) {
	args := m.Called(ctx, id, userID, list)
	result, _ := args.Get(0).(types.List)
	return result, args.Error(1)
}

func (m *MockListService) DeleteList(ctx context.Context, id, userID int64) error {
	return m.Called(ctx, id, userID).Error(0)
}

// setupListHandler simula o middleware de autenticação com o usuário 7.
func setupListHandler(serviceMock *MockListService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UserIDKey, int64(7))
		c.Next()
	})
	handler := NewListHandler(serviceMock)
	router.POST("/lists", handler.Store)
	router.GET("/lists", handler.Index)
	router.GET("/lists/:id", handler.Show)
	return router
}

// Cenário: criação válida responde 201 e envia o userID do contexto ao service.
func TestListStore_Success(t *testing.T) {
	serviceMock := new(MockListService)
	serviceMock.On("CreateList", mock.Anything, int64(7), types.List{Title: "Casa", Description: "Domésticas"}).
		Return(types.List{ID: 3, UserID: 7, Title: "Casa"}, nil).Once()

	response := postJSON(setupListHandler(serviceMock), "/lists", `{"title":"Casa","description":"Domésticas"}`)

	assert.Equal(t, http.StatusCreated, response.Code)
	assert.Contains(t, response.Body.String(), `"id":3`)
	serviceMock.AssertExpectations(t)
}

// Cenário: título curto demais responde 422 sem chamar o service.
func TestListStore_InvalidPayload(t *testing.T) {
	serviceMock := new(MockListService)

	response := postJSON(setupListHandler(serviceMock), "/lists", `{"title":"ab"}`)

	assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	serviceMock.AssertNotCalled(t, "CreateList", mock.Anything, mock.Anything, mock.Anything)
}

// Cenário: erro do service na criação responde 500.
func TestListStore_InternalError(t *testing.T) {
	serviceMock := new(MockListService)
	serviceMock.On("CreateList", mock.Anything, int64(7), mock.Anything).Return(types.List{}, errors.New("boom")).Once()

	response := postJSON(setupListHandler(serviceMock), "/lists", `{"title":"Casa"}`)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
}

// Cenário: listagem responde 200 com as listas do usuário.
func TestListIndex_Success(t *testing.T) {
	serviceMock := new(MockListService)
	serviceMock.On("GetUserLists", mock.Anything, int64(7)).Return([]types.List{{ID: 1}, {ID: 2}}, nil).Once()

	response := httptest.NewRecorder()
	setupListHandler(serviceMock).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/lists", nil))

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), `"id":2`)
}

// Cenário: erro do service na listagem responde 500.
func TestListIndex_InternalError(t *testing.T) {
	serviceMock := new(MockListService)
	serviceMock.On("GetUserLists", mock.Anything, int64(7)).Return(nil, errors.New("boom")).Once()

	response := httptest.NewRecorder()
	setupListHandler(serviceMock).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/lists", nil))

	assert.Equal(t, http.StatusInternalServerError, response.Code)
}

// Cenário: consulta por ID existente responde 200.
func TestListShow_Success(t *testing.T) {
	serviceMock := new(MockListService)
	serviceMock.On("GetListByID", mock.Anything, int64(3), int64(7)).Return(types.List{ID: 3, Title: "Casa"}, nil).Once()

	response := httptest.NewRecorder()
	setupListHandler(serviceMock).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/lists/3", nil))

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), `"title":"Casa"`)
}

// Cenário: ID não numérico responde 400 sem chamar o service.
func TestListShow_InvalidID(t *testing.T) {
	serviceMock := new(MockListService)

	response := httptest.NewRecorder()
	setupListHandler(serviceMock).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/lists/abc", nil))

	assert.Equal(t, http.StatusBadRequest, response.Code)
	serviceMock.AssertNotCalled(t, "GetListByID", mock.Anything, mock.Anything, mock.Anything)
}

// Cenário: lista inexistente (ou de outro usuário) responde 404.
func TestListShow_NotFound(t *testing.T) {
	serviceMock := new(MockListService)
	serviceMock.On("GetListByID", mock.Anything, int64(3), int64(7)).Return(types.List{}, types.ErrListNotFound).Once()

	response := httptest.NewRecorder()
	setupListHandler(serviceMock).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/lists/3", nil))

	assert.Equal(t, http.StatusNotFound, response.Code)
}

// Cenário: erro inesperado na consulta responde 500.
func TestListShow_InternalError(t *testing.T) {
	serviceMock := new(MockListService)
	serviceMock.On("GetListByID", mock.Anything, int64(3), int64(7)).Return(types.List{}, errors.New("boom")).Once()

	response := httptest.NewRecorder()
	setupListHandler(serviceMock).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/lists/3", nil))

	assert.Equal(t, http.StatusInternalServerError, response.Code)
}
