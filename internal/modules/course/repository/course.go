package repository

import (
	"TA-management/internal/modules/course/dto/request"
	"TA-management/internal/modules/course/dto/response"
)

type CourseRepository interface {
	GetAllCourse() ([]response.Course, error)
	CreateCourse(body request.CreateCourse) error
}
