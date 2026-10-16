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

type MockTagService struct {
	mock.Mock
}

func (m *MockTagService) CreateTag(ctx context.Context, userID int64, tag types.Tag) (types.Tag, error) {
	args := m.Called(ctx, userID, tag)
	result, _ := args.Get(0).(types.Tag)
	return result, args.Error(1)
}

func (m *MockTagService) UpdateTag(ctx context.Context, tagID, userID int64, tag types.Tag) (types.Tag, error) {
	args := m.Called(ctx, tagID, userID, tag)
	result, _ := args.Get(0).(types.Tag)
	return result, args.Error(1)
}

func (m *MockTagService) DeleteTag(ctx context.Context, tagID, userID int64) error {
	return m.Called(ctx, tagID, userID).Error(0)
}

func (m *MockTagService) GetTagsByList(ctx context.Context, listID, userID int64) ([]types.Tag, error) {
	args := m.Called(ctx, listID, userID)
	result, _ := args.Get(0).([]types.Tag)
	return result, args.Error(1)
}

func (m *MockTagService) GetTagById(ctx context.Context, tagID, userID int64) (types.Tag, error) {
	args := m.Called(ctx, tagID, userID)
	result, _ := args.Get(0).(types.Tag)
	return result, args.Error(1)
}

// setupTagHandler simula o middleware de autenticação com o usuário 7; withUser=false omite o userID.
func setupTagHandler(serviceMock *MockTagService, withUser bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	if withUser {
		router.Use(func(c *gin.Context) {
			c.Set(middleware.UserIDKey, int64(7))
			c.Next()
		})
	}
	handler := NewTagHandler(serviceMock)
	router.GET("/tags", handler.IndexByList)
	router.POST("/tags", handler.Store)
	router.GET("/tags/:id", handler.Show)
	router.PUT("/tags/:id", handler.Update)
	router.DELETE("/tags/:id", handler.Delete)
	return router
}

func putJSON(router *gin.Engine, path, body string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodPut, path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func doRequest(router *gin.Engine, method, path string) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(method, path, nil))
	return response
}

// Cenário: criação válida responde 201 com os dados detalhados da tag.
func TestTagStore_Success(t *testing.T) {
	serviceMock := new(MockTagService)
	serviceMock.On("CreateTag", mock.Anything, int64(7), types.Tag{ListID: 8, Name: "urgente", Color: "#ff0000"}).
		Return(types.Tag{ID: 5, ListID: 8, Name: "urgente", Color: "#ff0000"}, nil).Once()

	response := postJSON(setupTagHandler(serviceMock, true), "/tags", `{"list_id":8,"name":"urgente","color":"#ff0000"}`)

	assert.Equal(t, http.StatusCreated, response.Code)
	assert.JSONEq(t, `{"id":5,"list_id":8,"name":"urgente","color":"#ff0000"}`, response.Body.String())
}

// Cenário: requisição sem usuário autenticado responde 401 sem chamar o service.
func TestTagStore_Unauthorized(t *testing.T) {
	serviceMock := new(MockTagService)

	response := postJSON(setupTagHandler(serviceMock, false), "/tags", `{"list_id":8,"name":"urgente"}`)

	assert.Equal(t, http.StatusUnauthorized, response.Code)
	serviceMock.AssertNotCalled(t, "CreateTag", mock.Anything, mock.Anything, mock.Anything)
}

// Cenário: payload sem nome responde 422 sem chamar o service.
func TestTagStore_InvalidPayload(t *testing.T) {
	serviceMock := new(MockTagService)

	response := postJSON(setupTagHandler(serviceMock, true), "/tags", `{"list_id":8}`)

	assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	serviceMock.AssertNotCalled(t, "CreateTag", mock.Anything, mock.Anything, mock.Anything)
}

// Cenário: erros do service são traduzidos em status HTTP (404, 422 e 500).
func TestTagStore_ServiceErrors(t *testing.T) {
	tests := map[string]struct {
		err    error
		status int
	}{
		"lista não encontrada": {err: types.ErrListNotFound, status: http.StatusNotFound},
		"tag inválida":         {err: types.ErrInvalidTag, status: http.StatusUnprocessableEntity},
		"erro inesperado":      {err: errors.New("boom"), status: http.StatusInternalServerError},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			serviceMock := new(MockTagService)
			serviceMock.On("CreateTag", mock.Anything, int64(7), mock.Anything).Return(types.Tag{}, test.err).Once()

			response := postJSON(setupTagHandler(serviceMock, true), "/tags", `{"list_id":8,"name":"urgente"}`)

			assert.Equal(t, test.status, response.Code)
		})
	}
}

// Cenário: listagem por list_id responde 200 com as tags resumidas.
func TestTagIndexByList_Success(t *testing.T) {
	serviceMock := new(MockTagService)
	serviceMock.On("GetTagsByList", mock.Anything, int64(8), int64(7)).
		Return([]types.Tag{{ID: 1, Name: "a", Color: "#fff"}}, nil).Once()

	response := doRequest(setupTagHandler(serviceMock, true), http.MethodGet, "/tags?list_id=8")

	assert.Equal(t, http.StatusOK, response.Code)
	assert.JSONEq(t, `[{"id":1,"name":"a","color":"#fff"}]`, response.Body.String())
}

