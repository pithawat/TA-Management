package service_test

import (
	"TA-management/internal/modules/course/dto/request"
	"TA-management/internal/modules/course/repository/mocks"
	"TA-management/internal/modules/course/service"
	"TA-management/internal/modules/shared/dto/testutils"
	"database/sql"

	// "TA-management/internal/modules/shared/dto/testutils"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	testDB = testutils.InitTestDB()

	exitCode := m.Run()
	testDB.Close()
	os.Exit(exitCode)
}

func TestApplyJobPost(t *testing.T) {
	t.Run("T001_AllocationFull_ShouldReturnErrorAndNotCallStartDBTx", func(t *testing.T) {
		// When TA allocation is already full, ApplyJobPost should return an error
		// and must NOT proceed to start a DB transaction.
		mockRepo := new(mocks.CourseRepository)
		svc := service.NewCourseService(mockRepo, nil)

		jobID := 101
		input := request.ApplyJobPost{JobPostID: &jobID, StudentID: 1}

		mockRepo.On("CheckStudentJobpostStatus", input).Return(true, nil)
		mockRepo.On("GetTaAllocation", jobID).Return(5, nil)
		mockRepo.On("CountTaAllocation", jobID).Return(5, nil)

		res, err := svc.ApplyJobPost(input)
		assert.Nil(t, res)
		assert.EqualError(t, err, "ta allocation is full")
		mockRepo.AssertNotCalled(t, "StartDBTx")
	})

	t.Run("T002_DuplicateApply", func(t *testing.T) {

		mockRepo := new(mocks.CourseRepository)
		svc := service.NewCourseService(mockRepo, nil)

		jobID := 1
		input := request.ApplyJobPost{JobPostID: &jobID}

		mockRepo.On("CheckStudentJobpostStatus", input).Return(false, nil)

		res, err := svc.ApplyJobPost(input)
		assert.Nil(t, res)
		assert.EqualError(t, err, "already apply to this jobpost")
		mockRepo.AssertNotCalled(t, "StartDBTx")
	})

	t.Run("T003_ApplyJobPostSuccess", func(t *testing.T) {

		mockRepo := new(mocks.CourseRepository)
		svc := service.NewCourseService(mockRepo, nil)

		jobID := 1
		transcript := []byte("hello")
		input := request.ApplyJobPost{JobPostID: &jobID, TranscriptBytes: &transcript}

		mockRepo.On("CheckStudentJobpostStatus", input).Return(true, nil)
		mockRepo.On("GetTaAllocation", jobID).Return(5, nil)
		mockRepo.On("CountTaAllocation", jobID).Return(2, nil)

		mockTx := &sql.Tx{}
		mockRepo.On("StartDBTx").Return(mockTx, nil)
		mockRepo.On("UpsertTranscript", mockTx, input).Return(nil)
		mockRepo.On("UpdateStudentData", mockTx, input).Return(nil)
		mockRepo.On("InsertApplication", mockTx, input).Return(555, nil)
		mockRepo.On("CommitTx", mockTx).Return(nil)
		mockRepo.On("RollbackTx", mockTx).Return(nil)

		res, err := svc.ApplyJobPost(input)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, 555, res.Id)

		mockRepo.AssertExpectations(t)
	})
}

func TestApproveApplication(t *testing.T) {

	t.Run("T004_TaAllocationFull", func(t *testing.T) {

		mockRepo := new(mocks.CourseRepository)
		svc := service.NewCourseService(mockRepo, nil)

		mockRepo.On("GetApproveApplicationData", 1).Return(1, 1, 1, nil)
		mockRepo.On("GetTaAllocation", 1).Return(2, nil)
		mockRepo.On("CountTaAllocation", 1).Return(2, nil)

		res, err := svc.ApproveApplication(1)
		assert.Nil(t, res)
		assert.EqualError(t, err, "ta allocation is full")
		mockRepo.AssertNotCalled(t, "StartDBTx")
	})

	t.Run("TC005_ApproveSuccess", func(t *testing.T) {
		mockRepo := new(mocks.CourseRepository)
		svc := service.NewCourseService(mockRepo, nil)

		mockRepo.On("GetApproveApplicationData", 1).Return(1, 1, 1, nil)
		mockRepo.On("GetTaAllocation", 1).Return(4, nil)
		mockRepo.On("CountTaAllocation", 1).Return(2, nil).Once()

		mockTx := &sql.Tx{}
		mockRepo.On("StartDBTx").Return(mockTx, nil)
		mockRepo.On("UpdateApplicationStatus", mockTx, 1).Return(nil)
		mockRepo.On("InsertTaCourse", mockTx, 1, 1).Return(nil)
		mockRepo.On("CommitTx", mockTx).Return(nil)
		mockRepo.On("RollbackTx", mockTx).Return(nil)
		mockRepo.On("CountTaAllocation", 1).Return(4, nil).Once()
		mockRepo.On("UpdateJobPostStatus", 1).Return(nil)

		res, err := svc.ApproveApplication(1)
		if err != nil {
			t.Logf("error: %v", err)
		}
		assert.NotNil(t, res)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateCourse(t *testing.T) {
	mockRepo := new(mocks.CourseRepository)
	svc := service.NewCourseService(mockRepo, nil)
	dummy := request.CreateCourse{
		CourseName:      "Advance Golang",
		CourseCode:      "CS101",
		ProfessorID:     1,
		CourseProgramID: 1,
		Sec:             "A",
		SemesterID:      1,
		ClassdayID:      1,
		ClassStart:      "09:00",
		ClassEnd:        "12:00",
		WorkHour:        2,
	}
	t.Run("T006CreateDupllicateCourse", func(t *testing.T) {

		mockRepo.On("IsCourseExist", dummy).Return(1, nil)

		res, err := svc.CreateCourse(dummy)
		assert.Nil(t, res)
		assert.EqualError(t, err, "already have this course")
		mockRepo.AssertNotCalled(t, "CreateCourse")
	})

	t.Run("T007CreateCourseSuccess", func(t *testing.T) {
		localDummy := request.CreateCourse{
			CourseCode: "CS101",
			CourseName: "Intro to Programming",
		}
		mockRepo.On("IsCourseExist", localDummy).Return(0, nil)
		mockRepo.On("CreateCourse", localDummy).Return(1, nil)

		res, err := svc.CreateCourse(localDummy)
		assert.NotNil(t, res)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
