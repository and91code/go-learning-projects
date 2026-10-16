package list

import (
	"context"
	"fmt"
	"time"

	"github.com/seu-usuario/taskflow-backend/internal/types"
)

type CreateList interface {
	Execute(context.Context, int64, types.List) (types.List, error)
}

type createList struct{ repository types.ListRepository }

func NewCreateList(repository types.ListRepository) CreateList {
	return &createList{repository: repository}
}

func (u *createList) Execute(ctx context.Context, userID int64, item types.List) (types.List, error) {
	item.UserID, item.CreatedAt = userID, time.Now()
	if err := u.repository.Create(ctx, userID, &item); err != nil {
		return types.List{}, fmt.Errorf("create list: %w", err)
	}
	return item, nil
}
