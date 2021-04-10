package repository

import (
	"context"
	"log"
	"time"

	//courseproto "github.com/fahimsGit/basic-microservice/proto/course"
	pb "github.com/fahimsGit/basic-microservice/proto/student"
	"github.com/fahimsGit/basic-microservice/student/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type studentRepository struct {
	client *mongo.Client
	config config.DBConfig
}

// Create ...
func (repoService studentRepository) CreateStudent(student *pb.Student) (*pb.Student, error) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Minute)
	collection := repoService.client.Database(repoService.config.Dbname).Collection(repoService.config.Students)
	_, err := collection.InsertOne(ctx, student)
	if err != nil {
		return &pb.Student{}, err
	}
	return student, nil
}

func (repoService studentRepository) GetAllStudent() ([]*pb.Student, error) {
	var students []*pb.Student
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Minute)
	collection := repoService.client.Database(repoService.config.Dbname).Collection(repoService.config.Students)
	resultCursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		return students, err
	}
	for resultCursor.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem pb.Student
		err := resultCursor.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		students = append(students, &elem)
	}
	return students, nil
}
func (repoService studentRepository) CreateCourseEnrollment(enrolment *pb.Enrolment) (*pb.Enrolment, error) {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)
	collection := repoService.client.Database(repoService.config.Dbname).Collection(repoService.config.Enrolment)
	_, err := collection.InsertOne(ctx, enrolment)
	if err != nil {
		return enrolment, err
	}
	return enrolment, nil

}

func (repoService studentRepository) GetSingleStudent(studentId string) (*pb.Student, error) {
	var student *pb.Student
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Minute)
	collection := repoService.client.Database(repoService.config.Dbname).Collection(repoService.config.Students)
	err := collection.FindOne(ctx, bson.M{"id": studentId}).Decode(&student)
	if err != nil {
		return &pb.Student{
			Name:    "",
			Id:      studentId,
			Roll:    "",
			Session: "",
		}, err
	}
	return &pb.Student{
		Name:    student.Name,
		Id:      student.Id,
		Roll:    student.Roll,
		Session: student.Session,
	}, nil
}

func (repoService studentRepository) GetAllEnrollment(studentId string) ([]*pb.Enrolment, error) {
	var enrolment []*pb.Enrolment
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Minute)
	collection := repoService.client.Database(repoService.config.Dbname).Collection(repoService.config.Students)
	resultCursor, err := collection.Find(ctx, bson.M{"id": studentId})

	if err != nil {
		return enrolment, err
	}
	for resultCursor.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem pb.Enrolment
		err := resultCursor.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		enrolment = append(enrolment, &elem)
	}
	return enrolment, nil
}
func NewMongoRepository(dbClient *mongo.Client) StudentRepository {
	return &studentRepository{
		client: dbClient,
		config: config.New(),
	}
}
