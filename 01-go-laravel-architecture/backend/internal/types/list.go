package types

import (
	"context"
	"time"
)

type List struct {
	ID          int64
	UserID      int64
	Title       string
	Description string
	CreatedAt   time.Time
	Tasks       []Task
	Tags        []Tag
}

type ListRepository interface {
	Create(ctx context.Context, list *List) error
	FindByIDAndUserID(ctx context.Context, id, userID int64) (*List, error)
	FindByUserID(ctx context.Context, userID int64) ([]List, error)
	Update(ctx context.Context, list *List) error
	Delete(ctx context.Context, id, userID int64) error
}

type ListService interface {
	CreateList(ctx context.Context, userID int64, list List) (List, error)
	GetListByID(ctx context.Context, id, userID int64) (List, error)
	GetUserLists(ctx context.Context, userID int64) ([]List, error)
	UpdateList(ctx context.Context, id, userID int64, list List) (List, error)
	DeleteList(ctx context.Context, id, userID int64) error
}
