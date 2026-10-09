package dto

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/seu-usuario/taskflow-backend/internal/types"
)

type TaskDateTime time.Time

func (t *TaskDateTime) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("due_date deve ser uma string: %w", err)
	}

	for _, layout := range []string{time.RFC3339Nano, "2006-01-02 15:04:05"} {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			*t = TaskDateTime(parsed)
			return nil
		}
	}

	return fmt.Errorf("due_date deve estar no formato RFC3339 ou YYYY-MM-DD HH:MM:SS")
}

type CreateTaskRequest struct {
	ListID      int64         `json:"list_id" binding:"required,gt=0"`
	Title       string        `json:"title" binding:"required,min=3,max=150"`
	Description string        `json:"description" binding:"omitempty,max=65535"`
	Priority    string        `json:"priority" binding:"required,oneof=low medium high"`
	DueDate     *TaskDateTime `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       string        `json:"title" binding:"required,min=3,max=150"`
	Description string        `json:"description" binding:"omitempty,max=65535"`
	Status      string        `json:"status" binding:"required,oneof=pending in_progress completed"`
	Priority    string        `json:"priority" binding:"required,oneof=low medium high"`
	DueDate     *TaskDateTime `json:"due_date"`
}

type TaskResponse struct {
	ID          int64         `json:"id"`
	ListID      int64         `json:"list_id"`
	Title       string        `json:"title"`
	Description string        `json:"description,omitempty"`
	Status      string        `json:"status"`
	Priority    string        `json:"priority"`
	DueDate     *time.Time    `json:"due_date,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Tags        []TagResponse `json:"tags,omitempty"`
}

func (req CreateTaskRequest) ToDomain() types.Task {
	return types.Task{
		ListID:      req.ListID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    types.TaskPriority(req.Priority),
		DueDate:     taskDateTimeToTime(req.DueDate),
	}
}

func (req UpdateTaskRequest) ToDomain() types.Task {
	return types.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      types.TaskStatus(req.Status),
		Priority:    types.TaskPriority(req.Priority),
		DueDate:     taskDateTimeToTime(req.DueDate),
	}
}

func taskDateTimeToTime(value *TaskDateTime) *time.Time {
	if value == nil {
		return nil
	}
	converted := time.Time(*value)
	return &converted
}

func TaskResponseFromDomain(task types.Task) TaskResponse {
	response := TaskResponse{
		ID:          task.ID,
		ListID:      task.ListID,
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		Priority:    string(task.Priority),
		DueDate:     task.DueDate,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
	if len(task.Tags) > 0 {
		response.Tags = make([]TagResponse, 0, len(task.Tags))
		for _, tag := range task.Tags {
			response.Tags = append(response.Tags, TagResponse{ID: tag.ID, Name: tag.Name, Color: tag.Color})
		}
	}
	return response
}

func TaskResponsesFromDomain(tasks []types.Task) []TaskResponse {
	responses := make([]TaskResponse, 0, len(tasks))
	for _, task := range tasks {
		responses = append(responses, TaskResponseFromDomain(task))
	}
	return responses
}
