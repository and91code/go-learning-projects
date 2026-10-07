package service

import (
	"context"
	"time"

	"github.com/seu-usuario/taskflow-backend/internal/domain"
	"github.com/seu-usuario/taskflow-backend/internal/dto"
	"github.com/seu-usuario/taskflow-backend/internal/repository"
)

type ListService interface {
	CreateList(ctx context.Context, userID int64, req dto.CreateListRequest) (dto.ListResponse, error)
	GetListByID(ctx context.Context, id, userID int64) (dto.ListResponse, error)
	GetUserLists(ctx context.Context, userID int64) ([]dto.ListResponse, error)
	UpdateList(ctx context.Context, id, userID int64, req dto.UpdateListRequest) (dto.ListResponse, error)
	DeleteList(ctx context.Context, id, userID int64) error
}

type listService struct {
	repo repository.ListRepository
}

func NewListService(r repository.ListRepository) ListService {
	return &listService{repo: r}
}

func (s *listService) CreateList(ctx context.Context, userID int64, req dto.CreateListRequest) (dto.ListResponse, error) {
	list := &domain.List{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, list); err != nil {
		return dto.ListResponse{}, err
	}

	return dto.FormatListResponse(list.ID, list.UserID, list.Title, list.Description, list.CreatedAt), nil
}

func (s *listService) GetListByID(ctx context.Context, id, userID int64) (dto.ListResponse, error) {
	list, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return dto.ListResponse{}, err
	}

	return dto.FormatListResponse(list.ID, list.UserID, list.Title, list.Description, list.CreatedAt), nil
}

func (s *listService) GetUserLists(ctx context.Context, userID int64) ([]dto.ListResponse, error) {
	lists, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ListResponse, 0, len(lists))
	for _, list := range lists {
		responses = append(responses, dto.FormatListResponse(list.ID, list.UserID, list.Title, list.Description, list.CreatedAt))
	}

	return responses, nil
}

func (s *listService) UpdateList(ctx context.Context, id, userID int64, req dto.UpdateListRequest) (dto.ListResponse, error) {
	list, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return dto.ListResponse{}, err
	}

	list.Title = req.Title
	list.Description = req.Description

	if err := s.repo.Update(ctx, list); err != nil {
		return dto.ListResponse{}, err
	}

	return dto.FormatListResponse(list.ID, list.UserID, list.Title, list.Description, list.CreatedAt), nil
}

func (s *listService) DeleteList(ctx context.Context, id, userID int64) error {
	return s.repo.Delete(ctx, id, userID)
}
