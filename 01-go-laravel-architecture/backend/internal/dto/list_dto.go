package dto

import (
	"time"

	"github.com/seu-usuario/taskflow-backend/internal/types"
)

// CreateListRequest representa o payload para criação de lista (Form Request)
type CreateListRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"omitempty,max=255"`
}

// UpdateListRequest representa o payload de atualização de lista
type UpdateListRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"omitempty,max=255"`
}

// ListResponse representa a resposta formatada da lista (API Resource)
type ListResponse struct {
	ID          int64          `json:"id"`
	UserID      int64          `json:"user_id"`
	Title       string         `json:"title"`
	Description string         `json:"description,omitempty"`
	CreatedAt   string         `json:"created_at"`
	Tasks       []TaskResponse `json:"tasks,omitempty"`
	Tags        []TagResponse  `json:"tags,omitempty"`
}

// TagResponse representa o DTO simplificado da Tag
type TagResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (req CreateListRequest) ToDomain() types.List {
	return types.List{
		Title:       req.Title,
		Description: req.Description,
	}
}

func (req UpdateListRequest) ToDomain() types.List {
	return types.List{
		Title:       req.Title,
		Description: req.Description,
	}
}

func ListResponseFromDomain(list types.List) ListResponse {
	response := ListResponse{
		ID:          list.ID,
		UserID:      list.UserID,
		Title:       list.Title,
		Description: list.Description,
		CreatedAt:   list.CreatedAt.Format(time.RFC3339),
	}
	if len(list.Tasks) > 0 {
		response.Tasks = make([]TaskResponse, 0, len(list.Tasks))
		for _, task := range list.Tasks {
			response.Tasks = append(response.Tasks, TaskResponseFromDomain(task))
		}
	}
	if len(list.Tags) > 0 {
		response.Tags = make([]TagResponse, 0, len(list.Tags))
		for _, tag := range list.Tags {
			response.Tags = append(response.Tags, TagResponse{ID: tag.ID, Name: tag.Name, Color: tag.Color})
		}
	}
	return response
}

func ListResponsesFromDomain(lists []types.List) []ListResponse {
	responses := make([]ListResponse, 0, len(lists))
	for _, list := range lists {
		responses = append(responses, ListResponseFromDomain(list))
	}
	return responses
}
