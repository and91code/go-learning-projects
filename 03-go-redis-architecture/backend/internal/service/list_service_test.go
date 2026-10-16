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

func newListService(repo *MockListRepository) types.ListService {
	return NewListService(repo)
}

// Cenário: CreateList associa a lista ao usuário informado e devolve o ID gerado.
func TestCreateList_Success(t *testing.T) {
	repo := new(MockListRepository)
	repo.On("Create", mock.Anything, int64(7), mock.MatchedBy(func(list *types.List) bool {
		return list.UserID == 7 && list.Title == "Casa"
	})).Run(func(args mock.Arguments) {
		args.Get(2).(*types.List).ID = 3
	}).Return(nil).Once()

	list, err := newListService(repo).CreateList(context.Background(), 7, types.List{Title: "Casa"})

	assert.NoError(t, err)
	assert.Equal(t, int64(3), list.ID)
	assert.Equal(t, int64(7), list.UserID)
	repo.AssertExpectations(t)
}

// Cenário: CreateList propaga erro do repositório e não devolve lista parcial.
func TestCreateList_RepositoryError(t *testing.T) {
	dbErr := errors.New("database unavailable")
	repo := new(MockListRepository)
	repo.On("Create", mock.Anything, int64(7), mock.Anything).Return(dbErr).Once()

	list, err := newListService(repo).CreateList(context.Background(), 7, types.List{Title: "Casa"})

	assert.ErrorIs(t, err, dbErr)
	assert.Empty(t, list)
}

// Cenário: GetListByID devolve a lista do usuário.
func TestGetListByID_Success(t *testing.T) {
	repo := new(MockListRepository)
	repo.On("FindByIDAndUserID", mock.Anything, int64(3), int64(7)).
		Return(&types.List{ID: 3, UserID: 7}, nil).Once()

	list, err := newListService(repo).GetListByID(context.Background(), 3, 7)

	assert.NoError(t, err)
	assert.Equal(t, int64(3), list.ID)
}

// Cenário: GetListByID propaga lista não encontrada.
func TestGetListByID_NotFound(t *testing.T) {
	repo := new(MockListRepository)
	repo.On("FindByIDAndUserID", mock.Anything, int64(3), int64(7)).Return(nil, types.ErrListNotFound).Once()

	list, err := newListService(repo).GetListByID(context.Background(), 3, 7)

	assert.ErrorIs(t, err, types.ErrListNotFound)
	assert.Empty(t, list)
}

// Cenário: GetUserLists devolve as listas do usuário.
func TestGetUserLists_Success(t *testing.T) {
	repo := new(MockListRepository)
	repo.On("FindByUserID", mock.Anything, int64(7)).Return([]types.List{{ID: 1}, {ID: 2}}, nil).Once()

	lists, err := newListService(repo).GetUserLists(context.Background(), 7)

	assert.NoError(t, err)
	assert.Len(t, lists, 2)
}

// Cenário: UpdateList altera título e descrição e envia id e userID ao repositório.
func TestUpdateList_Success(t *testing.T) {
	repo := new(MockListRepository)
	repo.On("FindByIDAndUserID", mock.Anything, int64(3), int64(7)).
		Return(&types.List{ID: 3, UserID: 7, Title: "Antigo"}, nil).Once()
	repo.On("Update", mock.Anything, int64(3), int64(7), mock.MatchedBy(func(list *types.List) bool {
		return list.Title == "Novo" && list.Description == "Desc"
	})).Return(nil).Once()

	list, err := newListService(repo).UpdateList(context.Background(), 3, 7,
		types.List{Title: "Novo", Description: "Desc"})

	assert.NoError(t, err)
	assert.Equal(t, "Novo", list.Title)
	repo.AssertExpectations(t)
}

// Cenário: UpdateList não atualiza quando a lista não pertence ao usuário.
func TestUpdateList_NotFound(t *testing.T) {
	repo := new(MockListRepository)
	repo.On("FindByIDAndUserID", mock.Anything, int64(3), int64(7)).Return(nil, types.ErrListNotFound).Once()

	_, err := newListService(repo).UpdateList(context.Background(), 3, 7, types.List{Title: "Novo"})

	assert.ErrorIs(t, err, types.ErrListNotFound)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

// Cenário: UpdateList propaga erro ao gravar a atualização.
func TestUpdateList_RepositoryError(t *testing.T) {
	dbErr := errors.New("database unavailable")
	repo := new(MockListRepository)
	repo.On("FindByIDAndUserID", mock.Anything, int64(3), int64(7)).Return(&types.List{ID: 3}, nil).Once()
	repo.On("Update", mock.Anything, int64(3), int64(7), mock.Anything).Return(dbErr).Once()

	list, err := newListService(repo).UpdateList(context.Background(), 3, 7, types.List{Title: "Novo"})

	assert.ErrorIs(t, err, dbErr)
	assert.Empty(t, list)
}

// Cenário: DeleteList repassa id e userID ao repositório.
func TestDeleteList_Success(t *testing.T) {
	repo := new(MockListRepository)
	repo.On("Delete", mock.Anything, int64(3), int64(7)).Return(nil).Once()

	err := newListService(repo).DeleteList(context.Background(), 3, 7)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
