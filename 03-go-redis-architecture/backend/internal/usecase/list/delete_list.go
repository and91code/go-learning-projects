package list

import (
	"context"
	"fmt"

	"github.com/seu-usuario/taskflow-backend/internal/types"
)
type DeleteList interface { Execute(context.Context, int64, int64) error }
type deleteList struct{ repository types.ListRepository }
func NewDeleteList(r types.ListRepository) DeleteList { return &deleteList{r} }
func (u *deleteList) Execute(ctx context.Context, id, userID int64) error {
	if err := u.repository.Delete(ctx, id, userID); err != nil { return fmt.Errorf("delete list: %w", err) }
	return nil
}
