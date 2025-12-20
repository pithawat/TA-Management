package service

import (
	"TA-management/internal/modules/course/dto/request"
	courseResponse "TA-management/internal/modules/course/dto/response"
	"TA-management/internal/modules/shared/dto/response"
)

type CourseService interface {
	GetAllCourse() (*response.RequestDataResponse, error)
	CreateCourse(body request.CreateCourse) (response.CreateResponse, error)
	UpdateCourse(body request.UpdateCourse) (response.GeneralResponse, error)
	DeleteCourse(courseId int) (response.GeneralResponse, error)
	ApplyCourse(body request.ApplyCourse) (*response.CreateResponse, error)
	GetApplicationByStudentId(studentId int) (*response.RequestDataResponse, error)
	GetApplicationByCourseId(CourseId int) (*response.RequestDataResponse, error)
	GetApplicationDetail(ApplicationId int) (*response.RequestDataResponse, error)
	GetApplicationPdf(ApplicationId int) (*courseResponse.ApplicationTrancript, error)
	ApproveApplication(ApplicationId int) (*response.GeneralResponse, error)
}
