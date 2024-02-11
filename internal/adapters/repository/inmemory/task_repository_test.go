package inmemory

import (
	"context"
	"reflect"
	"testing"

	"github.com/jittakal/todo-api-hexagonal/internal/domain"
	"github.com/jittakal/todo-api-hexagonal/internal/errors"
)

func TestTaskRepository_Create(t *testing.T) {
	repo := NewTaskRepository()

	taskID, err := repo.Create(context.Background(), "Task 1")
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	if len(taskID) == 0 {
		t.Error("Expected non-empty task ID")
	}

	// Verify task exists in the repository
	task, err := repo.GetByID(context.Background(), taskID)
	if err != nil {
		t.Fatalf("Failed to get task by ID: %v", err)
	}

	if task.Title != "Task 1" {
		t.Errorf("Expected task title 'Task 1', got '%s'", task.Title)
	}
}

func TestTaskRepository_Update(t *testing.T) {
	repo := NewTaskRepository()

	taskID, err := repo.Create(context.Background(), "Task 1")
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	task := &domain.Task{ID: taskID, Title: "Updated Task"}
	err = repo.Update(context.Background(), task)
	if err != nil {
		t.Fatalf("Failed to update task: %v", err)
	}

	// Verify task was updated
	updatedTask, err := repo.GetByID(context.Background(), taskID)
	if err != nil {
		t.Fatalf("Failed to get task by ID: %v", err)
	}

	if !reflect.DeepEqual(task, updatedTask) {
		t.Error("Expected updated task, got different task")
	}
}

func TestTaskRepository_Delete(t *testing.T) {
	repo := NewTaskRepository()

	taskID, err := repo.Create(context.Background(), "Task 1")
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	err = repo.Delete(context.Background(), taskID)
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}

	// Verify task is not present after deletion
	_, err = repo.GetByID(context.Background(), taskID)
	if err != errors.ErrTaskNotFound {
		t.Errorf("Expected task not found error, got: %v", err)
	}
}

func TestTaskRepository_GetByID(t *testing.T) {
	repo := NewTaskRepository()

	taskID, err := repo.Create(context.Background(), "Task 1")
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Verify getting task by ID
	task, err := repo.GetByID(context.Background(), taskID)
	if err != nil {
		t.Fatalf("Failed to get task by ID: %v", err)
	}

	if task.Title != "Task 1" {
		t.Errorf("Expected task title 'Task 1', got '%s'", task.Title)
	}

	// Verify getting non-existing task
	_, err = repo.GetByID(context.Background(), "non-existing-id")
	if err != errors.ErrTaskNotFound {
		t.Errorf("Expected task not found error, got: %v", err)
	}
}

func TestTaskRepository_GetAll(t *testing.T) {
	repo := NewTaskRepository()

	// Create some tasks
	_, err := repo.Create(context.Background(), "Task 1")
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	_, err = repo.Create(context.Background(), "Task 2")
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Verify getting all tasks
	tasks, err := repo.GetAll(context.Background())
	if err != nil {
		t.Fatalf("Failed to get all tasks: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
}
