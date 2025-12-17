package service

import (
	"TA-management/internal/modules/course/dto/request"
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
		return response.GeneralResponse{Message: "Delete Failed!"}, err
	}
	return response.GeneralResponse{Message: "Delete Successful"}, err
}
