package service_test

import (
	"bytes"
	"errors"
	"testing"

	"TA-management/internal/modules/student/dto/request"
	"TA-management/internal/modules/student/dto/response"
	"TA-management/internal/modules/student/repository/mocks"
	"TA-management/internal/modules/student/service"

	"github.com/stretchr/testify/assert"
)

// ─── GetStudentProfile ────────────────────────────────────────────────────────

func TestGetStudentProfile(t *testing.T) {
	t.Run("T001_Success", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.StudentRepository)
		svc := service.NewStudentService(mockRepo)

		expected := &response.StudentProfile{
			StudentID:     1,
			FirstnameThai: "สมชาย",
			LastnameThai:  "ใจดี",
			Email:         "test@example.com",
			PhoneNumber:   "0812345678",
		}
		mockRepo.On("GetStudentByID", 1).Return(expected, nil)

		// Act
		result, err := svc.GetStudentProfile(1)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("T002_StudentNotFound", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.StudentRepository)
		svc := service.NewStudentService(mockRepo)

		mockRepo.On("GetStudentByID", 99).Return((*response.StudentProfile)(nil), errors.New("student not found"))

		// Act
		result, err := svc.GetStudentProfile(99)

		// Assert
		assert.Nil(t, result)
		assert.EqualError(t, err, "student not found")
		mockRepo.AssertExpectations(t)
	})
}

// ─── UpdateStudentProfile ─────────────────────────────────────────────────────

func TestUpdateStudentProfile(t *testing.T) {
	t.Run("T003_MissingThaiName_ShouldReturnError", func(t *testing.T) {
		// Arrange - the service validates before calling repo
		mockRepo := new(mocks.StudentRepository)
		svc := service.NewStudentService(mockRepo)

		req := &request.UpdateProfile{
			FirstnameThai: "",
			LastnameThai:  "ใจดี",
			PhoneNumber:   "0812345678",
		}

		// Act
		err := svc.UpdateStudentProfile(1, req)

		// Assert - should block on validation, repo never called
		assert.EqualError(t, err, "Thai name is required")
		mockRepo.AssertNotCalled(t, "UpdateStudent")
	})

	t.Run("T004_MissingPhoneNumber_ShouldReturnError", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.StudentRepository)
		svc := service.NewStudentService(mockRepo)

		req := &request.UpdateProfile{
			FirstnameThai: "สมชาย",
			LastnameThai:  "ใจดี",
			PhoneNumber:   "",
		}

		// Act
		err := svc.UpdateStudentProfile(1, req)

		// Assert
		assert.EqualError(t, err, "phone number is required")
		mockRepo.AssertNotCalled(t, "UpdateStudent")
	})

	t.Run("T005_Success", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.StudentRepository)
		svc := service.NewStudentService(mockRepo)

		req := &request.UpdateProfile{
			FirstnameThai: "สมชาย",
			LastnameThai:  "ใจดี",
			PhoneNumber:   "0812345678",
		}
		mockRepo.On("UpdateStudent", 1, "สมชาย", "ใจดี", "0812345678").Return(nil)

		// Act
		err := svc.UpdateStudentProfile(1, req)

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// ─── UploadDocument ───────────────────────────────────────────────────────────

func TestUploadDocument(t *testing.T) {
	t.Run("T006_FileSizeExceedsLimit", func(t *testing.T) {
		// Arrange - create a reader with >10MB data
		mockRepo := new(mocks.StudentRepository)
		svc := service.NewStudentService(mockRepo)

		oversizedData := make([]byte, 10*1024*1024+1) // 10MB + 1 byte
		reader := bytes.NewReader(oversizedData)

		// Act
		err := svc.UploadDocument(1, "transcript", reader, "transcript.pdf")

		// Assert
		assert.EqualError(t, err, "file size exceeds 10MB limit")
		mockRepo.AssertNotCalled(t, "UploadDocument")
	})

	t.Run("T007_Success", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.StudentRepository)
		svc := service.NewStudentService(mockRepo)

		data := []byte("fake-pdf-content")
		reader := bytes.NewReader(data)

		mockRepo.On("UploadDocument", 1, "transcript", data, "transcript.pdf").Return(nil)

		// Act
		err := svc.UploadDocument(1, "transcript", reader, "transcript.pdf")

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

// ─── GetDocument ──────────────────────────────────────────────────────────────

func TestGetDocument(t *testing.T) {
	t.Run("T008_Success", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.StudentRepository)
		svc := service.NewStudentService(mockRepo)

		expectedBytes := []byte("file-bytes")
		expectedName := "transcript.pdf"
		mockRepo.On("GetDocument", 1, "transcript").Return(expectedBytes, expectedName, nil)

		// Act
		fileBytes, fileName, err := svc.GetDocument(1, "transcript")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedBytes, fileBytes)
		assert.Equal(t, expectedName, fileName)
		mockRepo.AssertExpectations(t)
	})

	t.Run("T009_NotFound", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.StudentRepository)
		svc := service.NewStudentService(mockRepo)

		mockRepo.On("GetDocument", 1, "bank_account").Return([]byte(nil), "", errors.New("document not found"))

		// Act
		_, _, err := svc.GetDocument(1, "bank_account")

		// Assert
		assert.EqualError(t, err, "document not found")
		mockRepo.AssertExpectations(t)
	})
}

// ─── DeleteDocument ───────────────────────────────────────────────────────────

func TestDeleteDocument(t *testing.T) {
	t.Run("T010_Success", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.StudentRepository)
		svc := service.NewStudentService(mockRepo)

		mockRepo.On("DeleteDocument", 1, "student_card").Return(nil)

		// Act
		err := svc.DeleteDocument(1, "student_card")

		// Assert
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("T011_NotFound", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.StudentRepository)
		svc := service.NewStudentService(mockRepo)

		mockRepo.On("DeleteDocument", 99, "transcript").Return(errors.New("document not found"))

		// Act
		err := svc.DeleteDocument(99, "transcript")

		// Assert
		assert.EqualError(t, err, "document not found")
		mockRepo.AssertExpectations(t)
	})
}
