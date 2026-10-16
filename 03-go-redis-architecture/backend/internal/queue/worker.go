package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/seu-usuario/taskflow-backend/internal/logger"
	"github.com/seu-usuario/taskflow-backend/internal/pubsub"
	"github.com/seu-usuario/taskflow-backend/internal/types"
)

const (
	processingKey = "taskflow:processing"
	deadLetterKey = "taskflow:dead-letter"
	maxAttempts   = 5
)

type Dispatcher interface {
	Dispatch(context.Context, Job) (int64, error)
}

type Worker struct {
	client      *redis.Client
	dispatcher  Dispatcher
	broadcaster interface {
		Broadcast(context.Context, int64, []byte) error
	}
}

func NewWorker(client *redis.Client, dispatcher Dispatcher, broadcaster *pubsub.Broadcaster) *Worker {
	return &Worker{client: client, dispatcher: dispatcher, broadcaster: broadcaster}
}

type permanentError struct{ err error }

func (e permanentError) Error() string { return e.err.Error() }
func (e permanentError) Unwrap() error { return e.err }
func permanent(err error) error        { return permanentError{err: err} }

func isPermanent(err error) bool {
	var target permanentError
	return errors.As(err, &target) ||
		errors.Is(err, types.ErrInvalidTask) ||
		errors.Is(err, types.ErrInvalidTag) ||
		errors.Is(err, types.ErrTaskNotFound) ||
		errors.Is(err, types.ErrTagNotFound) ||
		errors.Is(err, types.ErrListNotFound)
}

func (w *Worker) Run(ctx context.Context) {
	if w.client == nil || w.dispatcher == nil {
		logger.Errorf("queue worker requires Redis client and dispatcher")
		return
	}
	w.recoverProcessing(ctx)
	for {
		raw, err := w.client.BRPopLPush(ctx, QueueKey, processingKey, 0).Result()
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			logger.Errorf("queue pop failed: %v", err)
			if !sleep(ctx, time.Second) {
				return
			}
			continue
		}
		var job Job
		if err := json.Unmarshal([]byte(raw), &job); err != nil {
			w.deadLetter(ctx, raw, fmt.Errorf("decode job: %w", err))
			continue
		}
		if err := validateJob(job); err != nil {
			w.deadLetter(ctx, raw, err)
			continue
		}
		entityID, err := w.dispatcher.Dispatch(ctx, job)
		if err != nil {
			logger.Errorf("job %s for user %d failed: %v", job.Action, job.UserID, err)
			if isPermanent(err) || job.Attempts+1 >= maxAttempts {
				w.deadLetter(ctx, raw, err)
				continue
			}
			job.Attempts++
			retry, marshalErr := json.Marshal(job)
			if marshalErr != nil {
				w.deadLetter(ctx, raw, marshalErr)
				continue
			}
			if !sleep(ctx, time.Second) {
				return
			}
			if err := w.client.LPush(ctx, QueueKey, retry).Err(); err != nil {
				logger.Errorf("requeue failed job: %v", err)
				continue
			}
			w.ack(ctx, raw)
			continue
		}

		resource, action := notificationAction(job.Action)
		notification := map[string]any{
			"id":            fmt.Sprintf("%d-%d", time.Now().UnixNano(), entityID),
			"message":       notificationMessage(resource, action),
			"resource_type": resource, "resource_id": entityID,
			"action": action, "timestamp": time.Now().UTC().Format(time.RFC3339Nano), "unread": true,
		}
		if w.client != nil {
			body, _ := json.Marshal(notification)
			key := fmt.Sprintf("taskflow:user:%d:notifications", job.UserID)
			pipe := w.client.TxPipeline()
			pipe.HSet(ctx, key, notification["id"], body)
			pipe.ZAdd(ctx, key+":index", redis.Z{Score: float64(time.Now().UnixNano()), Member: notification["id"]})
			if _, err := pipe.Exec(ctx); err != nil {
				logger.Errorf("persist notification for user %d: %v", job.UserID, err)
			}
		}
		event, _ := json.Marshal(notification)
		if err := w.broadcaster.Broadcast(ctx, job.UserID, event); err != nil {
			logger.Errorf("publish completion event for user %d: %v", job.UserID, err)
		}
		w.ack(ctx, raw)
	}
}

func validateJob(job Job) error {
	if job.UserID <= 0 {
		return permanent(fmt.Errorf("job user_id must be positive"))
	}
	switch job.Action {
	case "TASK_CREATE", "LIST_CREATE", "TAG_CREATE":
		if len(job.Payload) == 0 {
			return permanent(fmt.Errorf("job payload is required"))
		}
	case "TASK_UPDATE", "TASK_DELETE", "LIST_UPDATE", "LIST_DELETE", "TAG_UPDATE", "TAG_DELETE":
		if job.EntityID <= 0 {
			return permanent(fmt.Errorf("job entity_id must be positive"))
		}
		if job.Action == "TASK_UPDATE" || job.Action == "LIST_UPDATE" || job.Action == "TAG_UPDATE" {
			if len(job.Payload) == 0 {
				return permanent(fmt.Errorf("job payload is required"))
			}
		}
	default:
		return permanent(fmt.Errorf("unknown action %q", job.Action))
	}
	return nil
}

