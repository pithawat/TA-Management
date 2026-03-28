package service_test

import (
	"errors"
	"testing"
	"time"

	"TA-management/internal/modules/announce/dto/request"
	"TA-management/internal/modules/announce/dto/response"
	"TA-management/internal/modules/announce/repository/mocks"
	"TA-management/internal/modules/announce/service"

	"github.com/stretchr/testify/assert"
)

// GetEmailHistory ─────────────────────────────────────────────────────────────

func TestGetEmailHistory(t *testing.T) {
	t.Run("T001_Success", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.AnnouncementRepository)
		svc := service.NewAnnouncementService(mockRepo, nil)

		expected := &[]response.EmailHistory{
			{
				Id:           1,
				Subject:      "Test Subject",
				Body:         "Test Body",
				ReceivedName: "All TA",
				NReceived:    5,
				Status:       "success",
				CreatedDate:  time.Now(),
			},
		}
		mockRepo.On("GetEmailHistory").Return(expected, nil)

		// Act
		result, err := svc.GetEmailHistory()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("T002_RepositoryError", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.AnnouncementRepository)
		svc := service.NewAnnouncementService(mockRepo, nil)

		mockRepo.On("GetEmailHistory").Return((*[]response.EmailHistory)(nil), errors.New("db error"))

		// Act
		result, err := svc.GetEmailHistory()

		// Assert
		assert.Nil(t, result)
		assert.EqualError(t, err, "db error")
		mockRepo.AssertExpectations(t)
	})
}

// JoinDiscordChannel ──────────────────────────────────────────────────────────

func TestJoinDiscordChannel(t *testing.T) {
	t.Run("T003_GetDiscordRoleIDError_ShouldReturnError", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.AnnouncementRepository)
		// Pass nil for discordClient — this path should fail before calling it
		svc := service.NewAnnouncementService(mockRepo, nil)

		mockRepo.On("GetDiscordRoleID", 10).Return("", errors.New("no role found"))

		// Act
		inviteLink, err := svc.JoinDiscordChannel(10)

		// Assert
		assert.Empty(t, inviteLink)
		assert.EqualError(t, err, "no role found")
		mockRepo.AssertExpectations(t)
	})
}

// SendMailToAllCourse ─────────────────────────────────────────────────────────

func TestSendMailToAllCourse(t *testing.T) {
	t.Run("T004_RepositoryError_ShouldNotPanic", func(t *testing.T) {
		// Arrange - when repo returns error, function should return silently
		mockRepo := new(mocks.AnnouncementRepository)
		svc := service.NewAnnouncementService(mockRepo, nil)

		mockRepo.On("GetStudentEmailByCourseIDs").Return((*request.EmailRequest)(nil), errors.New("db error"))

		rq := request.MailForAllCourse{Subject: "Hello", Body: "Hi everyone"}

		// Act - should not panic
		assert.NotPanics(t, func() {
			svc.SendMailToAllCourse(rq)
		})
		mockRepo.AssertExpectations(t)
	})
}

// SendMailToCourse ────────────────────────────────────────────────────────────

func TestSendMailToCourse(t *testing.T) {
	t.Run("T005_RepositoryError_ShouldNotPanic", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.AnnouncementRepository)
		svc := service.NewAnnouncementService(mockRepo, nil)

		mockRepo.On("GetStudentEmailByCourseID", 5).
			Return((*request.EmailRequest)(nil), (*response.CourseDetail)(nil), errors.New("no course"))

		rq := request.MailForCourse{CourseID: 5, Subject: "Test", Body: "Body"}

		// Act - should not panic
		assert.NotPanics(t, func() {
			svc.SendMailToCourse(rq)
		})
		mockRepo.AssertExpectations(t)
	})
}

// SendMailToTA ────────────────────────────────────────────────────────────────

func TestSendMailToTA(t *testing.T) {
	t.Run("T006_RepositoryError_ShouldNotPanic", func(t *testing.T) {
		// Arrange
		mockRepo := new(mocks.AnnouncementRepository)
		svc := service.NewAnnouncementService(mockRepo, nil)

		mockRepo.On("GetStudentEmailByStudentID", 42).
			Return((*request.EmailRequest)(nil), errors.New("student email not found"))

		rq := request.MailForTA{StudentID: 42, Subject: "Test", Body: "Body"}

		// Act - should not panic
		assert.NotPanics(t, func() {
			svc.SendMailToTA(rq)
		})
		mockRepo.AssertExpectations(t)
	})
}
