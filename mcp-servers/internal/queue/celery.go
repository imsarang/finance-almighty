package queue

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CeleryPublisher struct {
	client *redis.Client
}

//  this to check in compile time whether CeleryPublisher contains all the methods listed by the Enqueuer Interface
var _ Enqueuer = (*CeleryPublisher)(nil) // compile time check

func NewCeleryPublisher(client *redis.Client) *CeleryPublisher{
	return &CeleryPublisher{client: client}
}

func (p *CeleryPublisher) Enqueue(ctx context.Context, task Task)(string ,error) {
	// enqueue task to celery
	if task.Task_type == "" {
		return "", errors.New("task type is required")
	}

	id := task.ID
	if id == "" {
		id = uuid.NewString()
	}

	queueName := task.Queue
	if queueName == "" {
		queueName = "celery"
	}

	args := task.Args
	if args == nil {
		args = []any{}
	}

	kwargs := task.Kwargs
	if kwargs == nil {
		kwargs = map[string]any{}
	}

	// Celery protocol v2 body: [args, kwargs, embed]
	body, err := json.Marshal([]any{
		args,
		kwargs,
		map[string]any{"callbacks": nil, "errbacks": nil, "chain": nil, "chord": nil},
	})
	if err != nil {
		return "", fmt.Errorf("encode body: %w", err)
	}

	host, _ := os.Hostname()
	envelope, err := json.Marshal(map[string]any{
		"body":             base64.StdEncoding.EncodeToString(body),
		"content-encoding": "utf-8",
		"content-type":     "application/json",
		"headers": map[string]any{
			"lang":      "py",
			"task":      task.Task_type,
			"id":        id,
			"root_id":   id,
			"parent_id": nil,
			"group":     nil,
			"retries":   0,
			"eta":       nil,
			"expires":   nil,
			"timelimit": []any{nil, nil},
			"origin":    fmt.Sprintf("gen%d@%s", os.Getpid(), host),
		},
		"properties": map[string]any{
			"correlation_id": id,
			"reply_to":       "",
			"delivery_mode":  2,
			"delivery_info":  map[string]any{"exchange": "", "routing_key": queueName},
			"priority":       0,
			"body_encoding":  "base64",
			"delivery_tag":   uuid.NewString(),
		},
	})
	if err != nil {
		return "", fmt.Errorf("encode envelope: %w", err)
	}

	if err := p.client.LPush(ctx, queueName, envelope).Err(); err != nil {
		return "", fmt.Errorf("push to queue %q: %w", queueName, err)
	}

	return id, nil
}

