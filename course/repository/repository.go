package repository

import (
	pb "github.com/fahimsGit/basic-microservice/proto/course"
)

type CourseRepository interface {
	CreateCourse(*pb.Course) (*pb.Course, error)
	GetSingleCourse(string) (*pb.Course, error)
}
