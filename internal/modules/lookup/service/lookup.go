package service

import "TA-management/internal/modules/lookup/dto/response"

type LookupService interface {
	GetCourseProgram() (*[]response.LookupResponse, error)
	GetClassday() (*[]response.LookupResponse, error)
	GetSemester() (*[]response.LookupResponse, error)
	GetGrade() (*[]response.LookupResponse, error)
}
