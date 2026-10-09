package service

import (
	"context"
	"time"

	"github.com/seu-usuario/taskflow-backend/internal/types"
)

type listService struct {
	repository types.ListRepository
}

var _ types.ListService = (*listService)(nil)

func NewListService(repository types.ListRepository) types.ListService {
	return &listService{repository: repository}
}

func (s *listService) CreateList(ctx context.Context, userID int64, list types.List) (types.List, error) {
	list.UserID = userID
	list.CreatedAt = time.Now()
	if err := s.repository.Create(ctx, &list); err != nil {
		return types.List{}, err
	}
	return list, nil
}

func (s *listService) GetListByID(ctx context.Context, id, userID int64) (types.List, error) {
	list, err := s.repository.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return types.List{}, err
	}
	return *list, nil
}

func (s *listService) GetUserLists(ctx context.Context, userID int64) ([]types.List, error) {
	return s.repository.FindByUserID(ctx, userID)
}

func (s *listService) UpdateList(ctx context.Context, id, userID int64, update types.List) (types.List, error) {
	list, err := s.repository.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return types.List{}, err
	}

	list.Title = update.Title
	list.Description = update.Description
	if err := s.repository.Update(ctx, list); err != nil {
		return types.List{}, err
	}
	return *list, nil
}

func (s *listService) DeleteList(ctx context.Context, id, userID int64) error {
	return s.repository.Delete(ctx, id, userID)
}