// Cenário: list_id ausente ou inválido responde 400.
func TestTagIndexByList_InvalidListID(t *testing.T) {
	serviceMock := new(MockTagService)

	response := doRequest(setupTagHandler(serviceMock, true), http.MethodGet, "/tags")

	assert.Equal(t, http.StatusBadRequest, response.Code)
	serviceMock.AssertNotCalled(t, "GetTagsByList", mock.Anything, mock.Anything, mock.Anything)
}

// Cenário: lista inexistente na listagem responde 404.
func TestTagIndexByList_ListNotFound(t *testing.T) {
	serviceMock := new(MockTagService)
	serviceMock.On("GetTagsByList", mock.Anything, int64(8), int64(7)).Return(nil, types.ErrListNotFound).Once()

	response := doRequest(setupTagHandler(serviceMock, true), http.MethodGet, "/tags?list_id=8")

	assert.Equal(t, http.StatusNotFound, response.Code)
}

// Cenário: consulta por ID responde 200 com os dados detalhados.
func TestTagShow_Success(t *testing.T) {
	serviceMock := new(MockTagService)
	serviceMock.On("GetTagById", mock.Anything, int64(5), int64(7)).
		Return(types.Tag{ID: 5, ListID: 8, Name: "a"}, nil).Once()

	response := doRequest(setupTagHandler(serviceMock, true), http.MethodGet, "/tags/5")

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), `"list_id":8`)
}

// Cenário: ID inválido responde 400 sem chamar o service.
func TestTagShow_InvalidID(t *testing.T) {
	serviceMock := new(MockTagService)

	response := doRequest(setupTagHandler(serviceMock, true), http.MethodGet, "/tags/abc")

	assert.Equal(t, http.StatusBadRequest, response.Code)
	serviceMock.AssertNotCalled(t, "GetTagById", mock.Anything, mock.Anything, mock.Anything)
}

// Cenário: tag inexistente responde 404.
func TestTagShow_NotFound(t *testing.T) {
	serviceMock := new(MockTagService)
	serviceMock.On("GetTagById", mock.Anything, int64(5), int64(7)).Return(types.Tag{}, types.ErrTagNotFound).Once()

	response := doRequest(setupTagHandler(serviceMock, true), http.MethodGet, "/tags/5")

	assert.Equal(t, http.StatusNotFound, response.Code)
}

// Cenário: atualização válida responde 200 e envia tagID e userID ao service.
func TestTagUpdate_Success(t *testing.T) {
	serviceMock := new(MockTagService)
	serviceMock.On("UpdateTag", mock.Anything, int64(5), int64(7), types.Tag{ListID: 8, Name: "novo"}).
		Return(types.Tag{ID: 5, ListID: 8, Name: "novo"}, nil).Once()

	response := putJSON(setupTagHandler(serviceMock, true), "/tags/5", `{"list_id":8,"name":"novo"}`)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), `"name":"novo"`)
}

// Cenário: ID inválido na atualização responde 400.
func TestTagUpdate_InvalidID(t *testing.T) {
	serviceMock := new(MockTagService)

	response := putJSON(setupTagHandler(serviceMock, true), "/tags/0", `{"list_id":8,"name":"novo"}`)

	assert.Equal(t, http.StatusBadRequest, response.Code)
}

// Cenário: payload inválido na atualização responde 422.
func TestTagUpdate_InvalidPayload(t *testing.T) {
	serviceMock := new(MockTagService)

	response := putJSON(setupTagHandler(serviceMock, true), "/tags/5", `{"name":"novo"}`)

	assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	serviceMock.AssertNotCalled(t, "UpdateTag", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

// Cenário: tag inexistente na atualização responde 404.
func TestTagUpdate_NotFound(t *testing.T) {
	serviceMock := new(MockTagService)
	serviceMock.On("UpdateTag", mock.Anything, int64(5), int64(7), mock.Anything).
		Return(types.Tag{}, types.ErrTagNotFound).Once()

	response := putJSON(setupTagHandler(serviceMock, true), "/tags/5", `{"list_id":8,"name":"novo"}`)

	assert.Equal(t, http.StatusNotFound, response.Code)
}

// Cenário: exclusão válida responde 204.
func TestTagDelete_Success(t *testing.T) {
	serviceMock := new(MockTagService)
	serviceMock.On("DeleteTag", mock.Anything, int64(5), int64(7)).Return(nil).Once()

	response := doRequest(setupTagHandler(serviceMock, true), http.MethodDelete, "/tags/5")

	assert.Equal(t, http.StatusNoContent, response.Code)
}

// Cenário: ID inválido na exclusão responde 400.
func TestTagDelete_InvalidID(t *testing.T) {
	serviceMock := new(MockTagService)

	response := doRequest(setupTagHandler(serviceMock, true), http.MethodDelete, "/tags/abc")

	assert.Equal(t, http.StatusBadRequest, response.Code)
	serviceMock.AssertNotCalled(t, "DeleteTag", mock.Anything, mock.Anything, mock.Anything)
}

// Cenário: tag inexistente na exclusão responde 404.
func TestTagDelete_NotFound(t *testing.T) {
	serviceMock := new(MockTagService)
	serviceMock.On("DeleteTag", mock.Anything, int64(5), int64(7)).Return(types.ErrTagNotFound).Once()

	response := doRequest(setupTagHandler(serviceMock, true), http.MethodDelete, "/tags/5")

	assert.Equal(t, http.StatusNotFound, response.Code)
}
