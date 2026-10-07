package dto

import "time"

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

// FormatListResponse converte uma Entidade de Domínio para o DTO de resposta
func FormatListResponse(id, userID int64, title, description string, createdAt time.Time) ListResponse {
	return ListResponse{
		ID:          id,
		UserID:      userID,
		Title:       title,
		Description: description,
		CreatedAt:   createdAt.Format(time.RFC3339),
	}
}
