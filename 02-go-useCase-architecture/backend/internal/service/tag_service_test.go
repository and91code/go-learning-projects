//go:build unit

package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/seu-usuario/taskflow-backend/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTagRepository struct {
	mock.Mock
}

func (m *MockTagRepository) Create(ctx context.Context, userID int64, tag *types.Tag) error {
	return m.Called(ctx, userID, tag).Error(0)
}

func (m *MockTagRepository) Update(ctx context.Context, tagID, userID int64, tag *types.Tag) error {
	return m.Called(ctx, tagID, userID, tag).Error(0)
}

func (m *MockTagRepository) Delete(ctx context.Context, tagID, userID int64) error {
	return m.Called(ctx, tagID, userID).Error(0)
}

func (m *MockTagRepository) FindByListIDAndUserID(ctx context.Context, listID, userID int64) ([]types.Tag, error) {
	args := m.Called(ctx, listID, userID)
	tags, _ := args.Get(0).([]types.Tag)
	return tags, args.Error(1)
}

func (m *MockTagRepository) FindByIDAndUserID(ctx context.Context, tagID, userID int64) (types.Tag, error) {
	args := m.Called(ctx, tagID, userID)
	tag, _ := args.Get(0).(types.Tag)
	return tag, args.Error(1)
}

func newTagService(tagRepo *MockTagRepository, listRepo *MockListRepository) types.TagService {
	return NewTagService(tagRepo, listRepo)
}

// Cenário: CreateTag valida a lista do usuário, remove espaços do nome e devolve o ID gerado.
func TestCreateTag_Success(t *testing.T) {
	tagRepo := new(MockTagRepository)
	listRepo := new(MockListRepository)
	listRepo.On("FindByIDAndUserID", mock.Anything, int64(8), int64(7)).Return(&types.List{ID: 8}, nil).Once()
	tagRepo.On("Create", mock.Anything, int64(7), mock.MatchedBy(func(tag *types.Tag) bool {
		return tag.ListID == 8 && tag.Name == "urgente"
	})).Run(func(args mock.Arguments) {
		args.Get(2).(*types.Tag).ID = 5
	}).Return(nil).Once()

	tag, err := newTagService(tagRepo, listRepo).CreateTag(context.Background(), 7,
		types.Tag{ListID: 8, Name: "  urgente  ", Color: "#ff0000"})

	assert.NoError(t, err)
	assert.Equal(t, int64(5), tag.ID)
	tagRepo.AssertExpectations(t)
	listRepo.AssertExpectations(t)
}

// Cenário: CreateTag rejeita dados inválidos sem consultar repositórios.
func TestCreateTag_InvalidInput(t *testing.T) {
	tests := map[string]struct {
		userID int64
		tag    types.Tag
	}{
		"usuário inválido": {userID: 0, tag: types.Tag{ListID: 8, Name: "a"}},
		"lista inválida":   {userID: 7, tag: types.Tag{ListID: 0, Name: "a"}},
		"nome vazio":       {userID: 7, tag: types.Tag{ListID: 8, Name: "   "}},
		"nome muito longo": {userID: 7, tag: types.Tag{ListID: 8, Name: strings.Repeat("a", 51)}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tagRepo := new(MockTagRepository)
			listRepo := new(MockListRepository)

			tag, err := newTagService(tagRepo, listRepo).CreateTag(context.Background(), test.userID, test.tag)

			assert.ErrorIs(t, err, types.ErrInvalidTag)
			assert.Empty(t, tag)
			listRepo.AssertNotCalled(t, "FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything)
		})
	}
}

// Cenário: CreateTag não cria quando a lista não pertence ao usuário.
func TestCreateTag_ListNotFound(t *testing.T) {
	tagRepo := new(MockTagRepository)
	listRepo := new(MockListRepository)
	listRepo.On("FindByIDAndUserID", mock.Anything, int64(8), int64(7)).Return(nil, types.ErrListNotFound).Once()

	_, err := newTagService(tagRepo, listRepo).CreateTag(context.Background(), 7, types.Tag{ListID: 8, Name: "a"})

	assert.ErrorIs(t, err, types.ErrListNotFound)
	tagRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything, mock.Anything)
}

