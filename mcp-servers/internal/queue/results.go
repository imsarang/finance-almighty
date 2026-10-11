package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type ReadResultReader struct {
	client *redis.Client
}

var _ ResultReader = (*ReadResultReader)(nil)

func NewRedisResultReader(c *redis.Client) *ReadResultReader {
	return &ReadResultReader{client: c}
}

func (p *ReadResultReader) Result(ctx context.Context, id string)(Result, error){
	// get task status from the queue
	raw, err := p.client.Get(ctx, "celery-task-meta-"+id).Bytes()
	if errors.Is(err, redis.Nil) {
		return Result{TaskId: id, Status: StatusPending}, nil
	}
	if err != nil {
		return Result{}, fmt.Errorf("get result: %w", err)
	}

	var res Result
	if err := json.Unmarshal(raw, &res); err != nil {
		return Result{}, fmt.Errorf("decode result: %w", err)
	}
	res.TaskId = id

	return res, nil
}
