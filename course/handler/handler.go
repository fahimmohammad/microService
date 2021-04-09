package handler

import (
	"context"

	"github.com/fahimsGit/basic-microservice/course/repository"
	pb "github.com/fahimsGit/basic-microservice/proto/course"
	"github.com/google/uuid"
)

type service struct {
	repoistory repository.CourseRepository
}

func (s *service) CreateCourse(ctx context.Context, req *pb.RequestCreateCourse) (*pb.ResponseCreateCourse, error) {
	course := req.GetCourse()
	course.Id = uuid.New().String()
	course, err := s.repoistory.CreateCourse(course)

	if err != nil {
		return &pb.ResponseCreateCourse{
			Status: &pb.Status{
				Success: false,
				Error:   err.Error(),
			},
		}, nil
	}

	return &pb.ResponseCreateCourse{
		Course: course,
		Status: &pb.Status{Success: true, Error: ""},
	}, nil
}

func (s *service) GetSingleCourse(ctx context.Context, req *pb.RequestGetSingleCourse) (*pb.ResponseGetSingleCourse, error) {

	id := req.GetCourseId()
	course, err := s.repoistory.GetSingleCourse(id)
	if err != nil {
		return &pb.ResponseGetSingleCourse{
			Course: course,
			Status: &pb.Status{},
		}, err
	}
	return &pb.ResponseGetSingleCourse{
		Course: course,
		Status: &pb.Status{},
	}, nil
}

func NewService(repo repository.CourseRepository) pb.CourseServiceServer {
	return &service{
		repoistory: repo,
	}
}
