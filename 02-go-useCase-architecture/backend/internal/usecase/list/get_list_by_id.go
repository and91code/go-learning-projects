package list

import (
	"context"
	"fmt"

	"github.com/seu-usuario/taskflow-backend/internal/types"
)

type GetListByID interface { Execute(context.Context, int64, int64) (types.List, error) }
type getListByID struct{ repository types.ListRepository }
func NewGetListByID(repository types.ListRepository) GetListByID { return &getListByID{repository} }
func (u *getListByID) Execute(ctx context.Context, id, userID int64) (types.List, error) {
	item, err := u.repository.FindByIDAndUserID(ctx, id, userID)
	if err != nil { return types.List{}, fmt.Errorf("get list: %w", err) }
	return *item, nil
}
