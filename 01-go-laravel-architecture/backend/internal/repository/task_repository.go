package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/seu-usuario/taskflow-backend/internal/types"
)

type mysqlTaskRepository struct {
	db *sql.DB
}

var _ types.TaskRepository = (*mysqlTaskRepository)(nil)

func NewTaskRepository(db *sql.DB) types.TaskRepository {
	return &mysqlTaskRepository{db: db}
}

func (r *mysqlTaskRepository) Create(ctx context.Context, task *types.Task) error {
	query := `
		INSERT INTO tasks (list_id, title, description, status, priority, due_date)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(ctx, query,
		task.ListID, task.Title, task.Description, task.Status, task.Priority, task.DueDate,
	)
	if err != nil {
		return fmt.Errorf("taskRepository.Create: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("taskRepository.Create LastInsertId: %w", err)
	}

	task.ID = id
	return nil
}

func (r *mysqlTaskRepository) FindByIDAndUserID(ctx context.Context, taskID, userID int64) (*types.Task, error) {
	query := `
		SELECT t.id, t.list_id, t.title, t.description, t.status, t.priority, t.due_date, t.created_at, t.updated_at
		FROM tasks t
		INNER JOIN lists l ON l.id = t.list_id
		WHERE t.id = ? AND l.user_id = ?
	`

	var t types.Task
	var desc sql.NullString
	var dueDate sql.NullTime

	err := r.db.QueryRowContext(ctx, query, taskID, userID).Scan(
		&t.ID, &t.ListID, &t.Title, &desc, &t.Status, &t.Priority, &dueDate, &t.CreatedAt, &t.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.ErrTaskNotFound
		}
		return nil, fmt.Errorf("taskRepository.FindByIDAndUserID: %w", err)
	}

	if desc.Valid {
		t.Description = desc.String
	}
	if dueDate.Valid {
		t.DueDate = &dueDate.Time
	}

	return &t, nil
}

func (r *mysqlTaskRepository) FindByListIDAndUserID(ctx context.Context, listID, userID int64) ([]types.Task, error) {
	query := `
		SELECT t.id, t.list_id, t.title, t.description, t.status, t.priority, t.due_date, t.created_at, t.updated_at
		FROM tasks t
		INNER JOIN lists l ON l.id = t.list_id
		WHERE t.list_id = ? AND l.user_id = ?
		ORDER BY t.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, listID, userID)
	if err != nil {
		return nil, fmt.Errorf("taskRepository.FindByListIDAndUserID: %w", err)
	}
	defer rows.Close()

	var tasks []types.Task
	for rows.Next() {
		var t types.Task
		var desc sql.NullString
		var dueDate sql.NullTime

		if err := rows.Scan(&t.ID, &t.ListID, &t.Title, &desc, &t.Status, &t.Priority, &dueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("taskRepository.FindByListIDAndUserID scan: %w", err)
		}

		if desc.Valid {
			t.Description = desc.String
		}
		if dueDate.Valid {
			t.DueDate = &dueDate.Time
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *mysqlTaskRepository) Update(ctx context.Context, task *types.Task, userID int64) error {
	query := `
		UPDATE tasks t
		INNER JOIN lists l ON l.id = t.list_id
		SET t.title = ?, t.description = ?, t.status = ?, t.priority = ?, t.due_date = ?
		WHERE t.id = ? AND l.user_id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		task.Title, task.Description, task.Status, task.Priority, task.DueDate, task.ID, userID,
	)
	if err != nil {
		return fmt.Errorf("taskRepository.Update: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("taskRepository.Update RowsAffected: %w", err)
	}
	if rows == 0 {
		if _, err := r.FindByIDAndUserID(ctx, task.ID, userID); err != nil {
			return err
		}
	}

	return nil
}

func (r *mysqlTaskRepository) Delete(ctx context.Context, taskID, userID int64) error {
	query := `
		DELETE t FROM tasks t
		INNER JOIN lists l ON l.id = t.list_id
		WHERE t.id = ? AND l.user_id = ?
	`

	result, err := r.db.ExecContext(ctx, query, taskID, userID)
	if err != nil {
		return fmt.Errorf("taskRepository.Delete: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("taskRepository.Delete RowsAffected: %w", err)
	}
	if rows == 0 {
		return types.ErrTaskNotFound
	}

	return nil
}
