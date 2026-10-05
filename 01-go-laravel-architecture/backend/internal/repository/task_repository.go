package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) CreateTask(ctx context.Context, title, priority string) (int64, error) {
	result, err := r.db.ExecContext(
		ctx,
		"INSERT INTO tasks (title, priority, status) VALUES (?, ?, ?)",
		title,
		priority,
		"pending",
	)
	if err != nil {
		return 0, fmt.Errorf("create task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get created task id: %w", err)
	}

	return id, nil
}
