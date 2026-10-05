package service

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/seu-usuario/taskflow-backend/internal/dto"
)

var ErrInvalidTask = errors.New("invalid task")

type TaskRepository interface {
	CreateTask(ctx context.Context, title, priority string) (int64, error)
}

type TaskService interface {
	CreateTask(ctx context.Context, request dto.CreateTaskRequest) (dto.TaskResponse, error)
}

type taskService struct {
	repository TaskRepository
}

func NewTaskService(repository TaskRepository) TaskService {
	return &taskService{repository: repository}
}

func (s *taskService) CreateTask(ctx context.Context, request dto.CreateTaskRequest) (dto.TaskResponse, error) {
	title := strings.TrimSpace(request.Title)
	if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 100 {
		return dto.TaskResponse{}, ErrInvalidTask
	}
	if request.Priority != "low" && request.Priority != "medium" && request.Priority != "high" {
		return dto.TaskResponse{}, ErrInvalidTask
	}

	id, err := s.repository.CreateTask(ctx, title, request.Priority)
	if err != nil {
		return dto.TaskResponse{}, err
	}

	return dto.TaskResponse{
		ID:       id,
		Title:    title,
		Priority: request.Priority,
		Status:   "pending",
	}, nil
}
