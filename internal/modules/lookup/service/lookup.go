package service

import (
	"TA-management/internal/modules/lookup/dto/request"
	"TA-management/internal/modules/lookup/dto/response"
	"TA-management/internal/modules/ta_duty/entity"
)

type LookupService interface {
	GetCourseProgram() (*[]response.LookupResponse, error)
	GetClassday() (*[]response.LookupResponse, error)
	GetSemester() (*[]response.SemesterResponse, error)
	GetSemesterDropdown() (*[]response.LookupResponse, error)
	GetGrade() (*[]response.LookupResponse, error)
	GetProfessors() (*[]response.LookupResponse, error)
	SyncOfficialHoliday(apiKey string, url string) error
	GetHolidays(month int, year int) ([]entity.Holiday, error)
	AddSpecialHoliday(req request.CreateHoliday) error
	DeleteHoliday(id int) error
	GetTA(searchVal string) (*[]response.TaDetail, error)
	GetAvailableMonths(month int) (*[]response.AvailableMonth, error)
	GetTranscript(studentID int) (*response.PdfFile, error)
	GetBankAccount(studentID int) (*response.PdfFile, error)
	GetStudentCard(studentID int) (*response.PdfFile, error)
	AddSemester(rq request.CreateSemester) error
	UpdateSemester(rq request.UpdateSemester) (*[]response.SemesterResponse, error)
	SetSemesterActive(semesterID int) error
}
