package repository

import (
	//courseproto "github.com/fahimsGit/basic-microservice/proto/course"
	pb "github.com/fahimsGit/basic-microservice/proto/student"
)

type StudentRepository interface {
	CreateStudent(*pb.Student) (*pb.Student, error)
	GetAllStudent() ([]*pb.Student, error)
	CreateCourseEnrollment(*pb.RequestCreateCourseEnrollment, string) (*pb.ResponseCreateCourseEnrollment, error)
}
