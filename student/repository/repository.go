package repository

import (
	pb "github.com/fahimsGit/basic-microservice/proto/student"
)

type StudentRepository interface {
	CreateStudent(*pb.Student) (*pb.Student, error)
	GetAllStudent() ([]*pb.Student, error)
}
