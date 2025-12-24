package service

import (
	"TA-management/internal/modules/course/dto/request"
	courseResponse "TA-management/internal/modules/course/dto/response"
	"TA-management/internal/modules/shared/dto/response"
)

type CourseService interface {
	GetAllCourse() (*response.RequestDataResponse, error)
	GetAllCourseByStudentId(studentId int) (*response.RequestDataResponse, error)
	CreateCourse(body request.CreateCourse) (response.CreateResponse, error)
	UpdateCourse(body request.UpdateCourse) (response.GeneralResponse, error)
	DeleteCourse(courseId int) (response.GeneralResponse, error)
	CreateJobPost(body request.CreateJobPost) (response.CreateResponse, error)
	UpdateJobPost(body request.UpdateJobPost) (response.GeneralResponse, error)
	DeleteJobPost(jobPostId int) (response.GeneralResponse, error)
	ApplyJobPost(body request.ApplyJobPost) (*response.CreateResponse, error)
	GetApplicationByStudentId(studentId int) (*response.RequestDataResponse, error)
	GetApplicationByCourseId(CourseId int) (*response.RequestDataResponse, error)
	GetApplicationDetail(ApplicationId int) (*response.RequestDataResponse, error)
	GetApplicationTranscriptPdf(ApplicationId int) (*courseResponse.PdfFile, error)
	GetApplicationBankAccountPdf(ApplicationId int) (*courseResponse.PdfFile, error)
	GetApplicationStudentCardPdf(ApplicationId int) (*courseResponse.PdfFile, error)
	ApproveApplication(ApplicationId int) (*response.GeneralResponse, error)
	GetProfessorCourse(ProfessorId int) (*response.RequestDataResponse, error)
	GetApplicationByProfessorId(professorId int) (*response.RequestDataResponse, error)
}
