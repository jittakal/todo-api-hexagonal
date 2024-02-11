package grpc

import (
	"context"

	"github.com/jittakal/todo-api-hexagonal/api/proto"
	"github.com/jittakal/todo-api-hexagonal/internal/domain"
	"github.com/jittakal/todo-api-hexagonal/internal/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCTaskHandler struct {
	proto.UnimplementedTaskServiceServer
	repo domain.TaskRepository
}

func NewGRPCTaskHandler(repository domain.TaskRepository) *GRPCTaskHandler {
	return &GRPCTaskHandler{
		repo: repository,
	}
}

func (s *GRPCTaskHandler) GetTask(ctx context.Context, req *proto.GetTaskRequest) (*proto.Task, error) {
	task, err := s.repo.GetByID(ctx, req.GetId())
	if err != nil {
		if err == errors.ErrTaskNotFound {
			return nil, status.Errorf(codes.NotFound, "Task not found")
		} else {
			return nil, status.Errorf(codes.Internal, "Failed to get task")
		}
	}

	protoTask := &proto.Task{
		Id:    task.ID,
		Title: task.Title,
		Done:  task.Done,
	}
	return protoTask, nil
}
