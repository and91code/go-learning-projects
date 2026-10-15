//go:build unit

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/seu-usuario/taskflow-backend/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *types.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*types.User, error) {
	args := m.Called(ctx, email)
	user, _ := args.Get(0).(*types.User)
	return user, args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id int64) (*types.User, error) {
	args := m.Called(ctx, id)
	user, _ := args.Get(0).(*types.User)
	return user, args.Error(1)
}

func (m *MockUserRepository) UpdateProfile(ctx context.Context, id int64, name, email string) error {
	return m.Called(ctx, id, name, email).Error(0)
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, id int64, passwordHash string) error {
	return m.Called(ctx, id, passwordHash).Error(0)
}

// Cenário: GetUserByID devolve o usuário encontrado.
func TestGetUserByID_Success(t *testing.T) {
	repo := new(MockUserRepository)
	repo.On("FindByID", mock.Anything, int64(1)).Return(&types.User{ID: 1, Name: "Ana"}, nil).Once()

	user, err := NewAuthService(repo).GetUserByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, "Ana", user.Name)
	repo.AssertExpectations(t)
}

// Cenário: GetUserByID propaga o erro do repositório (ex.: usuário inexistente).
func TestGetUserByID_NotFound(t *testing.T) {
	repo := new(MockUserRepository)
	repo.On("FindByID", mock.Anything, int64(1)).Return(nil, types.ErrUserNotFound).Once()

	user, err := NewAuthService(repo).GetUserByID(context.Background(), 1)

	assert.ErrorIs(t, err, types.ErrUserNotFound)
	assert.Empty(t, user)
}

func TestUpdateProfile_Success(t *testing.T) {
	repo := new(MockUserRepository)
	repo.On("FindByEmail", mock.Anything, "ana@example.com").Return(nil, types.ErrUserNotFound).Once()
	repo.On("UpdateProfile", mock.Anything, int64(1), "Ana Silva", "ana@example.com").Return(nil).Once()
	repo.On("FindByID", mock.Anything, int64(1)).
		Return(&types.User{ID: 1, Name: "Ana Silva", Email: "ana@example.com"}, nil).Once()

	user, err := NewAuthService(repo).UpdateProfile(context.Background(), 1,
		types.UpdateProfileInput{Name: "Ana Silva", Email: "ana@example.com"})

	assert.NoError(t, err)
	assert.Equal(t, "Ana Silva", user.Name)
	repo.AssertExpectations(t)
}

func TestUpdateProfile_EmailAlreadyExists(t *testing.T) {
	repo := new(MockUserRepository)
	repo.On("FindByEmail", mock.Anything, "other@example.com").
		Return(&types.User{ID: 2}, nil).Once()

	_, err := NewAuthService(repo).UpdateProfile(context.Background(), 1,
		types.UpdateProfileInput{Name: "Ana", Email: "other@example.com"})

	assert.ErrorIs(t, err, types.ErrEmailExists)
	repo.AssertNotCalled(t, "UpdateProfile", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestUpdatePassword_Success(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret1"), bcrypt.MinCost)
	repo := new(MockUserRepository)
	repo.On("FindByID", mock.Anything, int64(1)).
		Return(&types.User{ID: 1, PasswordHash: string(hash)}, nil).Once()
	repo.On("UpdatePassword", mock.Anything, int64(1), mock.MatchedBy(func(passwordHash string) bool {
		return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte("newsecret")) == nil
	})).Return(nil).Once()

	err := NewAuthService(repo).UpdatePassword(context.Background(), 1,
		types.UpdatePasswordInput{CurrentPassword: "secret1", NewPassword: "newsecret"})

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestUpdatePassword_WrongCurrentPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret1"), bcrypt.MinCost)
	repo := new(MockUserRepository)
	repo.On("FindByID", mock.Anything, int64(1)).
		Return(&types.User{ID: 1, PasswordHash: string(hash)}, nil).Once()

	err := NewAuthService(repo).UpdatePassword(context.Background(), 1,
		types.UpdatePasswordInput{CurrentPassword: "wrong", NewPassword: "newsecret"})

	assert.ErrorIs(t, err, types.ErrInvalidCurrentPassword)
	repo.AssertNotCalled(t, "UpdatePassword", mock.Anything, mock.Anything, mock.Anything)
}

// Cenário: Register cria o usuário com senha criptografada e devolve um token JWT.
func TestRegister_Success(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	repo := new(MockUserRepository)
	repo.On("FindByEmail", mock.Anything, "ana@example.com").Return(nil, types.ErrUserNotFound).Once()
	repo.On("Create", mock.Anything, mock.MatchedBy(func(user *types.User) bool {
		return user.Email == "ana@example.com" && user.PasswordHash != "secret1" &&
			bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("secret1")) == nil
	})).Run(func(args mock.Arguments) {
		args.Get(1).(*types.User).ID = 1
	}).Return(nil).Once()

	result, err := NewAuthService(repo).Register(context.Background(),
		types.RegisterInput{Name: "Ana", Email: "ana@example.com", Password: "secret1"})

	assert.NoError(t, err)
	assert.NotEmpty(t, result.Token)
	assert.Equal(t, int64(1), result.User.ID)
	repo.AssertExpectations(t)
}

