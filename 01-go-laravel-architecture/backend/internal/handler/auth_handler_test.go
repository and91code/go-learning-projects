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

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) GetUserByID(ctx context.Context, id int64) (types.User, error) {
	args := m.Called(ctx, id)
	user, _ := args.Get(0).(types.User)
	return user, args.Error(1)
}

func (m *MockAuthService) Register(ctx context.Context, input types.RegisterInput) (types.AuthResult, error) {
	args := m.Called(ctx, input)
	result, _ := args.Get(0).(types.AuthResult)
	return result, args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, input types.LoginInput) (types.AuthResult, error) {
	args := m.Called(ctx, input)
	result, _ := args.Get(0).(types.AuthResult)
	return result, args.Error(1)
}

func (m *MockAuthService) UpdateProfile(ctx context.Context, id int64, input types.UpdateProfileInput) (types.User, error) {
	args := m.Called(ctx, id, input)
	user, _ := args.Get(0).(types.User)
	return user, args.Error(1)
}

func (m *MockAuthService) UpdatePassword(ctx context.Context, id int64, input types.UpdatePasswordInput) error {
	return m.Called(ctx, id, input).Error(0)
}

func setupAuthHandler(serviceMock *MockAuthService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set(middleware.UserIDKey, int64(1))
		c.Next()
	})
	handler := NewAuthHandler(serviceMock)
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)
	router.PUT("/user/profile", handler.UpdateProfile)
	router.PUT("/user/password", handler.UpdatePassword)
	return router
}

func TestAuthUpdateProfile_Success(t *testing.T) {
	serviceMock := new(MockAuthService)
	input := types.UpdateProfileInput{Name: "Ana Silva", Email: "ana@example.com"}
	serviceMock.On("UpdateProfile", mock.Anything, int64(1), input).
		Return(types.User{ID: 1, Name: input.Name, Email: input.Email}, nil).Once()

	response := requestJSON(setupAuthHandler(serviceMock), http.MethodPut, "/user/profile",
		`{"name":"Ana Silva","email":"ana@example.com"}`)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), `"name":"Ana Silva"`)
	serviceMock.AssertExpectations(t)
}

func TestAuthUpdateProfile_InvalidPayload(t *testing.T) {
	serviceMock := new(MockAuthService)

	response := requestJSON(setupAuthHandler(serviceMock), http.MethodPut, "/user/profile", `{"name":"A","email":"invalido"}`)

	assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	serviceMock.AssertNotCalled(t, "UpdateProfile", mock.Anything, mock.Anything, mock.Anything)
}

func TestAuthUpdateProfile_EmailExists(t *testing.T) {
	serviceMock := new(MockAuthService)
	serviceMock.On("UpdateProfile", mock.Anything, int64(1), mock.Anything).
		Return(types.User{}, types.ErrEmailExists).Once()

	response := requestJSON(setupAuthHandler(serviceMock), http.MethodPut, "/user/profile",
		`{"name":"Ana Silva","email":"ana@example.com"}`)

	assert.Equal(t, http.StatusConflict, response.Code)
}

func TestAuthUpdatePassword_Success(t *testing.T) {
	serviceMock := new(MockAuthService)
	input := types.UpdatePasswordInput{CurrentPassword: "secret1", NewPassword: "newsecret"}
	serviceMock.On("UpdatePassword", mock.Anything, int64(1), input).Return(nil).Once()

	response := requestJSON(setupAuthHandler(serviceMock), http.MethodPut, "/user/password",
		`{"current_password":"secret1","new_password":"newsecret"}`)

	assert.Equal(t, http.StatusOK, response.Code)
	serviceMock.AssertExpectations(t)
}

func TestAuthUpdatePassword_InvalidPayload(t *testing.T) {
	serviceMock := new(MockAuthService)

	response := requestJSON(setupAuthHandler(serviceMock), http.MethodPut, "/user/password",
		`{"current_password":"secret1","new_password":"short"}`)

	assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	serviceMock.AssertNotCalled(t, "UpdatePassword", mock.Anything, mock.Anything, mock.Anything)
}

func TestAuthUpdatePassword_WrongCurrentPassword(t *testing.T) {
	serviceMock := new(MockAuthService)
	serviceMock.On("UpdatePassword", mock.Anything, int64(1), mock.Anything).
		Return(types.ErrInvalidCurrentPassword).Once()

	response := requestJSON(setupAuthHandler(serviceMock), http.MethodPut, "/user/password",
		`{"current_password":"wrongpass","new_password":"newsecret"}`)

	assert.Equal(t, http.StatusUnauthorized, response.Code)
}

