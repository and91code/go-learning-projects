package list

import (
	"context"
	"fmt"
	"github.com/seu-usuario/taskflow-backend/internal/types"
)
type UpdateList interface { Execute(context.Context, int64, int64, types.List) (types.List, error) }
type updateList struct{ repository types.ListRepository }
func NewUpdateList(r types.ListRepository) UpdateList { return &updateList{r} }
func (u *updateList) Execute(ctx context.Context, id, userID int64, update types.List) (types.List, error) {
	item, err := u.repository.FindByIDAndUserID(ctx, id, userID)
	if err != nil { return types.List{}, fmt.Errorf("find list for update: %w", err) }
	item.Title, item.Description = update.Title, update.Description
	if err := u.repository.Update(ctx, id, userID, item); err != nil { return types.List{}, fmt.Errorf("update list: %w", err) }
	return *item, nil
}