// Cenário: CreateTag propaga erro do repositório e não devolve tag parcial.
func TestCreateTag_RepositoryError(t *testing.T) {
	dbErr := errors.New("database unavailable")
	tagRepo := new(MockTagRepository)
	listRepo := new(MockListRepository)
	listRepo.On("FindByIDAndUserID", mock.Anything, int64(8), int64(7)).Return(&types.List{ID: 8}, nil).Once()
	tagRepo.On("Create", mock.Anything, int64(7), mock.Anything).Return(dbErr).Once()

	tag, err := newTagService(tagRepo, listRepo).CreateTag(context.Background(), 7, types.Tag{ListID: 8, Name: "a"})

	assert.ErrorIs(t, err, dbErr)
	assert.Empty(t, tag)
}

// Cenário: UpdateTag altera lista, nome e cor, enviando tagID e userID ao repositório.
func TestUpdateTag_Success(t *testing.T) {
	tagRepo := new(MockTagRepository)
	listRepo := new(MockListRepository)
	tagRepo.On("FindByIDAndUserID", mock.Anything, int64(5), int64(7)).
		Return(types.Tag{ID: 5, ListID: 8, Name: "antigo"}, nil).Once()
	listRepo.On("FindByIDAndUserID", mock.Anything, int64(9), int64(7)).Return(&types.List{ID: 9}, nil).Once()
	tagRepo.On("Update", mock.Anything, int64(5), int64(7), mock.MatchedBy(func(tag *types.Tag) bool {
		return tag.ListID == 9 && tag.Name == "novo" && tag.Color == "#00ff00"
	})).Return(nil).Once()

	tag, err := newTagService(tagRepo, listRepo).UpdateTag(context.Background(), 5, 7,
		types.Tag{ListID: 9, Name: "novo", Color: "#00ff00"})

	assert.NoError(t, err)
	assert.Equal(t, "novo", tag.Name)
	tagRepo.AssertExpectations(t)
	listRepo.AssertExpectations(t)
}

// Cenário: UpdateTag rejeita dados inválidos sem consultar repositórios.
func TestUpdateTag_InvalidInput(t *testing.T) {
	tagRepo := new(MockTagRepository)
	listRepo := new(MockListRepository)

	tag, err := newTagService(tagRepo, listRepo).UpdateTag(context.Background(), 0, 7, types.Tag{ListID: 8, Name: "a"})

	assert.ErrorIs(t, err, types.ErrInvalidTag)
	assert.Empty(t, tag)
	tagRepo.AssertNotCalled(t, "FindByIDAndUserID", mock.Anything, mock.Anything, mock.Anything)
}

