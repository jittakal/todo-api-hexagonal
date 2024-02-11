package domain

import (
	"context"
)

type TaskRepository interface {
	Create(ctx context.Context, taskTitle string) (string, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Task, error)
	GetAll(ctx context.Context) ([]*Task, error)
}
