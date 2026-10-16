package pubsub

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Broadcaster struct{ client *redis.Client }

func NewBroadcaster(client *redis.Client) *Broadcaster {
	return &Broadcaster{client: client}
}

func (b *Broadcaster) Broadcast(ctx context.Context, userID int64, event []byte) error {
	if b == nil || b.client == nil {
		return fmt.Errorf("pubsub broadcaster is not configured")
	}
	if userID <= 0 {
		return fmt.Errorf("user ID must be positive")
	}
	if err := b.client.Publish(ctx, fmt.Sprintf("user:%d:events", userID), event).Err(); err != nil {
		return fmt.Errorf("publish event: %w", err)
	}
	return nil
}
