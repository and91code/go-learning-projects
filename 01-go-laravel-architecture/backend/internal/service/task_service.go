package service

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/seu-usuario/taskflow-backend/internal/types"
)

type taskService struct {
	taskRepository types.TaskRepository
	listRepository types.ListRepository
}

var _ types.TaskService = (*taskService)(nil)

func NewTaskService(taskRepository types.TaskRepository, listRepository types.ListRepository) types.TaskService {
	return &taskService{taskRepository: taskRepository, listRepository: listRepository}
}

func (s *taskService) CreateTask(ctx context.Context, userID int64, task types.Task) (types.Task, error) {
	task.Title = strings.TrimSpace(task.Title)
	if userID <= 0 || task.ListID <= 0 || !validTaskTitle(task.Title) || !validTaskPriority(task.Priority) {
		return types.Task{}, types.ErrInvalidTask
	}
	if _, err := s.listRepository.FindByIDAndUserID(ctx, task.ListID, userID); err != nil {
		return types.Task{}, err
	}

	task.Status = types.StatusPending
	task.CreatedAt = time.Now()
	task.UpdatedAt = task.CreatedAt
	if err := s.taskRepository.Create(ctx, &task); err != nil {
		return types.Task{}, err
	}
	return task, nil
}

func (s *taskService) GetTasksByList(ctx context.Context, listID, userID int64) ([]types.Task, error) {
	if listID <= 0 || userID <= 0 {
		return nil, types.ErrInvalidTask
	}
	return s.taskRepository.FindByListIDAndUserID(ctx, listID, userID)
}

func (s *taskService) GetTaskByID(ctx context.Context, taskID, userID int64) (types.Task, error) {
	if taskID <= 0 || userID <= 0 {
		return types.Task{}, types.ErrInvalidTask
	}
	task, err := s.taskRepository.FindByIDAndUserID(ctx, taskID, userID)
	if err != nil {
		return types.Task{}, err
	}
	return *task, nil
}

func (s *taskService) UpdateTask(ctx context.Context, taskID, userID int64, update types.Task) (types.Task, error) {
	update.Title = strings.TrimSpace(update.Title)
	if taskID <= 0 || userID <= 0 || !validTaskTitle(update.Title) ||
		!validTaskPriority(update.Priority) || !validTaskStatus(update.Status) {
		return types.Task{}, types.ErrInvalidTask
	}

	task, err := s.taskRepository.FindByIDAndUserID(ctx, taskID, userID)
	if err != nil {
		return types.Task{}, err
	}
	task.Title = update.Title
	task.Description = strings.TrimSpace(update.Description)
	task.Status = update.Status
	task.Priority = update.Priority
	task.DueDate = update.DueDate
	task.UpdatedAt = time.Now()

	if err := s.taskRepository.Update(ctx, task, userID); err != nil {
		return types.Task{}, err
	}
	return *task, nil
}

func (s *taskService) DeleteTask(ctx context.Context, taskID, userID int64) error {
	if taskID <= 0 || userID <= 0 {
		return types.ErrInvalidTask
	}
	return s.taskRepository.Delete(ctx, taskID, userID)
}

func validTaskTitle(title string) bool {
	length := utf8.RuneCountInString(title)
	return length >= 3 && length <= 150
}

func validTaskPriority(priority types.TaskPriority) bool {
	return priority == types.PriorityLow || priority == types.PriorityMedium || priority == types.PriorityHigh
}

func validTaskStatus(status types.TaskStatus) bool {
	return status == types.StatusPending || status == types.StatusInProgress || status == types.StatusCompleted
}
