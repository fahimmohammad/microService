package repository

import (
	"context"
	"time"

	pb "github.com/haquenafeem/basic-microservice/proto/todo"
	"github.com/haquenafeem/basic-microservice/todo/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type todoRepository struct {
	client *mongo.Client
	config config.Config
}

// Create ...
func (tr todoRepository) Create(todo *pb.Todo) (*pb.Todo, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := tr.client.Database(tr.config.Database).Collection(tr.config.TodoColl)
	_, err := collection.InsertOne(ctx, todo)
	if err != nil {
		return &pb.Todo{}, err
	}
	return todo, nil
}

// NewMongoRepository returns a new todo repo for mongo database
func NewMongoRepository(dbClient *mongo.Client) TodoRepository {
	return &todoRepository{
		client: dbClient,
		config: config.New(),
	}
}
