package handler

import (
	"context"

	pb "github.com/fahimsGit/basic-microservice/proto/student"
	"github.com/fahimsGit/basic-microservice/student/repository"
	"github.com/google/uuid"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type service struct {
	repoistory repository.StudentRepository
}

func (s *service) CreateStudent(ctx context.Context, req *pb.RequestCreateStudent) (*pb.ResponseCreateStudent, error) {
	student := req.GetStudent()
	student.Id = uuid.New().String()
	student, err := s.repoistory.CreateStudent(student)

	if err != nil {
		return &pb.ResponseCreateStudent{
			Status: &pb.Status{
				Success: false,
				Error:   err.Error(),
			},
		}, nil
	}

	return &pb.ResponseCreateStudent{
		Name:   student.Name,
		Id:     student.Id,
		Status: &pb.Status{Success: true, Error: ""},
	}, nil
}

func (s *service) GetAllStudent(ctx context.Context, empty *emptypb.Empty) (*pb.ResponseGetAllStudent, error) {
	students, err := s.repoistory.GetAllStudent()
	if err != nil {
		return &pb.ResponseGetAllStudent{
			Students: students,
			Status:   &pb.Status{},
		}, err
	}
	return &pb.ResponseGetAllStudent{
		Students: students,
		Status:   &pb.Status{},
	}, nil
}

func NewService(repo repository.StudentRepository) pb.StudentServiceServer {
	return &service{
		repoistory: repo,
	}
}
