package types

import (
	"context"
	"time"
)

type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
)

type TaskPriority string

const (
	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
)

type Task struct {
	ID          int64
	ListID      int64
	Title       string
	Description string
	Status      TaskStatus
	Priority    TaskPriority
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Tags        []Tag
}

type TaskRepository interface {
	Create(ctx context.Context, userID int64, task *Task) error
	Update(ctx context.Context, taskID, userID int64, task *Task) error
	Delete(ctx context.Context, taskID, userID int64) error
	FindByListIDAndUserID(ctx context.Context, listID, userID int64) ([]Task, error)
	FindByIDAndUserID(ctx context.Context, taskID, userID int64) (*Task, error)
}

type TaskService interface {
	CreateTask(ctx context.Context, userID int64, task Task) (Task, error)
	UpdateTask(ctx context.Context, taskID, userID int64, task Task) (Task, error)
	DeleteTask(ctx context.Context, taskID, userID int64) error
	GetTasksByList(ctx context.Context, listID, userID int64) ([]Task, error)
	GetTaskByID(ctx context.Context, taskID, userID int64) (Task, error)
}
