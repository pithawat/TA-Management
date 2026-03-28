package repository

import (
	"TA-management/internal/modules/course/dto/request"
	"TA-management/internal/modules/course/dto/response"
	sharedResponse "TA-management/internal/modules/shared/dto/response"
	"database/sql"
)

type CourseRepository interface {
	GetAllJobPost() ([]response.JobPost, error)
	GetAllJobPostAllStatus() ([]response.JobPost, error)
	GetAllJobPostByStudentId(studentId int) ([]response.JobPost, error)
	GetAllCourse() ([]response.Course, error)
	GetProfessorCourse(professorId int) ([]response.Course, error)
	CreateCourse(body request.CreateCourse) (int, error)
	IsCourseExist(body request.CreateCourse) (int, error)
	UpdateCourse(body request.UpdateCourse) error
	DeleteCourse(id int) error
	CreateJobPost(body request.CreateJobPost) (int, error)
	UpdateJobPost(body request.UpdateJobPost) error
	GetJobPostByID(jobPostId int) (*response.JobPost, error)
	DeleteJobPost(jobPostId int) error
	InsertApplication(tx *sql.Tx, body request.ApplyJobPost) (int, error)
	CheckStudentJobpostStatus(body request.ApplyJobPost) (bool, error)
	GetTaAllocation(jobPostID int) (int, error)
	CountTaAllocation(jobPostID int) (int, error)
	UpsertTranscript(tx *sql.Tx, body request.ApplyJobPost) error
	UpsertBankAccount(tx *sql.Tx, body request.ApplyJobPost) error
	UpsertStudentCard(tx *sql.Tx, body request.ApplyJobPost) error
	UpdateStudentData(tx *sql.Tx, body request.ApplyJobPost) error
	GetApplicationByStudentId(studentId int) ([]response.Application, error)
	GetAllTimeApprovedCoursesByStudentId(studentId int) ([]response.Application, error)
	GetApplicationByCourseId(courseId int) ([]response.Application, error)
	GetApplicationDetail(ApplicationId int) (*response.Application, error)
	GetApplicationTranscriptPdf(ApplicationId int) (*response.PdfFile, error)
	GetApplicationBankAccountPdf(ApplicationId int) (*response.PdfFile, error)
	GetApplicationStudentCardPdf(ApplicationId int) (*response.PdfFile, error)
	GetApplicationByProfessorId(professorId int) ([]response.Application, error)
	UpdateApplicationStatus(tx *sql.Tx, ApplicationId int) error
	InsertTaCourse(tx *sql.Tx, studentId int, courseId int) error
	UpdateJobPostStatus(jobPostId int) error
	GetApproveApplicationData(applicationData int) (int, int, int, error)
	RejectApplication(rq request.RejectApplication) error
	UpdateCourseDiscord(courseId int, roleId string, channelId string, channelName string) error
	GetDiscordRoleByCourseId(courseId int) (string, error)
	SoftDeleteExpiredData() error
	GetExpiredSemesters() (*sharedResponse.RequestDataResponse, error)
	GetCoursesBySemesterID(semesterID int) (*sharedResponse.RequestDataResponse, error)
	StartDBTx() (*sql.Tx, error)
	CommitTx(tx *sql.Tx) error
	RollbackTx(tx *sql.Tx) error
}
