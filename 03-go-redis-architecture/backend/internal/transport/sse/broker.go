package sse

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/seu-usuario/taskflow-backend/internal/logger"
)

type Broker struct {
	client  *redis.Client
	mu      sync.RWMutex
	nextID  uint64
	clients map[int64]map[uint64]chan []byte
}

func NewBroker(client *redis.Client) *Broker {
	return &Broker{client: client, clients: make(map[int64]map[uint64]chan []byte)}
}

func (b *Broker) Subscribe(userID int64) (<-chan []byte, func(), error) {
	if b == nil || b.client == nil {
		return nil, nil, fmt.Errorf("SSE broker is not configured")
	}
	if userID <= 0 {
		return nil, nil, fmt.Errorf("user ID must be positive")
	}

	b.mu.Lock()
	b.nextID++
	id := b.nextID
	events := make(chan []byte, 16)
	if b.clients[userID] == nil {
		b.clients[userID] = make(map[uint64]chan []byte)
	}
	b.clients[userID][id] = events
	b.mu.Unlock()

	var once sync.Once
	unsubscribe := func() {
		once.Do(func() {
			b.mu.Lock()
			delete(b.clients[userID], id)
			if len(b.clients[userID]) == 0 {
				delete(b.clients, userID)
			}
			b.mu.Unlock()
			close(events)
		})
	}
	return events, unsubscribe, nil
}

func (b *Broker) Run(ctx context.Context) error {
	if b == nil || b.client == nil {
		return fmt.Errorf("SSE broker is not configured")
	}
	sub := b.client.PSubscribe(ctx, "user:*:events")
	defer sub.Close()
	if _, err := sub.Receive(ctx); err != nil {
		return fmt.Errorf("subscribe to user events: %w", err)
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case message, ok := <-sub.Channel():
			if !ok {
				if ctx.Err() != nil {
					return nil
				}
				return fmt.Errorf("Redis event subscription closed")
			}
			userID, err := userIDFromChannel(message.Channel)
			if err != nil {
				return err
			}
			b.broadcast(userID, []byte(message.Payload))
		}
	}
}

func (b *Broker) broadcast(userID int64, payload []byte) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, events := range b.clients[userID] {
		select {
		case events <- payload:
		default:
			logger.Warnf("dropping SSE event for slow client of user %d", userID)
		}
	}
}

func userIDFromChannel(channel string) (int64, error) {
	parts := strings.Split(channel, ":")
	if len(parts) != 3 || parts[0] != "user" || parts[2] != "events" {
		return 0, fmt.Errorf("invalid Redis event channel %q", channel)
	}
	userID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || userID <= 0 {
		return 0, fmt.Errorf("invalid user ID in Redis event channel %q", channel)
	}
	return userID, nil
}