func (w *Worker) recoverProcessing(ctx context.Context) {
	for ctx.Err() == nil {
		raw, err := w.client.RPop(ctx, processingKey).Result()
		if errors.Is(err, redis.Nil) {
			return
		}
		if err != nil {
			logger.Errorf("recover processing queue: %v", err)
			return
		}
		if err := w.client.LPush(ctx, QueueKey, raw).Err(); err != nil {
			logger.Errorf("restore interrupted job: %v", err)
			_ = w.client.LPush(ctx, processingKey, raw).Err()
			return
		}
	}
}

func (w *Worker) deadLetter(ctx context.Context, raw string, cause error) {
	logger.Errorf("discarding queue job: %v", cause)
	if err := w.client.LPush(ctx, deadLetterKey, raw).Err(); err != nil {
		logger.Errorf("store dead-letter job: %v", err)
		return
	}
	w.ack(ctx, raw)
}

func (w *Worker) ack(ctx context.Context, raw string) {
	if err := w.client.LRem(ctx, processingKey, 1, raw).Err(); err != nil {
		logger.Errorf("acknowledge queue job: %v", err)
	}
}

func sleep(ctx context.Context, delay time.Duration) bool {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

type DispatcherFuncs struct {
	CreateTask interface {
		Execute(context.Context, int64, types.Task) (types.Task, error)
	}
	UpdateTask interface {
		Execute(context.Context, int64, int64, types.Task) (types.Task, error)
	}
	DeleteTask interface {
		Execute(context.Context, int64, int64) error
	}
	CreateList interface {
		Execute(context.Context, int64, types.List) (types.List, error)
	}
	UpdateList interface {
		Execute(context.Context, int64, int64, types.List) (types.List, error)
	}
	DeleteList interface {
		Execute(context.Context, int64, int64) error
	}
	CreateTag interface {
		Execute(context.Context, int64, types.Tag) (types.Tag, error)
	}
	UpdateTag interface {
		Execute(context.Context, int64, int64, types.Tag) (types.Tag, error)
	}
	DeleteTag interface {
		Execute(context.Context, int64, int64) error
	}
}

func notificationAction(action string) (string, string) {
	parts := strings.SplitN(action, "_", 2)
	if len(parts) != 2 {
		return "task", "created"
	}
	verb := map[string]string{"CREATE": "created", "UPDATE": "updated", "DELETE": "deleted"}[parts[1]]
	return strings.ToLower(parts[0]), verb
}

func notificationMessage(resource, action string) string {
	noun := map[string]string{"list": "Lista", "task": "Tarefa", "tag": "Tag"}[resource]
	verb := map[string]string{"created": "criada", "updated": "atualizada", "deleted": "excluída"}[action]
	return noun + " " + verb + " com sucesso."
}

func (d DispatcherFuncs) Dispatch(ctx context.Context, job Job) (int64, error) {
	decode := func(target any) error {
		if err := json.Unmarshal(job.Payload, target); err != nil {
			return permanent(fmt.Errorf("decode %s payload: %w", job.Action, err))
		}
		return nil
	}
	switch job.Action {
	case "TASK_CREATE":
		var v types.Task
		if err := decode(&v); err != nil {
			return 0, err
		}
		result, err := d.CreateTask.Execute(ctx, job.UserID, v)
		return result.ID, err
	case "TASK_UPDATE":
		var v types.Task
		if err := decode(&v); err != nil {
			return 0, err
		}
		result, err := d.UpdateTask.Execute(ctx, job.EntityID, job.UserID, v)
		return result.ID, err
	case "TASK_DELETE":
		return job.EntityID, d.DeleteTask.Execute(ctx, job.EntityID, job.UserID)
	case "LIST_CREATE":
		var v types.List
		if err := decode(&v); err != nil {
			return 0, err
		}
		result, err := d.CreateList.Execute(ctx, job.UserID, v)
		return result.ID, err
	case "LIST_UPDATE":
		var v types.List
		if err := decode(&v); err != nil {
			return 0, err
		}
		result, err := d.UpdateList.Execute(ctx, job.EntityID, job.UserID, v)
		return result.ID, err
	case "LIST_DELETE":
		return job.EntityID, d.DeleteList.Execute(ctx, job.EntityID, job.UserID)
	case "TAG_CREATE":
		var v types.Tag
		if err := decode(&v); err != nil {
			return 0, err
		}
		result, err := d.CreateTag.Execute(ctx, job.UserID, v)
		return result.ID, err
	case "TAG_UPDATE":
		var v types.Tag
		if err := decode(&v); err != nil {
			return 0, err
		}
		result, err := d.UpdateTag.Execute(ctx, job.EntityID, job.UserID, v)
		return result.ID, err
	case "TAG_DELETE":
		return job.EntityID, d.DeleteTag.Execute(ctx, job.EntityID, job.UserID)
	default:
		return 0, permanent(fmt.Errorf("unknown action %q", job.Action))
	}
}
