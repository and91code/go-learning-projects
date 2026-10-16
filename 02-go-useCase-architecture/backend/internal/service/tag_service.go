package service

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/seu-usuario/taskflow-backend/internal/types"
)

type tagService struct {
	tagRepository  types.TagRepository
	listRepository types.ListRepository
}

func NewTagService(tagRepository types.TagRepository, listRepository types.ListRepository) types.TagService {
	return &tagService{
		tagRepository:  tagRepository,
		listRepository: listRepository,
	}
}

var _ types.TagService = (*tagService)(nil)

func (s *tagService) CreateTag(ctx context.Context, userID int64, tag types.Tag) (types.Tag, error) {
	tag.Name = strings.TrimSpace(tag.Name)
	if userID <= 0 || tag.ListID <= 0 || !validTagName(tag.Name) {
		return types.Tag{}, types.ErrInvalidTag
	}

	if _, err := s.listRepository.FindByIDAndUserID(ctx, tag.ListID, userID); err != nil {
		return types.Tag{}, err
	}

	if err := s.tagRepository.Create(ctx, userID, &tag); err != nil {
		return types.Tag{}, err
	}
	return tag, nil
}

func (s *tagService) UpdateTag(ctx context.Context, tagID, userID int64, update types.Tag) (types.Tag, error) {
	update.Name = strings.TrimSpace(update.Name)
	if tagID <= 0 || userID <= 0 || update.ListID <= 0 || !validTagName(update.Name) {
		return types.Tag{}, types.ErrInvalidTag
	}

	tag, err := s.tagRepository.FindByIDAndUserID(ctx, tagID, userID)
	if err != nil {
		return types.Tag{}, err
	}
	if _, err := s.listRepository.FindByIDAndUserID(ctx, update.ListID, userID); err != nil {
		return types.Tag{}, err
	}

	tag.ListID = update.ListID
	tag.Name = update.Name
	tag.Color = update.Color
	if err := s.tagRepository.Update(ctx, tagID, userID, &tag); err != nil {
		return types.Tag{}, err
	}
	return tag, nil
}

func validTagName(name string) bool {
	length := utf8.RuneCountInString(name)
	return length > 0 && length <= 50
}

func (s *tagService) DeleteTag(ctx context.Context, tagID, userID int64) error {
	if tagID <= 0 || userID <= 0 {
		return types.ErrInvalidTag
	}
	return s.tagRepository.Delete(ctx, tagID, userID)
}

func (s *tagService) GetTagsByList(ctx context.Context, listID, userID int64) ([]types.Tag, error) {
	if listID <= 0 || userID <= 0 {
		return nil, types.ErrInvalidTag
	}
	if _, err := s.listRepository.FindByIDAndUserID(ctx, listID, userID); err != nil {
		return nil, err
	}
	return s.tagRepository.FindByListIDAndUserID(ctx, listID, userID)
}

func (s *tagService) GetTagById(ctx context.Context, tagID, userID int64) (types.Tag, error) {
	if tagID <= 0 || userID <= 0 {
		return types.Tag{}, types.ErrInvalidTag
	}
	return s.tagRepository.FindByIDAndUserID(ctx, tagID, userID)
}