func postJSON(router *gin.Engine, path, body string) *httptest.ResponseRecorder {
	return requestJSON(router, http.MethodPost, path, body)
}

func requestJSON(router *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

// Cenário: cadastro válido responde 201 com token e dados públicos do usuário.
func TestAuthRegister_Success(t *testing.T) {
	serviceMock := new(MockAuthService)
	input := types.RegisterInput{Name: "Ana", Email: "ana@example.com", Password: "secret1"}
	serviceMock.On("Register", mock.Anything, input).
		Return(types.AuthResult{Token: "jwt", User: types.User{ID: 1, Name: "Ana", Email: "ana@example.com"}}, nil).Once()

	response := postJSON(setupAuthHandler(serviceMock), "/register",
		`{"name":"Ana","email":"ana@example.com","password":"secret1"}`)

	assert.Equal(t, http.StatusCreated, response.Code)
	assert.Contains(t, response.Body.String(), `"token":"jwt"`)
	serviceMock.AssertExpectations(t)
}

// Cenário: payload inválido (e-mail malformado) responde 422 sem chamar o service.
func TestAuthRegister_InvalidPayload(t *testing.T) {
	serviceMock := new(MockAuthService)

	response := postJSON(setupAuthHandler(serviceMock), "/register",
		`{"name":"Ana","email":"invalido","password":"secret1"}`)

	assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	serviceMock.AssertNotCalled(t, "Register", mock.Anything, mock.Anything)
}

// Cenário: e-mail já cadastrado responde 409.
func TestAuthRegister_EmailExists(t *testing.T) {
	serviceMock := new(MockAuthService)
	serviceMock.On("Register", mock.Anything, mock.Anything).Return(types.AuthResult{}, types.ErrEmailExists).Once()

	response := postJSON(setupAuthHandler(serviceMock), "/register",
		`{"name":"Ana","email":"ana@example.com","password":"secret1"}`)

	assert.Equal(t, http.StatusConflict, response.Code)
}

// Cenário: erro inesperado no cadastro responde 500.
func TestAuthRegister_InternalError(t *testing.T) {
	serviceMock := new(MockAuthService)
	serviceMock.On("Register", mock.Anything, mock.Anything).Return(types.AuthResult{}, errors.New("boom")).Once()

	response := postJSON(setupAuthHandler(serviceMock), "/register",
		`{"name":"Ana","email":"ana@example.com","password":"secret1"}`)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
}

// Cenário: login válido responde 200 com token.
func TestAuthLogin_Success(t *testing.T) {
	serviceMock := new(MockAuthService)
	serviceMock.On("Login", mock.Anything, types.LoginInput{Email: "ana@example.com", Password: "secret1"}).
		Return(types.AuthResult{Token: "jwt"}, nil).Once()

	response := postJSON(setupAuthHandler(serviceMock), "/login", `{"email":"ana@example.com","password":"secret1"}`)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), `"token":"jwt"`)
}

// Cenário: payload sem senha responde 422 sem chamar o service.
func TestAuthLogin_InvalidPayload(t *testing.T) {
	serviceMock := new(MockAuthService)

	response := postJSON(setupAuthHandler(serviceMock), "/login", `{"email":"ana@example.com"}`)

	assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	serviceMock.AssertNotCalled(t, "Login", mock.Anything, mock.Anything)
}

// Cenário: credenciais inválidas respondem 401.
func TestAuthLogin_InvalidCredentials(t *testing.T) {
	serviceMock := new(MockAuthService)
	serviceMock.On("Login", mock.Anything, mock.Anything).Return(types.AuthResult{}, types.ErrInvalidCredentials).Once()

	response := postJSON(setupAuthHandler(serviceMock), "/login", `{"email":"ana@example.com","password":"x"}`)

	assert.Equal(t, http.StatusUnauthorized, response.Code)
}

// Cenário: erro inesperado no login responde 500.
func TestAuthLogin_InternalError(t *testing.T) {
	serviceMock := new(MockAuthService)
	serviceMock.On("Login", mock.Anything, mock.Anything).Return(types.AuthResult{}, errors.New("boom")).Once()

	response := postJSON(setupAuthHandler(serviceMock), "/login", `{"email":"ana@example.com","password":"x"}`)

	assert.Equal(t, http.StatusInternalServerError, response.Code)
}
