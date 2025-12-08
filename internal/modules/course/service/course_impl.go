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

func (s CourseServiceImplementation) CreateCourse(body request.CreateCourse) error {
	err := s.repo.CreateCourse(body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
