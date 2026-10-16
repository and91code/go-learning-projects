package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/redis/go-redis/v9"
	"github.com/seu-usuario/taskflow-backend/internal/types"
)

const QueueKey = "taskflow:queue"

type Job struct {
	Action   string          `json:"action"`
	UserID   int64           `json:"user_id"`
	EntityID int64           `json:"entity_id,omitempty"`
	Payload  json.RawMessage `json:"payload,omitempty"`
	Attempts int             `json:"attempts,omitempty"`
}

type Producer interface {
	Enqueue(context.Context, Job) error
}

type RedisProducer struct{ client *redis.Client }

func NewProducer(client *redis.Client) *RedisProducer { return &RedisProducer{client: client} }

func (p *RedisProducer) Enqueue(ctx context.Context, job Job) error {
	if p == nil || p.client == nil {
		return fmt.Errorf("queue producer is not configured")
	}
	body, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("encode queue job: %w", err)
	}
	if err := p.client.LPush(ctx, QueueKey, body).Err(); err != nil {
		return fmt.Errorf("enqueue job: %w", err)
	}
	return nil
}

func ValidateTask(action string, userID, entityID int64, task types.Task) error {
	titleLength := runeLength(strings.TrimSpace(task.Title))
	if userID <= 0 || (action == "TASK_CREATE" && task.ListID <= 0) ||
		(action == "TASK_UPDATE" && entityID <= 0) ||
		titleLength < 3 || titleLength > 150 ||
		(task.Priority != types.PriorityLow && task.Priority != types.PriorityMedium && task.Priority != types.PriorityHigh) {
		return types.ErrInvalidTask
	}
	if action == "TASK_UPDATE" && task.Status != types.StatusPending && task.Status != types.StatusInProgress && task.Status != types.StatusCompleted {
		return types.ErrInvalidTask
	}
	seen := make(map[int64]struct{}, len(task.Tags))
	for _, tag := range task.Tags {
		if tag.ID <= 0 {
			return types.ErrInvalidTask
		}
		if _, exists := seen[tag.ID]; exists {
			return types.ErrInvalidTask
		}
		seen[tag.ID] = struct{}{}
	}
	return nil
}

func ValidateList(userID, entityID int64, list types.List) error {
	titleLength := runeLength(strings.TrimSpace(list.Title))
	if userID <= 0 || entityID < 0 || titleLength < 3 || titleLength > 100 || runeLength(list.Description) > 255 {
		return fmt.Errorf("invalid list")
	}
	return nil
}

func ValidateTag(action string, userID, entityID int64, tag types.Tag) error {
	nameLength := runeLength(strings.TrimSpace(tag.Name))
	if userID <= 0 || tag.ListID <= 0 || nameLength == 0 || nameLength > 50 || runeLength(tag.Color) > 7 ||
		(action == "TAG_UPDATE" && entityID <= 0) {
		return types.ErrInvalidTag
	}
	return nil
}

func runeLength(value string) int { return utf8.RuneCountInString(value) }
