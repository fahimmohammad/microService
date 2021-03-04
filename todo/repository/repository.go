package repository

import (
	pb "github.com/haquenafeem/basic-microservice/proto/todo"
)

// TodoRepository ...
type TodoRepository interface {
	Create(*pb.Todo) (*pb.Todo, error)
}
