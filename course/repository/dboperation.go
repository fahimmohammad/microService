package repository

import (
	"context"
	"time"

	"github.com/fahimsGit/basic-microservice/course/config"
	pb "github.com/fahimsGit/basic-microservice/proto/course"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type courseRepository struct {
	client *mongo.Client
	config config.DBConfig
}

// Create ...
func (repoService courseRepository) CreateCourse(course *pb.Course) (*pb.Course, error) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Minute)
	collection := repoService.client.Database(repoService.config.Dbname).Collection(repoService.config.Course)
	_, err := collection.InsertOne(ctx, course)
	if err != nil {
		return &pb.Course{}, err
	}
	return course, nil
}

func (repoService courseRepository) GetSingleCourse(courseId string) (*pb.Course, error) {
	var course *pb.Course
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Minute)
	collection := repoService.client.Database(repoService.config.Dbname).Collection(repoService.config.Course)
	err := collection.FindOne(ctx, bson.M{"id": courseId}).Decode(&course)
	if err != nil {
		return course, err
	}
	return course, nil
}

// NewMongoRepository returns a new todo repo for mongo database
func NewMongoRepository(dbClient *mongo.Client) CourseRepository {
	return &courseRepository{
		client: dbClient,
		config: config.New(),
	}
}
