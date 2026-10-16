package task

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/seu-usuario/taskflow-backend/internal/types"
)

type Create interface {
	Execute(context.Context, int64, types.Task) (types.Task, error)
}
type Update interface {
	Execute(context.Context, int64, int64, types.Task) (types.Task, error)
}
type Delete interface {
	Execute(context.Context, int64, int64) error
}
type GetByList interface {
	Execute(context.Context, int64, int64) ([]types.Task, error)
}
type GetByID interface {
	Execute(context.Context, int64, int64) (types.Task, error)
}
type repositories struct {
	tasks types.TaskRepository
	lists types.ListRepository
}

func New(tasks types.TaskRepository, lists types.ListRepository) (Create, Update, Delete, GetByList, GetByID) {
	r := &repositories{tasks: tasks, lists: lists}
	return create{r}, update{r}, remove{r}, byList{r}, byID{r}
}

type create struct{ *repositories }

func (u create) Execute(ctx context.Context, userID int64, task types.Task) (types.Task, error) {
	task.Title = strings.TrimSpace(task.Title)
	if userID <= 0 || task.ListID <= 0 || !validTitle(task.Title) || !validPriority(task.Priority) || !validTags(task.Tags) {
		return types.Task{}, types.ErrInvalidTask
	}
	if _, err := u.lists.FindByIDAndUserID(ctx, task.ListID, userID); err != nil {
		return types.Task{}, err
	}
	task.Status = types.StatusPending
	task.CreatedAt = time.Now()
	task.UpdatedAt = task.CreatedAt
	if err := u.tasks.Create(ctx, userID, &task); err != nil {
		return types.Task{}, err
	}
	return task, nil
}

type update struct{ *repositories }

func (u update) Execute(ctx context.Context, taskID, userID int64, update types.Task) (types.Task, error) {
	update.Title = strings.TrimSpace(update.Title)
	if taskID <= 0 || userID <= 0 || !validTitle(update.Title) || !validPriority(update.Priority) || !validStatus(update.Status) || !validTags(update.Tags) {
		return types.Task{}, types.ErrInvalidTask
	}
	task, err := u.tasks.FindByIDAndUserID(ctx, taskID, userID)
	if err != nil {
		return types.Task{}, err
	}
	task.Title = update.Title
	task.Description = strings.TrimSpace(update.Description)
	task.Status = update.Status
	task.Priority = update.Priority
	task.DueDate = update.DueDate
	if update.Tags != nil {
		task.Tags = update.Tags
	}
	task.UpdatedAt = time.Now()
	if err := u.tasks.Update(ctx, taskID, userID, task); err != nil {
		return types.Task{}, err
	}
	return *task, nil
}

type remove struct{ *repositories }

func (u remove) Execute(ctx context.Context, taskID, userID int64) error {
	if taskID <= 0 || userID <= 0 {
		return types.ErrInvalidTask
	}
	return u.tasks.Delete(ctx, taskID, userID)
}

type byList struct{ *repositories }

func (u byList) Execute(ctx context.Context, listID, userID int64) ([]types.Task, error) {
	if listID <= 0 || userID <= 0 {
		return nil, types.ErrInvalidTask
	}
	return u.tasks.FindByListIDAndUserID(ctx, listID, userID)
}

type byID struct{ *repositories }

func (u byID) Execute(ctx context.Context, taskID, userID int64) (types.Task, error) {
	if taskID <= 0 || userID <= 0 {
		return types.Task{}, types.ErrInvalidTask
	}
	task, err := u.tasks.FindByIDAndUserID(ctx, taskID, userID)
	if err != nil {
		return types.Task{}, err
	}
	return *task, nil
}
func validTitle(title string) bool { n := utf8.RuneCountInString(title); return n >= 3 && n <= 150 }
func validPriority(p types.TaskPriority) bool {
	return p == types.PriorityLow || p == types.PriorityMedium || p == types.PriorityHigh
}
func validStatus(s types.TaskStatus) bool {
	return s == types.StatusPending || s == types.StatusInProgress || s == types.StatusCompleted
}
func validTags(tags []types.Tag) bool {
	seen := make(map[int64]struct{}, len(tags))
	for _, tag := range tags {
		if tag.ID <= 0 {
			return false
		}
		if _, ok := seen[tag.ID]; ok {
			return false
		}
		seen[tag.ID] = struct{}{}
	}
	return true
}
