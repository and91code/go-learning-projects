package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/seu-usuario/taskflow-backend/internal/types"
)

type mysqlTagRepository struct {
	db *sql.DB
}

var _ types.TagRepository = (*mysqlTagRepository)(nil)

func NewTagRepository(db *sql.DB) types.TagRepository {
	return &mysqlTagRepository{db: db}
}

func (r *mysqlTagRepository) Create(ctx context.Context, userID int64, tag *types.Tag) error {
	query := `
		INSERT INTO tags (list_id, name, color)
		SELECT l.id, ?, ?
		FROM lists l
		WHERE l.id = ? AND l.user_id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		tag.Name, tag.Color, tag.ListID, userID,
	)
	if err != nil {
		return fmt.Errorf("tagRepository.Create: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("tagRepository.Create RowsAffected: %w", err)
	}
	if rows == 0 {
		return types.ErrListNotFound
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("tagRepository.Create LastInsertId: %w", err)
	}

	tag.ID = id
	return nil
}

func (r *mysqlTagRepository) Update(ctx context.Context, tagID, userID int64, tag *types.Tag) error {
	query := `
		UPDATE tags t
		INNER JOIN lists l ON l.id = t.list_id
		SET t.list_id = ?, t.name = ?, t.color = ?
		WHERE t.id = ? AND l.user_id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		tag.ListID, tag.Name, tag.Color, tagID, userID,
	)
	if err != nil {
		return fmt.Errorf("tagRepository.Update: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("tagRepository.Update RowsAffected: %w", err)
	}
	if rows == 0 {
		if _, err := r.FindByIDAndUserID(ctx, tagID, userID); err != nil {
			return err
		}
		var existingListID int64
		err := r.db.QueryRowContext(ctx,
			`SELECT id FROM lists WHERE id = ? AND user_id = ?`,
			tag.ListID, userID,
		).Scan(&existingListID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return types.ErrListNotFound
			}
			return fmt.Errorf("tagRepository.Update check list: %w", err)
		}
	}

	return nil
}

func (r *mysqlTagRepository) Delete(ctx context.Context, tagID, userID int64) error {
	query := `
		DELETE t FROM tags t
		INNER JOIN lists l ON l.id = t.list_id
		WHERE t.id = ? AND l.user_id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		tagID, userID,
	)
	if err != nil {
		return fmt.Errorf("tagRepository.Delete: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("tagRepository.Delete RowsAffected: %w", err)
	}
	if rows == 0 {
		return types.ErrTagNotFound
	}

	return nil
}

func (r *mysqlTagRepository) FindByIDAndUserID(ctx context.Context, tagID, userID int64) (types.Tag, error) {
	query := `
		SELECT t.id, t.list_id, t.name, t.color
		FROM tags t
		INNER JOIN lists l ON l.id = t.list_id
		WHERE t.id = ? AND l.user_id = ?
	`

	var t types.Tag

	err := r.db.QueryRowContext(ctx, query, tagID, userID).Scan(
		&t.ID, &t.ListID, &t.Name, &t.Color,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.Tag{}, types.ErrTagNotFound
		}
		return types.Tag{}, fmt.Errorf("TagRepository.FindByIDAndUserID: %w", err)
	}

	return t, nil
}

func (r *mysqlTagRepository) FindByListIDAndUserID(ctx context.Context, listID, userID int64) ([]types.Tag, error) {
	query := `
		SELECT t.id, t.list_id, t.name, t.color
		FROM tags t
		INNER JOIN lists l ON l.id = t.list_id
		WHERE t.list_id = ? AND l.user_id = ?
		ORDER BY t.id ASC
	`

	rows, err := r.db.QueryContext(ctx, query,
		listID, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("tagRepository.FindByListIDAndUserID: %w", err)
	}
	defer rows.Close()

	var tags []types.Tag
	for rows.Next() {
		var tag types.Tag
		if err := rows.Scan(&tag.ID, &tag.ListID, &tag.Name, &tag.Color); err != nil {
			return nil, fmt.Errorf("tagRepository.FindByListIDAndUserID scan: %w", err)
		}
		tags = append(tags, tag)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("tagRepository.FindByListIDAndUserID rows: %w", err)
	}

	return tags, nil
}