// Cenário: UpdateTag falha quando a tag não pertence ao usuário.
func TestUpdateTag_TagNotFound(t *testing.T) {
	tagRepo := new(MockTagRepository)
	listRepo := new(MockListRepository)
	tagRepo.On("FindByIDAndUserID", mock.Anything, int64(5), int64(7)).Return(types.Tag{}, types.ErrTagNotFound).Once()

	_, err := newTagService(tagRepo, listRepo).UpdateTag(context.Background(), 5, 7, types.Tag{ListID: 8, Name: "a"})

	assert.ErrorIs(t, err, types.ErrTagNotFound)
	tagRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

// Cenário: UpdateTag falha quando a nova lista não pertence ao usuário.
func TestUpdateTag_NewListNotFound(t *testing.T) {
	tagRepo := new(MockTagRepository)
	listRepo := new(MockListRepository)
	tagRepo.On("FindByIDAndUserID", mock.Anything, int64(5), int64(7)).Return(types.Tag{ID: 5}, nil).Once()
	listRepo.On("FindByIDAndUserID", mock.Anything, int64(9), int64(7)).Return(nil, types.ErrListNotFound).Once()

	_, err := newTagService(tagRepo, listRepo).UpdateTag(context.Background(), 5, 7, types.Tag{ListID: 9, Name: "a"})

	assert.ErrorIs(t, err, types.ErrListNotFound)
	tagRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

// Cenário: UpdateTag propaga erro ao gravar a atualização.
func TestUpdateTag_RepositoryError(t *testing.T) {
	dbErr := errors.New("database unavailable")
	tagRepo := new(MockTagRepository)
	listRepo := new(MockListRepository)
	tagRepo.On("FindByIDAndUserID", mock.Anything, int64(5), int64(7)).Return(types.Tag{ID: 5}, nil).Once()
	listRepo.On("FindByIDAndUserID", mock.Anything, int64(9), int64(7)).Return(&types.List{ID: 9}, nil).Once()
	tagRepo.On("Update", mock.Anything, int64(5), int64(7), mock.Anything).Return(dbErr).Once()

	tag, err := newTagService(tagRepo, listRepo).UpdateTag(context.Background(), 5, 7, types.Tag{ListID: 9, Name: "a"})

	assert.ErrorIs(t, err, dbErr)
	assert.Empty(t, tag)
}

// Cenário: DeleteTag repassa tagID e userID ao repositório.
func TestDeleteTag_Success(t *testing.T) {
	tagRepo := new(MockTagRepository)
	tagRepo.On("Delete", mock.Anything, int64(5), int64(7)).Return(nil).Once()

	err := newTagService(tagRepo, new(MockListRepository)).DeleteTag(context.Background(), 5, 7)

	assert.NoError(t, err)
	tagRepo.AssertExpectations(t)
}

// Cenário: DeleteTag rejeita IDs inválidos sem chamar o repositório.
func TestDeleteTag_InvalidInput(t *testing.T) {
	tagRepo := new(MockTagRepository)

	err := newTagService(tagRepo, new(MockListRepository)).DeleteTag(context.Background(), 0, 7)

	assert.ErrorIs(t, err, types.ErrInvalidTag)
	tagRepo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything, mock.Anything)
}

// Cenário: GetTagsByList valida a lista do usuário e devolve suas tags.
func TestGetTagsByList_Success(t *testing.T) {
	tagRepo := new(MockTagRepository)
	listRepo := new(MockListRepository)
	listRepo.On("FindByIDAndUserID", mock.Anything, int64(8), int64(7)).Return(&types.List{ID: 8}, nil).Once()
	tagRepo.On("FindByListIDAndUserID", mock.Anything, int64(8), int64(7)).
		Return([]types.Tag{{ID: 1}, {ID: 2}}, nil).Once()

	tags, err := newTagService(tagRepo, listRepo).GetTagsByList(context.Background(), 8, 7)

	assert.NoError(t, err)
	assert.Len(t, tags, 2)
}

// Cenário: GetTagsByList rejeita IDs inválidos.
func TestGetTagsByList_InvalidInput(t *testing.T) {
	tags, err := newTagService(new(MockTagRepository), new(MockListRepository)).GetTagsByList(context.Background(), 0, 7)

	assert.ErrorIs(t, err, types.ErrInvalidTag)
	assert.Nil(t, tags)
}

// Cenário: GetTagsByList não lista tags de uma lista que não pertence ao usuário.
func TestGetTagsByList_ListNotFound(t *testing.T) {
	tagRepo := new(MockTagRepository)
	listRepo := new(MockListRepository)
	listRepo.On("FindByIDAndUserID", mock.Anything, int64(8), int64(7)).Return(nil, types.ErrListNotFound).Once()

	_, err := newTagService(tagRepo, listRepo).GetTagsByList(context.Background(), 8, 7)

	assert.ErrorIs(t, err, types.ErrListNotFound)
	tagRepo.AssertNotCalled(t, "FindByListIDAndUserID", mock.Anything, mock.Anything, mock.Anything)
}

// Cenário: GetTagById devolve a tag do usuário.
func TestGetTagById_Success(t *testing.T) {
	tagRepo := new(MockTagRepository)
	tagRepo.On("FindByIDAndUserID", mock.Anything, int64(5), int64(7)).Return(types.Tag{ID: 5}, nil).Once()

	tag, err := newTagService(tagRepo, new(MockListRepository)).GetTagById(context.Background(), 5, 7)

	assert.NoError(t, err)
	assert.Equal(t, int64(5), tag.ID)
}

// Cenário: GetTagById rejeita IDs inválidos.
func TestGetTagById_InvalidInput(t *testing.T) {
	tag, err := newTagService(new(MockTagRepository), new(MockListRepository)).GetTagById(context.Background(), 5, 0)

	assert.ErrorIs(t, err, types.ErrInvalidTag)
	assert.Empty(t, tag)
}
