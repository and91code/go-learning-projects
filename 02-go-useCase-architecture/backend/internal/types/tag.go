package types

import "context"

type Tag struct {
	ID     int64
	ListID int64
	Name   string
	Color  string
}

type TagRepository interface {
	Create(ctx context.Context, userID int64, tag *Tag) error
	Update(ctx context.Context, tagID, userID int64, tag *Tag) error
	Delete(ctx context.Context, tagID int64, userID int64) error
	FindByListIDAndUserID(ctx context.Context, listID, userID int64) ([]Tag, error)
	FindByIDAndUserID(ctx context.Context, tagID, userID int64) (Tag, error)
}

type TagService interface {
	CreateTag(ctx context.Context, userID int64, tag Tag) (Tag, error)
	UpdateTag(ctx context.Context, tagID, userID int64, tag Tag) (Tag, error)
	DeleteTag(ctx context.Context, tagID int64, userID int64) error
	GetTagsByList(ctx context.Context, listID, userID int64) ([]Tag, error)
	GetTagById(ctx context.Context, tagID, userID int64) (Tag, error)
}
