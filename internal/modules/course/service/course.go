package service

import (
	"TA-management/internal/modules/course/dto/request"
	"TA-management/internal/modules/course/dto/response"
)

type CourseService interface {
	GetAllCourse() (*response.GeneralResponse, error)
	CreateCourse(body request.CreateCourse) error
}
