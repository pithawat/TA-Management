package service

import (
	"TA-management/internal/modules/course/dto/request"
	courseResponse "TA-management/internal/modules/course/dto/response"
	"TA-management/internal/modules/course/repository"
	"TA-management/internal/modules/shared/dto/response"
	"fmt"
)

type CourseServiceImplementation struct {
	repo repository.CourseRepository
}

func NewCourseService(repo repository.CourseRepository) CourseServiceImplementation {
	return CourseServiceImplementation{repo: repo}
}

func (s CourseServiceImplementation) GetAllCourse() (*response.RequestDataResponse, error) {

	courses, err := s.repo.GetAllCourse()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	response := response.RequestDataResponse{
		Data:    courses,
		Message: "Success",
	}

	return &response, nil
}

func (s CourseServiceImplementation) CreateCourse(body request.CreateCourse) (response.CreateResponse, error) {
	id, err := s.repo.CreateCourse(body)
	if err != nil {
		fmt.Println(err)
		return response.CreateResponse{
			Message: "Create Failed!",
		}, err
	}

	return response.CreateResponse{
		Message: "Created successfully!",
		Id:      id,
	}, nil
}

func (s CourseServiceImplementation) UpdateCourse(body request.UpdateCourse) (response.GeneralResponse, error) {
	err := s.repo.UpdateCourse(body)
	if err != nil {
		return response.GeneralResponse{Message: "Update Failed!"}, err
	}
	return response.GeneralResponse{Message: "Update Successful"}, err
}

func (s CourseServiceImplementation) DeleteCourse(id int) (response.GeneralResponse, error) {
	err := s.repo.DeleteCourse(id)
	if err != nil {
		fmt.Println(err)
		return response.GeneralResponse{Message: "Delete Failed!"}, err
	}
	return response.GeneralResponse{Message: "Delete Successful"}, err
}

func (s CourseServiceImplementation) ApplyCourse(body request.ApplyCourse) (*response.CreateResponse, error) {
	id, err := s.repo.ApplyCourse(body)
	if err != nil {
		return nil, err
	}
	return &response.CreateResponse{
		Message: "Apply course successfully",
		Id:      id,
	}, nil
}

func (s CourseServiceImplementation) GetApplicationByStudentId(studentId int) (*response.RequestDataResponse, error) {
	applications, err := s.repo.GetApplicationByStudentId(studentId)
	if err != nil {
		return nil, err
	}
	return &response.RequestDataResponse{
		Data:    applications,
		Message: "GET success",
	}, nil

}

func (s CourseServiceImplementation) GetApplicationByCourseId(courseId int) (*response.RequestDataResponse, error) {
	applications, err := s.repo.GetApplicationByCourseId(courseId)
	if err != nil {
		return nil, err
	}
	return &response.RequestDataResponse{
		Data:    applications,
		Message: "GET success",
	}, nil
}

func (s CourseServiceImplementation) GetApplicationDetail(applicationId int) (*response.RequestDataResponse, error) {
	application, err := s.repo.GetApplicationDetail(applicationId)
	if err != nil {
		return nil, nil
	}
	return &response.RequestDataResponse{
		Data:    application,
		Message: "GET SUCCESS",
	}, nil

}

func (s CourseServiceImplementation) GetApplicationPdf(applicationId int) (*courseResponse.ApplicationTrancript, error) {
	applicationPdf, err := s.repo.GetApplicationPdf(applicationId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return applicationPdf, nil
}

func (s CourseServiceImplementation) ApproveApplication(applicationId int) (*response.GeneralResponse, error) {
	err := s.repo.ApproveApplication(applicationId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &response.GeneralResponse{
		Message: "Approved application Successful",
	}, nil
}

func (s CourseServiceImplementation) GetProfessorCourse(professorId int) (*response.RequestDataResponse, error) {
	courses, err := s.repo.GetProfessorCourse(professorId)
	if err != nil {
		return nil, err
	}
	response := response.RequestDataResponse{
		Data:    courses,
		Message: "Success",
	}

	return &response, nil
}
