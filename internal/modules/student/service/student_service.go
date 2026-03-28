package service

import (
	"fmt"
	"io"

	"TA-management/internal/modules/student/dto/request"
	"TA-management/internal/modules/student/dto/response"
	"TA-management/internal/modules/student/repository"
)

type StudentService interface {
	GetStudentProfile(studentID int) (*response.StudentProfile, error)
	UpdateStudentProfile(studentID int, req *request.UpdateProfile) error
	UploadDocument(studentID int, docType string, file io.Reader, fileName string) error
	GetDocument(studentID int, docType string) ([]byte, string, error)
	DeleteDocument(studentID int, docType string) error
}

type studentServiceImpl struct {
	repo repository.StudentRepository
}

func NewStudentService(repo repository.StudentRepository) StudentService {
	return &studentServiceImpl{repo: repo}
}

func (s *studentServiceImpl) GetStudentProfile(studentID int) (*response.StudentProfile, error) {
	return s.repo.GetStudentByID(studentID)
}

func (s *studentServiceImpl) UpdateStudentProfile(studentID int, req *request.UpdateProfile) error {
	// Validate input
	if req.FirstnameThai == "" || req.LastnameThai == "" {
		return fmt.Errorf("Thai name is required")
	}
	if req.PhoneNumber == "" {
		return fmt.Errorf("phone number is required")
	}

	return s.repo.UpdateStudent(studentID, req.FirstnameThai, req.LastnameThai, req.PhoneNumber)
}

func (s *studentServiceImpl) UploadDocument(studentID int, docType string, file io.Reader, fileName string) error {
	// Read file bytes
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Validate file size (max 10MB)
	if len(fileBytes) > 10*1024*1024 {
		return fmt.Errorf("file size exceeds 10MB limit")
	}

	return s.repo.UploadDocument(studentID, docType, fileBytes, fileName)
}

func (s *studentServiceImpl) GetDocument(studentID int, docType string) ([]byte, string, error) {
	return s.repo.GetDocument(studentID, docType)
}

func (s *studentServiceImpl) DeleteDocument(studentID int, docType string) error {
	return s.repo.DeleteDocument(studentID, docType)
}
