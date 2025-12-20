package service

import (
	"TA-management/internal/modules/lookup/dto/response"
	"TA-management/internal/modules/lookup/repository"
)

type LookupServiceImplementation struct {
	repo repository.LookupRepository
}

func NewLookupService(repo repository.LookupRepository) LookupServiceImplementation {
	return LookupServiceImplementation{repo: repo}
}

func (s LookupServiceImplementation) GetCourseProgram() (*[]response.LookupResponse, error) {
	result, err := s.repo.GetCourseProgram()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s LookupServiceImplementation) GetClassday() (*[]response.LookupResponse, error) {
	result, err := s.repo.GetClassday()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s LookupServiceImplementation) GetSemester() (*[]response.LookupResponse, error) {
	result, err := s.repo.GetSemester()
	if err != nil {
		return nil, err
	}
	return result, nil
}
