package handler

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	pb "github.com/haquenafeem/basic-microservice/proto/todo"
	"github.com/haquenafeem/basic-microservice/todo/repository"
)

type service struct {
	repo repository.TodoRepository
}

// CreateTodo ...
func (s *service) CreateTodo(ctx context.Context, r *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {
	todo := r.GetTodo()
	todo.ID = uuid.New().String()
	todo.CreatedAt = ptypes.TimestampNow()

	todo, err := s.repo.Create(todo)

	if err != nil {
		return &pb.CreateTodoResponse{
			Status: &pb.Status{
				Success: false,
				Error:   err.Error(),
			},
		}, nil
	}

	return &pb.CreateTodoResponse{
		Todo: todo,
		Status: &pb.Status{
			Success: true,
			Error:   "",
		},
	}, nil
}

// NewService returns a new todo server
func NewService(repo repository.TodoRepository) pb.TodoServiceServer {
	return &service{
		repo: repo,
	}
}
