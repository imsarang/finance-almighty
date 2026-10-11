package queue

import "context"

type Task struct {
	ID string
	Task_type string
	Queue string
	Args []any
	Kwargs map[string]any
}

type Enqueuer interface {
	Enqueue(ctx context.Context, t Task)(taskId string, err error)
}

type Status string

const (
	StatusPending Status = "PENDING"
	StatusSuccess Status = "SUCCESS"
	StatusFailure Status = "FAILURE"
)

type Result struct {
	TaskId string
	Status Status
	Result any `jsonschema:"task result returned by the worker"`
}

type ResultReader interface {
	Result(ctx context.Context, taskID string)(Result, error)
}