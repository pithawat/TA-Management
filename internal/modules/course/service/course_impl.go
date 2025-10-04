package service

import (
	"TA-management/internal/modules/course/dto/request"
	"TA-management/internal/modules/course/dto/response"
	"TA-management/internal/modules/course/repository"
	"fmt"
)

type CourseServiceImplementation struct {
	repo repository.CourseRepository
}

func NewCourseService(repo repository.CourseRepository) CourseServiceImplementation {
	return CourseServiceImplementation{repo: repo}
}

func (s CourseServiceImplementation) GetAllCourse() (*response.GeneralResponse, error) {

	courses, err := s.repo.GetAllCourse()
	if err != nil {
		return nil, err
	}
	response := response.GeneralResponse{
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
