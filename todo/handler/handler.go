package handler

import (
	"context"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	pb "github.com/haquenafeem/basic-microservice/proto/todo"
)

type service struct {
}

func (s *service) CreateTodo(ctx context.Context, r *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {

	todo := r.GetTodo()
	todo.ID = uuid.New().String()
	todo.CreatedAt = ptypes.TimestampNow()
	return &pb.CreateTodoResponse{
		Todo: r.GetTodo(),
		Status: &pb.Status{
			Success: true,
			Error:   "",
		},
	}, nil
}

// NewService returns a new todo server
func NewService() pb.TodoServiceServer {
	return &service{}
}
