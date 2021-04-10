package handler

import (
	"context"
	"fmt"

	courseProto "github.com/fahimsGit/basic-microservice/proto/course"
	pb "github.com/fahimsGit/basic-microservice/proto/student"
	"github.com/fahimsGit/basic-microservice/student/repository"
	"github.com/google/uuid"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type service struct {
	repoistory repository.StudentRepository
	client     courseProto.CourseServiceClient
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

func (s *service) CreateCourseEnrollment(ctx context.Context, req *pb.RequestCreateCourseEnrollment) (*pb.ResponseCreateCourseEnrollment, error) {
	var enrolment *pb.Enrolment

	grpcReq := &courseProto.RequestGetSingleCourse{
		CourseId: req.GetCourseId(),
	}
	course, err := s.client.GetSingleCourse(context.Background(), grpcReq)
	if err != nil {
		fmt.Println(err)
		return &pb.ResponseCreateCourseEnrollment{
			Enrolment: &pb.Enrolment{},
			Status: &pb.Status{
				Success: false,
				Error:   "No course found",
			},
		}, err
	}
	std, err := s.repoistory.GetSingleStudent(req.StudentId)
	if err != nil {
		fmt.Println(err)
		return &pb.ResponseCreateCourseEnrollment{
			Enrolment: &pb.Enrolment{},
			Status: &pb.Status{
				Success: false,
				Error:   "No student found",
			},
		}, err
	}

	enrolment = &pb.Enrolment{
		Id:         std.Id,
		Name:       std.Name,
		CourseId:   course.Course.Id,
		CourseName: course.Course.Name,
	}
	enrolment.CourseName = course.Course.Name

	reponse, err := s.repoistory.CreateCourseEnrollment(enrolment)

	if err != nil {
		return &pb.ResponseCreateCourseEnrollment{
			Enrolment: reponse,
			Status:    &pb.Status{Success: false, Error: err.Error()},
		}, nil
	}

	return &pb.ResponseCreateCourseEnrollment{
		Enrolment: reponse,
		Status:    &pb.Status{},
	}, nil
}
func (s *service) GetAllEnrollment(ctx context.Context, req *pb.RequestGetAllEnrollment) (*pb.ResponseGetAllEnrollment, error) {
	studentId := req.GetId()
	response, err := s.repoistory.GetAllEnrollment(studentId)
	if err != nil {
		return &pb.ResponseGetAllEnrollment{
			Enrolment: response,
			Status: &pb.Status{
				Success: false,
				Error:   "",
			},
		}, err
	}
	return &pb.ResponseGetAllEnrollment{
		Enrolment: response,
		Status: &pb.Status{
			Success: true,
			Error:   "",
		},
	}, nil
}
func NewService(repo repository.StudentRepository, client courseProto.CourseServiceClient) pb.StudentServiceServer {
	return &service{
		repoistory: repo,
		client:     client,
	}
}
