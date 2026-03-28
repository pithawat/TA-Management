package service_test

import (
	"errors"
	"testing"

	"TA-management/internal/modules/ta_duty/dto/request"
	"TA-management/internal/modules/ta_duty/dto/response"
	"TA-management/internal/modules/ta_duty/repository/mocks"
	"TA-management/internal/modules/ta_duty/service"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func newTaDutyService(mockRepo *mocks.TaDutyRepository) service.TaDutyService {
	logger, _ := zap.NewDevelopment()
	return service.NewTaDutyServiceImplementation(mockRepo, logger.Sugar())
}

// ─── GetTADutyRoadmap ─────────────────────────────────────────────────────────

func TestGetTADutyRoadmap(t *testing.T) {
	t.Run("T001_Success", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.TaDutyRepository)
		svc := newTaDutyService(mockRepo)

		expected := &[]response.DutyChecklistItem{
			{Date: "2025-03-01", Status: "done", IsChecked: true},
			{Date: "2025-03-08", Status: "pending", IsChecked: false},
		}
		mockRepo.On("GetTADutyRoadmap", 10, 5).Return(expected, nil)

		// Act
		result, err := svc.GetTADutyRoadmap(10, 5)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("T002_RepositoryError", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.TaDutyRepository)
		svc := newTaDutyService(mockRepo)

		mockRepo.On("GetTADutyRoadmap", 10, 5).Return((*[]response.DutyChecklistItem)(nil), errors.New("db error"))

		// Act
		result, err := svc.GetTADutyRoadmap(10, 5)

		// Assert
		assert.Nil(t, result)
		assert.EqualError(t, err, "db error")
		mockRepo.AssertExpectations(t)
	})
}

// ─── MarkDutyAsDone ───────────────────────────────────────────────────────────

func TestMarkDutyAsDone(t *testing.T) {
	t.Run("T003_InvalidDateFormat_ShouldReturnError", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.TaDutyRepository)
		svc := newTaDutyService(mockRepo)

		// Act - pass a date string that is not YYYY-MM-DD
		result, err := svc.MarkDutyAsDone(10, 5, "01/03/2025")

		// Assert
		assert.Nil(t, result)
		assert.EqualError(t, err, "invalid date format")
		mockRepo.AssertNotCalled(t, "MarkDutyAsDone")
	})

	t.Run("T004_FutureDate_ShouldReturnError", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.TaDutyRepository)
		svc := newTaDutyService(mockRepo)

		// Act - use a clearly far future date
		result, err := svc.MarkDutyAsDone(10, 5, "2099-12-31")

		// Assert
		assert.Nil(t, result)
		assert.EqualError(t, err, "cannot check off future duties")
		mockRepo.AssertNotCalled(t, "MarkDutyAsDone")
	})

	t.Run("T005_ValidPastDate_Success", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.TaDutyRepository)
		svc := newTaDutyService(mockRepo)

		mockRepo.On("MarkDutyAsDone", 10, 5, "2020-01-01").Return(nil)

		// Act - a date well in the past
		result, err := svc.MarkDutyAsDone(10, 5, "2020-01-01")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "MarkDuty as Done successfully", result.Message)
		mockRepo.AssertExpectations(t)
	})

	t.Run("T006_RepositoryError_OnSave", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.TaDutyRepository)
		svc := newTaDutyService(mockRepo)

		mockRepo.On("MarkDutyAsDone", 10, 5, "2020-01-01").Return(errors.New("db error on save"))

		// Act
		result, err := svc.MarkDutyAsDone(10, 5, "2020-01-01")

		// Assert
		assert.Nil(t, result)
		assert.EqualError(t, err, "db error on save")
		mockRepo.AssertExpectations(t)
	})
}

// ─── ExportPaymentReport ──────────────────────────────────────────────────────

func TestExportPaymentReport(t *testing.T) {
	t.Run("T007_RepositoryError_ShouldReturnError", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.TaDutyRepository)
		svc := newTaDutyService(mockRepo)

		rq := request.ExportPaymentReportRequest{CourseID: 1, Month: 3, Year: 2025}

		mockRepo.On("GetTADutyDataExportPayment", 1, 3).
			Return((*[]request.CreatePaymentData)(nil), (*request.CourseDutyData)(nil), errors.New("query error"))

		// Act
		buf, courseData, err := svc.ExportPaymentReport(rq)

		// Assert
		assert.Nil(t, buf)
		assert.Nil(t, courseData)
		assert.EqualError(t, err, "query error")
		mockRepo.AssertExpectations(t)
	})
}

// ─── ExportSignatureSheet ─────────────────────────────────────────────────────

func TestExportSignatureSheet(t *testing.T) {
	t.Run("T008_RepositoryError_ShouldReturnError", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.TaDutyRepository)
		svc := newTaDutyService(mockRepo)

		rq := request.ExportSignatureSheet{CourseID: 1, Month: 3, Year: 2025}

		mockRepo.On("GetTADutyDataExportSignature", 1, 3).
			Return((*request.CreateSignatureSheet)(nil), (*request.CourseDutyData)(nil), errors.New("query error"))

		// Act
		buf, courseData, err := svc.ExportSignatureSheet(rq)

		// Assert
		assert.Nil(t, buf)
		assert.Nil(t, courseData)
		assert.EqualError(t, err, "query error")
		mockRepo.AssertExpectations(t)
	})
}
