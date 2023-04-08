package tasks

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
)

const TypeHello = "hello"

type HelloPayload struct {
}

func NewHelloTask() (*asynq.Task, error) {
	return asynq.NewTask(TypeHello, []byte("{}"), asynq.MaxRetry(10), asynq.Timeout(20*time.Hour), asynq.Retention(30*time.Minute)), nil
}

func HelloTask(ctx context.Context, t *asynq.Task) error {
	return nil
}
