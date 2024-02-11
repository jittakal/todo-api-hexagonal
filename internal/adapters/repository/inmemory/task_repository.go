package inmemory

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/jittakal/todo-api-hexagonal/internal/domain"
	"github.com/jittakal/todo-api-hexagonal/internal/errors"
)

type TaskRepository struct {
	tasks map[string]*domain.Task
	mu    sync.RWMutex
}

var _ domain.TaskRepository = (*TaskRepository)(nil)

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		tasks: make(map[string]*domain.Task),
	}
}

// Create handles the creation of a new task.
func (r *TaskRepository) Create(ctx context.Context, taskTitle string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task := &domain.Task{
		ID:    uuid.New().String(),
		Title: taskTitle,
	}
	r.tasks[task.ID] = task

	return task.ID, nil
}

// Update handles the updation of an existing task.
func (r *TaskRepository) Update(ctx context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tasks[task.ID] = task
	return nil
}

// Delete handles the deletion of an existing task.
func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.tasks, id)
	return nil
}

// GetByID return task details by provided ID.
func (r *TaskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, ok := r.tasks[id]
	if !ok {
		return &domain.Task{}, errors.ErrTaskNotFound
	}
	return task, nil
}

// GetByID return all the task details.
func (r *TaskRepository) GetAll(ctx context.Context) ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}
