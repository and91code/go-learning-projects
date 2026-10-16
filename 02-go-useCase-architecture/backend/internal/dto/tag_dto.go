package dto

import "github.com/seu-usuario/taskflow-backend/internal/types"

type CreateTagRequest struct {
	ListID int64  `json:"list_id" binding:"required,gt=0"`
	Name   string `json:"name" binding:"required,min=1,max=50"`
	Color  string `json:"color" binding:"omitempty,max=7"`
}

type UpdateTagRequest struct {
	ListID int64  `json:"list_id" binding:"required,gt=0"`
	Name   string `json:"name" binding:"required,min=1,max=50"`
	Color  string `json:"color" binding:"omitempty,max=7"`
}

type TagResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type TagDetailResponse struct {
	ID     int64  `json:"id"`
	ListID int64  `json:"list_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}

func (req CreateTagRequest) ToDomain() types.Tag {
	return types.Tag{
		ListID: req.ListID,
		Name:   req.Name,
		Color:  req.Color,
	}
}

func (req UpdateTagRequest) ToDomain() types.Tag {
	return types.Tag{
		ListID: req.ListID,
		Name:   req.Name,
		Color:  req.Color,
	}
}

func TagResponseFromDomain(tag types.Tag) TagResponse {
	return TagResponse{
		ID:    tag.ID,
		Name:  tag.Name,
		Color: tag.Color,
	}
}

func TagDetailResponseFromDomain(tag types.Tag) TagDetailResponse {
	return TagDetailResponse{
		ID:     tag.ID,
		ListID: tag.ListID,
		Name:   tag.Name,
		Color:  tag.Color,
	}
}

func TagResponsesFromDomain(tags []types.Tag) []TagResponse {
	responses := make([]TagResponse, 0, len(tags))
	for _, tag := range tags {
		responses = append(responses, TagResponseFromDomain(tag))
	}
	return responses
}
