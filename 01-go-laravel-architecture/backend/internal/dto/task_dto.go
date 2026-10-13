package dto

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/seu-usuario/taskflow-backend/internal/types"
)

type TaskDateTime time.Time

const taskDateTimeLayout = "2006-01-02 15:04"

func (t *TaskDateTime) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("due_date deve ser uma string: %w", err)
	}

	parsed, err := time.ParseInLocation(taskDateTimeLayout, value, time.UTC)
	if err != nil {
		return fmt.Errorf("due_date deve estar no formato YYYY-MM-DD HH:mm: %w", err)
	}
	*t = TaskDateTime(parsed)
	return nil
}

func (t TaskDateTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).UTC().Format(taskDateTimeLayout))
}

type CreateTaskRequest struct {
	ListID      int64         `json:"list_id" binding:"required,gt=0"`
	Title       string        `json:"title" binding:"required,min=3,max=150"`
	Description string        `json:"description" binding:"omitempty,max=65535"`
	Priority    string        `json:"priority" binding:"required,oneof=low medium high"`
	DueDate     *TaskDateTime `json:"due_date"`
	TagIDs      *[]int64      `json:"tag_ids"`
}

type UpdateTaskRequest struct {
	Title       string        `json:"title" binding:"required,min=3,max=150"`
	Description string        `json:"description" binding:"omitempty,max=65535"`
	Status      string        `json:"status" binding:"required,oneof=pending in_progress completed"`
	Priority    string        `json:"priority" binding:"required,oneof=low medium high"`
	DueDate     *TaskDateTime `json:"due_date"`
	TagIDs      *[]int64      `json:"tag_ids"`
}

type TaskResponse struct {
	ID          int64         `json:"id"`
	ListID      int64         `json:"list_id"`
	Title       string        `json:"title"`
	Description string        `json:"description,omitempty"`
	Status      string        `json:"status"`
	Priority    string        `json:"priority"`
	DueDate     *TaskDateTime `json:"due_date,omitempty"`
	CreatedAt   TaskDateTime  `json:"created_at"`
	UpdatedAt   TaskDateTime  `json:"updated_at"`
	Tags        []TagResponse `json:"tags"`
}

func (req CreateTaskRequest) ToDomain() types.Task {
	return types.Task{
		ListID:      req.ListID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    types.TaskPriority(req.Priority),
		DueDate:     taskDateTimeToTime(req.DueDate),
		Tags:        taskTagsFromIDs(req.TagIDs),
	}
}

func (req UpdateTaskRequest) ToDomain() types.Task {
	return types.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      types.TaskStatus(req.Status),
		Priority:    types.TaskPriority(req.Priority),
		DueDate:     taskDateTimeToTime(req.DueDate),
		Tags:        taskTagsFromIDs(req.TagIDs),
	}
}

func taskTagsFromIDs(tagIDs *[]int64) []types.Tag {
	if tagIDs == nil {
		return nil
	}
	tags := make([]types.Tag, 0, len(*tagIDs))
	for _, id := range *tagIDs {
		tags = append(tags, types.Tag{ID: id})
	}
	return tags
}

func taskDateTimeToTime(value *TaskDateTime) *time.Time {
	if value == nil {
		return nil
	}
	converted := time.Time(*value)
	return &converted
}

func taskTimeToDateTime(value time.Time) TaskDateTime {
	return TaskDateTime(value)
}

func TaskResponseFromDomain(task types.Task) TaskResponse {
	var dueDate *TaskDateTime
	if task.DueDate != nil {
		formattedDueDate := taskTimeToDateTime(*task.DueDate)
		dueDate = &formattedDueDate
	}

	response := TaskResponse{
		ID:          task.ID,
		ListID:      task.ListID,
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		Priority:    string(task.Priority),
		DueDate:     dueDate,
		CreatedAt:   taskTimeToDateTime(task.CreatedAt),
		UpdatedAt:   taskTimeToDateTime(task.UpdatedAt),
		Tags:        make([]TagResponse, 0, len(task.Tags)),
	}
	for _, tag := range task.Tags {
		response.Tags = append(response.Tags, TagResponseFromDomain(tag))
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