// Cenário: Register rejeita e-mail já cadastrado sem criar usuário.
func TestRegister_EmailAlreadyExists(t *testing.T) {
	repo := new(MockUserRepository)
	repo.On("FindByEmail", mock.Anything, "ana@example.com").Return(&types.User{ID: 1}, nil).Once()

	_, err := NewAuthService(repo).Register(context.Background(),
		types.RegisterInput{Name: "Ana", Email: "ana@example.com", Password: "secret1"})

	assert.ErrorIs(t, err, types.ErrEmailExists)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// Cenário: Register propaga falha inesperada ao verificar o e-mail.
func TestRegister_FindByEmailError(t *testing.T) {
	dbErr := errors.New("database unavailable")
	repo := new(MockUserRepository)
	repo.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, dbErr).Once()

	_, err := NewAuthService(repo).Register(context.Background(),
		types.RegisterInput{Email: "ana@example.com", Password: "secret1"})

	assert.ErrorIs(t, err, dbErr)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// Cenário: Register propaga falha ao gravar o usuário.
func TestRegister_CreateError(t *testing.T) {
	dbErr := errors.New("database unavailable")
	repo := new(MockUserRepository)
	repo.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, types.ErrUserNotFound).Once()
	repo.On("Create", mock.Anything, mock.Anything).Return(dbErr).Once()

	_, err := NewAuthService(repo).Register(context.Background(),
		types.RegisterInput{Email: "ana@example.com", Password: "secret1"})

	assert.ErrorIs(t, err, dbErr)
}

// Cenário: Register falha quando JWT_SECRET não está configurado.
func TestRegister_MissingJWTSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "")
	repo := new(MockUserRepository)
	repo.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, types.ErrUserNotFound).Once()
	repo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

	result, err := NewAuthService(repo).Register(context.Background(),
		types.RegisterInput{Email: "ana@example.com", Password: "secret1"})

	assert.ErrorContains(t, err, "JWT_SECRET")
	assert.Empty(t, result)
}

// Cenário: Login com credenciais corretas devolve token e usuário.
func TestLogin_Success(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret1"), bcrypt.MinCost)
	repo := new(MockUserRepository)
	repo.On("FindByEmail", mock.Anything, "ana@example.com").
		Return(&types.User{ID: 1, Email: "ana@example.com", PasswordHash: string(hash)}, nil).Once()

	result, err := NewAuthService(repo).Login(context.Background(),
		types.LoginInput{Email: "ana@example.com", Password: "secret1"})

	assert.NoError(t, err)
	assert.NotEmpty(t, result.Token)
	assert.Equal(t, int64(1), result.User.ID)
}

// Cenário: Login com e-mail inexistente devolve credenciais inválidas (sem revelar que o usuário não existe).
func TestLogin_UnknownEmail(t *testing.T) {
	repo := new(MockUserRepository)
	repo.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, types.ErrUserNotFound).Once()

	_, err := NewAuthService(repo).Login(context.Background(),
		types.LoginInput{Email: "x@example.com", Password: "secret1"})

	assert.ErrorIs(t, err, types.ErrInvalidCredentials)
}

// Cenário: Login com senha errada devolve credenciais inválidas.
func TestLogin_WrongPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret1"), bcrypt.MinCost)
	repo := new(MockUserRepository)
	repo.On("FindByEmail", mock.Anything, mock.Anything).
		Return(&types.User{ID: 1, PasswordHash: string(hash)}, nil).Once()

	_, err := NewAuthService(repo).Login(context.Background(),
		types.LoginInput{Email: "ana@example.com", Password: "wrong"})

	assert.ErrorIs(t, err, types.ErrInvalidCredentials)
}

// Cenário: Login propaga falha inesperada do repositório.
func TestLogin_RepositoryError(t *testing.T) {
	dbErr := errors.New("database unavailable")
	repo := new(MockUserRepository)
	repo.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, dbErr).Once()

	_, err := NewAuthService(repo).Login(context.Background(),
		types.LoginInput{Email: "ana@example.com", Password: "secret1"})

	assert.ErrorIs(t, err, dbErr)
}
