package repository

import (
	"TA-management/internal/modules/lookup/dto/request"
	"TA-management/internal/modules/lookup/dto/response"
	"TA-management/internal/modules/ta_duty/entity"
)

type LookupRepository interface {
	GetCourseProgram() (*[]response.LookupResponse, error)
	GetClassday() (*[]response.LookupResponse, error)
	GetSemester() (*[]response.SemesterResponse, error)
	GetSemesterDropdown() (*[]response.LookupResponse, error)
	GetGrade() (*[]response.LookupResponse, error)
	GetProfessors() (*[]response.LookupResponse, error)
	SyncOfficialHoliday(holidays []request.CreateHoliday) error
	GetHolidaysByMonth(month int, year int) ([]entity.Holiday, error)
	DeleteHoliday(id int) error
	AddSpecialHoliday(holiday request.CreateHoliday) error
	GetTA(searchVal string) (*[]response.TaDetail, error)
	GetAvailableMonths(courseId int) (*[]response.AvailableMonth, error)
	GetStudentCard(studentID int) (*response.PdfFile, error)
	GetTranscript(studentID int) (*response.PdfFile, error)
	GetBankAccount(studentID int) (*response.PdfFile, error)
	AddSemester(rq request.CreateSemester) error
	UpdateSemester(rq request.UpdateSemester) error
	SetSemesterActive(semesterID int) error
}
