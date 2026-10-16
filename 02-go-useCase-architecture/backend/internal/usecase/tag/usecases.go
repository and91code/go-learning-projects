package tag

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/seu-usuario/taskflow-backend/internal/types"
)

type Create interface {
	Execute(context.Context, int64, types.Tag) (types.Tag, error)
}
type Update interface {
	Execute(context.Context, int64, int64, types.Tag) (types.Tag, error)
}
type Delete interface {
	Execute(context.Context, int64, int64) error
}
type GetByList interface {
	Execute(context.Context, int64, int64) ([]types.Tag, error)
}
type GetByID interface {
	Execute(context.Context, int64, int64) (types.Tag, error)
}
type repositories struct {
	tags  types.TagRepository
	lists types.ListRepository
}

func New(tags types.TagRepository, lists types.ListRepository) (Create, Update, Delete, GetByList, GetByID) {
	r := &repositories{tags: tags, lists: lists}
	return create{r}, update{r}, remove{r}, byList{r}, byID{r}
}

type create struct{ *repositories }

func (u create) Execute(ctx context.Context, userID int64, tag types.Tag) (types.Tag, error) {
	tag.Name = strings.TrimSpace(tag.Name)
	if userID <= 0 || tag.ListID <= 0 || !validName(tag.Name) {
		return types.Tag{}, types.ErrInvalidTag
	}
	if _, err := u.lists.FindByIDAndUserID(ctx, tag.ListID, userID); err != nil {
		return types.Tag{}, err
	}
	if err := u.tags.Create(ctx, userID, &tag); err != nil {
		return types.Tag{}, err
	}
	return tag, nil
}

type update struct{ *repositories }

func (u update) Execute(ctx context.Context, tagID, userID int64, update types.Tag) (types.Tag, error) {
	update.Name = strings.TrimSpace(update.Name)
	if tagID <= 0 || userID <= 0 || update.ListID <= 0 || !validName(update.Name) {
		return types.Tag{}, types.ErrInvalidTag
	}
	tag, err := u.tags.FindByIDAndUserID(ctx, tagID, userID)
	if err != nil {
		return types.Tag{}, err
	}
	if _, err := u.lists.FindByIDAndUserID(ctx, update.ListID, userID); err != nil {
		return types.Tag{}, err
	}
	tag.ListID, tag.Name, tag.Color = update.ListID, update.Name, update.Color
	if err := u.tags.Update(ctx, tagID, userID, &tag); err != nil {
		return types.Tag{}, err
	}
	return tag, nil
}

type remove struct{ *repositories }

func (u remove) Execute(ctx context.Context, tagID, userID int64) error {
	if tagID <= 0 || userID <= 0 {
		return types.ErrInvalidTag
	}
	return u.tags.Delete(ctx, tagID, userID)
}

type byList struct{ *repositories }

func (u byList) Execute(ctx context.Context, listID, userID int64) ([]types.Tag, error) {
	if listID <= 0 || userID <= 0 {
		return nil, types.ErrInvalidTag
	}
	if _, err := u.lists.FindByIDAndUserID(ctx, listID, userID); err != nil {
		return nil, err
	}
	return u.tags.FindByListIDAndUserID(ctx, listID, userID)
}

type byID struct{ *repositories }

func (u byID) Execute(ctx context.Context, tagID, userID int64) (types.Tag, error) {
	if tagID <= 0 || userID <= 0 {
		return types.Tag{}, types.ErrInvalidTag
	}
	return u.tags.FindByIDAndUserID(ctx, tagID, userID)
}
func validName(name string) bool { n := utf8.RuneCountInString(name); return n > 0 && n <= 50 }
