package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/seu-usuario/taskflow-backend/internal/types"
)

type mysqlUserRepository struct {
	db *sql.DB
}

var _ types.UserRepository = (*mysqlUserRepository)(nil)

func NewUserRepository(db *sql.DB) types.UserRepository {
	return &mysqlUserRepository{db: db}
}

func (r *mysqlUserRepository) Create(ctx context.Context, user *types.User) error {
	query := `INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("userRepository.Create: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("userRepository.Create LastInsertId: %w", err)
	}

	user.ID = id
	return nil
}

func (r *mysqlUserRepository) FindByEmail(ctx context.Context, email string) (*types.User, error) {
	query := `SELECT id, name, email, password_hash, created_at FROM users WHERE email = ?`

	var u types.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.ErrUserNotFound
		}
		return nil, fmt.Errorf("userRepository.FindByEmail: %w", err)
	}

	return &u, nil
}

func (r *mysqlUserRepository) FindByID(ctx context.Context, id int64) (*types.User, error) {
	query := `SELECT id, name, email, password_hash, created_at FROM users WHERE id = ?`

	var u types.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.ErrUserNotFound
		}
		return nil, fmt.Errorf("userRepository.FindByID: %w", err)
	}

	return &u, nil
}

func (r *mysqlUserRepository) UpdateProfile(ctx context.Context, id int64, name, email string) error {
	query := `UPDATE users SET name = ?, email = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, name, email, id)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return types.ErrEmailExists
		}
		return fmt.Errorf("userRepository.UpdateProfile: %w", err)
	}
	return nil
}

func (r *mysqlUserRepository) UpdatePassword(ctx context.Context, id int64, passwordHash string) error {
	query := `UPDATE users SET password_hash = ? WHERE id = ?`
	if _, err := r.db.ExecContext(ctx, query, passwordHash, id); err != nil {
		return fmt.Errorf("userRepository.UpdatePassword: %w", err)
	}
	return nil
}
