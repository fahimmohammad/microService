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
	courseId := req.GetCourseId()
	studentId := req.GetStudentId()
	grpcReq := &courseProto.RequestGetSingleCourse{
		CourseId: courseId,
	}
	course, err := s.client.GetSingleCourse(context.Background(), grpcReq)
	if err != nil {
		fmt.Println(err)
		return &pb.ResponseCreateCourseEnrollment{
			Id:         studentId,
			Name:       "",
			CourseId:   courseId,
			CourseName: "",
			Status: &pb.Status{
				Success: false,
				Error:   "Course Not found",
			},
		}, err
	}

	reponse, err := s.repoistory.CreateCourseEnrollment(req, course.Course.Name)

	if err != nil {
		return &pb.ResponseCreateCourseEnrollment{
			Id:         studentId,
			Name:       "",
			CourseId:   courseId,
			CourseName: "",
			Status:     &pb.Status{Success: false, Error: err.Error()},
		}, nil
	}

	return &pb.ResponseCreateCourseEnrollment{
		Id:         reponse.Id,
		Name:       "StudentName",
		CourseId:   reponse.CourseId,
		CourseName: reponse.CourseName,
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
