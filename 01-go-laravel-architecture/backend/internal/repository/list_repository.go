package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/seu-usuario/taskflow-backend/internal/domain"
)

var ErrListNotFound = errors.New("lista não encontrada")

type ListRepository interface {
	Create(ctx context.Context, list *domain.List) error
	FindByIDAndUserID(ctx context.Context, id, userID int64) (*domain.List, error)
	FindByUserID(ctx context.Context, userID int64) ([]domain.List, error)
	Update(ctx context.Context, list *domain.List) error
	Delete(ctx context.Context, id, userID int64) error
}

type mysqlListRepository struct {
	db *sql.DB
}

func NewListRepository(db *sql.DB) ListRepository {
	return &mysqlListRepository{db: db}
}

func (r *mysqlListRepository) Create(ctx context.Context, list *domain.List) error {
	query := `INSERT INTO lists (user_id, title, description) VALUES (?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, list.UserID, list.Title, list.Description)
	if err != nil {
		return fmt.Errorf("repository.Create: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("repository.Create LastInsertId: %w", err)
	}

	list.ID = id
	return nil
}

func (r *mysqlListRepository) FindByIDAndUserID(ctx context.Context, id, userID int64) (*domain.List, error) {
	query := `SELECT id, user_id, title, description, created_at FROM lists WHERE id = ? AND user_id = ?`

	var list domain.List
	var desc sql.NullString // Trata descrição opcional (NULL no MySQL)

	err := r.db.QueryRowContext(ctx, query, id, userID).Scan(
		&list.ID,
		&list.UserID,
		&list.Title,
		&desc,
		&list.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrListNotFound
		}
		return nil, fmt.Errorf("repository.FindByIDAndUserID: %w", err)
	}

	if desc.Valid {
		list.Description = desc.String
	}

	return &list, nil
}

func (r *mysqlListRepository) FindByUserID(ctx context.Context, userID int64) ([]domain.List, error) {
	query := `SELECT id, user_id, title, description, created_at FROM lists WHERE user_id = ? ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("repository.FindByUserID: %w", err)
	}
	defer rows.Close()

	var lists []domain.List
	for rows.Next() {
		var l domain.List
		var desc sql.NullString

		if err := rows.Scan(&l.ID, &l.UserID, &l.Title, &desc, &l.CreatedAt); err != nil {
			return nil, fmt.Errorf("repository.FindByUserID scan: %w", err)
		}

		if desc.Valid {
			l.Description = desc.String
		}

		lists = append(lists, l)
	}

	return lists, nil
}

func (r *mysqlListRepository) Update(ctx context.Context, list *domain.List) error {
	query := `UPDATE lists SET title = ?, description = ? WHERE id = ? AND user_id = ?`

	result, err := r.db.ExecContext(ctx, query, list.Title, list.Description, list.ID, list.UserID)
	if err != nil {
		return fmt.Errorf("repository.Update: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository.Update RowsAffected: %w", err)
	}

	if rows == 0 {
		return ErrListNotFound
	}

	return nil
}

func (r *mysqlListRepository) Delete(ctx context.Context, id, userID int64) error {
	query := `DELETE FROM lists WHERE id = ? AND user_id = ?`

	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("repository.Delete: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository.Delete RowsAffected: %w", err)
	}

	if rows == 0 {
		return ErrListNotFound
	}

	return nil
}
