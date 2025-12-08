package service

import (
	"TA-management/internal/modules/course/dto/request"
	"TA-management/internal/modules/shared/dto/response"
)

type CourseService interface {
	GetAllCourse() (*response.RequestDataResponse, error)
	CreateCourse(body request.CreateCourse) error
}
