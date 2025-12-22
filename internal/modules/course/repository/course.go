package repository

import (
	"TA-management/internal/modules/course/dto/request"
	"TA-management/internal/modules/course/dto/response"
)

type CourseRepository interface {
	GetAllCourse() ([]response.Course, error)
	GetAllCourseByStudentId(studentId int) ([]response.Course, error)
	GetProfessorCourse(professorId int) ([]response.Course, error)
	CreateCourse(body request.CreateCourse) (int, error)
	UpdateCourse(body request.UpdateCourse) error
	DeleteCourse(id int) error
	ApplyJobPost(body request.ApplyJobPost) (int, error)
	GetApplicationByStudentId(studentId int) ([]response.Application, error)
	GetApplicationByCourseId(courseId int) ([]response.Application, error)
	GetApplicationDetail(ApplicationId int) (*response.Application, error)
	GetApplicationTranscriptPdf(ApplicationId int) (*response.PdfFile, error)
	GetApplicationBankAccountPdf(ApplicationId int) (*response.PdfFile, error)
	GetApplicationStudentCardPdf(ApplicationId int) (*response.PdfFile, error)
	GetApplicationByProfessorId(professorId int) ([]response.Application, error)
	ApproveApplication(ApplicationId int) error
}
