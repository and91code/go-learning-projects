package list

import (
	"context"
	"fmt"

	"github.com/seu-usuario/taskflow-backend/internal/types"
)

type ListLists interface { Execute(context.Context, int64) ([]types.List, error) }
type listLists struct{ repository types.ListRepository }
func NewListLists(repository types.ListRepository) ListLists { return &listLists{repository} }
func (u *listLists) Execute(ctx context.Context, userID int64) ([]types.List, error) {
	items, err := u.repository.FindByUserID(ctx, userID)
	if err != nil { return nil, fmt.Errorf("list lists: %w", err) }
	return items, nil
}
