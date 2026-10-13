package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/seu-usuario/taskflow-backend/internal/types"
)

type mysqlTaskRepository struct {
	db *sql.DB
}

var _ types.TaskRepository = (*mysqlTaskRepository)(nil)

func NewTaskRepository(db *sql.DB) types.TaskRepository {
	return &mysqlTaskRepository{db: db}
}

func (r *mysqlTaskRepository) Create(ctx context.Context, userID int64, task *types.Task) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("taskRepository.Create begin transaction: %w", err)
	}
	defer rollbackTaskTransaction(tx, &err, "taskRepository.Create")

	query := `
		INSERT INTO tasks (list_id, title, description, status, priority, due_date)
		SELECT l.id, ?, ?, ?, ?, ?
		FROM lists l
		WHERE l.id = ? AND l.user_id = ?
	`
	result, err := tx.ExecContext(ctx, query,
		task.Title, task.Description, task.Status, task.Priority, task.DueDate, task.ListID, userID,
	)
	if err != nil {
		return fmt.Errorf("taskRepository.Create: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("taskRepository.Create RowsAffected: %w", err)
	}
	if rows == 0 {
		return types.ErrListNotFound
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("taskRepository.Create LastInsertId: %w", err)
	}

	task.ID = id
	task.Tags, err = writeTaskTagLinks(ctx, tx, task.ID, task.ListID, task.Tags, false)
	if err != nil {
		return fmt.Errorf("taskRepository.Create tags: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("taskRepository.Create commit: %w", err)
	}
	return nil
}

func (r *mysqlTaskRepository) Update(ctx context.Context, taskID, userID int64, task *types.Task) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("taskRepository.Update begin transaction: %w", err)
	}
	defer rollbackTaskTransaction(tx, &err, "taskRepository.Update")

	query := `
		UPDATE tasks t
		INNER JOIN lists l ON l.id = t.list_id
		SET t.title = ?, t.description = ?, t.status = ?, t.priority = ?, t.due_date = ?
		WHERE t.id = ? AND l.user_id = ?
	`

	result, err := tx.ExecContext(ctx, query,
		task.Title, task.Description, task.Status, task.Priority, task.DueDate, taskID, userID,
	)
	if err != nil {
		return fmt.Errorf("taskRepository.Update: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("taskRepository.Update RowsAffected: %w", err)
	}
	if rows == 0 {
		var existingTaskID int64
		err := tx.QueryRowContext(ctx, `
			SELECT t.id
			FROM tasks t
			INNER JOIN lists l ON l.id = t.list_id
			WHERE t.id = ? AND l.user_id = ?
		`, taskID, userID).Scan(&existingTaskID)
		if errors.Is(err, sql.ErrNoRows) {
			return types.ErrTaskNotFound
		}
		if err != nil {
			return fmt.Errorf("taskRepository.Update check task: %w", err)
		}
	}

	if task.Tags != nil {
		task.Tags, err = writeTaskTagLinks(ctx, tx, taskID, task.ListID, task.Tags, true)
		if err != nil {
			return fmt.Errorf("taskRepository.Update tags: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("taskRepository.Update commit: %w", err)
	}
	return nil
}

func rollbackTaskTransaction(tx *sql.Tx, operationErr *error, operation string) {
	if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
		*operationErr = errors.Join(*operationErr, fmt.Errorf("%s rollback: %w", operation, err))
	}
}

func writeTaskTagLinks(
	ctx context.Context,
	tx *sql.Tx,
	taskID, listID int64,
	tags []types.Tag,
	replace bool,
) ([]types.Tag, error) {
	validatedTags, err := findTagsForTask(ctx, tx, listID, tags)
	if err != nil {
		return nil, err
	}

	if replace {
		if _, err := tx.ExecContext(ctx, `DELETE FROM task_tags WHERE task_id = ?`, taskID); err != nil {
			return nil, fmt.Errorf("delete existing task tag links: %w", err)
		}
	}

	for _, tag := range validatedTags {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO task_tags (task_id, tag_id) VALUES (?, ?)`,
			taskID, tag.ID,
		); err != nil {
			return nil, fmt.Errorf("insert task tag link: %w", err)
		}
	}
	return validatedTags, nil
}

func findTagsForTask(ctx context.Context, tx *sql.Tx, listID int64, tags []types.Tag) ([]types.Tag, error) {
	if len(tags) == 0 {
		return []types.Tag{}, nil
	}

	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(tags)), ",")
	query := fmt.Sprintf(
		`SELECT id, list_id, name, color FROM tags WHERE list_id = ? AND id IN (%s)`,
		placeholders,
	)
	args := make([]any, 0, len(tags)+1)
	args = append(args, listID)
	seen := make(map[int64]struct{}, len(tags))
	for _, tag := range tags {
		if tag.ID <= 0 {
			return nil, types.ErrInvalidTask
		}
		if _, exists := seen[tag.ID]; exists {
			return nil, types.ErrInvalidTask
		}
		seen[tag.ID] = struct{}{}
		args = append(args, tag.ID)
	}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("find tags for task: %w", err)
	}
	foundTags := make(map[int64]types.Tag, len(tags))
	for rows.Next() {
		var tag types.Tag
		if err := rows.Scan(&tag.ID, &tag.ListID, &tag.Name, &tag.Color); err != nil {
			return nil, errors.Join(
				fmt.Errorf("scan tags for task: %w", err),
				rows.Close(),
			)
		}
		foundTags[tag.ID] = tag
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Join(
			fmt.Errorf("read tags for task: %w", err),
			rows.Close(),
		)
	}
	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("close tags for task: %w", err)
	}

	validatedTags := make([]types.Tag, 0, len(tags))
	for _, tag := range tags {
		found, exists := foundTags[tag.ID]
		if !exists {
			return nil, types.ErrTagNotFound
		}
		validatedTags = append(validatedTags, found)
	}
	return validatedTags, nil
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

func (r *mysqlTaskRepository) FindByIDAndUserID(ctx context.Context, taskID, userID int64) (*types.Task, error) {
	query := `
		SELECT t.id, t.list_id, t.title, t.description, t.status, t.priority, t.due_date, t.created_at, t.updated_at,
			g.id, g.name, g.color
		FROM tasks t
		INNER JOIN lists l ON l.id = t.list_id
		LEFT JOIN task_tags tt ON tt.task_id = t.id
		LEFT JOIN tags g ON g.id = tt.tag_id AND g.list_id = t.list_id
		WHERE t.id = ? AND l.user_id = ?
		ORDER BY g.id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, taskID, userID)
	if err != nil {
		return nil, fmt.Errorf("taskRepository.FindByIDAndUserID: %w", err)
	}
	defer rows.Close()

	tasks, err := scanTaskRows(rows)
	if err != nil {
		return nil, fmt.Errorf("taskRepository.FindByIDAndUserID scan: %w", err)
	}
	if len(tasks) == 0 {
		return nil, types.ErrTaskNotFound
	}

	return &tasks[0], nil
}

func (r *mysqlTaskRepository) FindByListIDAndUserID(ctx context.Context, listID, userID int64) ([]types.Task, error) {
	query := `
		SELECT t.id, t.list_id, t.title, t.description, t.status, t.priority, t.due_date, t.created_at, t.updated_at,
			g.id, g.name, g.color
		FROM tasks t
		INNER JOIN lists l ON l.id = t.list_id
		LEFT JOIN task_tags tt ON tt.task_id = t.id
		LEFT JOIN tags g ON g.id = tt.tag_id AND g.list_id = t.list_id
		WHERE t.list_id = ? AND l.user_id = ?
		ORDER BY t.created_at DESC, t.id DESC, g.id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, listID, userID)
	if err != nil {
		return nil, fmt.Errorf("taskRepository.FindByListIDAndUserID: %w", err)
	}
	defer rows.Close()

	tasks, err := scanTaskRows(rows)
	if err != nil {
		return nil, fmt.Errorf("taskRepository.FindByListIDAndUserID scan: %w", err)
	}
	return tasks, nil
}

func scanTaskRows(rows *sql.Rows) ([]types.Task, error) {
	var tasks []types.Task
	taskIndexes := make(map[int64]int)
	for rows.Next() {
		var t types.Task
		var desc sql.NullString
		var dueDate sql.NullTime
		var tagID sql.NullInt64
		var tagName, tagColor sql.NullString

		if err := rows.Scan(
			&t.ID, &t.ListID, &t.Title, &desc, &t.Status, &t.Priority, &dueDate, &t.CreatedAt, &t.UpdatedAt,
			&tagID, &tagName, &tagColor,
		); err != nil {
			return nil, err
		}

		if desc.Valid {
			t.Description = desc.String
		}
		if dueDate.Valid {
			t.DueDate = &dueDate.Time
		}

		index, exists := taskIndexes[t.ID]
		if !exists {
			index = len(tasks)
			taskIndexes[t.ID] = index
			tasks = append(tasks, t)
		}
		if tagID.Valid {
			tasks[index].Tags = append(tasks[index].Tags, types.Tag{
				ID:     tagID.Int64,
				ListID: t.ListID,
				Name:   tagName.String,
				Color:  tagColor.String,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}
