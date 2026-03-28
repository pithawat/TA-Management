// File: TA-management/internal/modules/course/service/course_service_integration_test.go

package service_test

import (
	"TA-management/internal/constants"
	"TA-management/internal/modules/course/dto/request"
	"TA-management/internal/modules/course/repository"
	"TA-management/internal/modules/course/service" // Assuming this is the service location
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	// Assuming testDB and cleanDB are accessible from the repository_test setup,
	// OR you've moved TestMain to a root package.
	// **For this example, we assume testDB is globally accessible or passed in setup.**
)

// NOTE: In a real project, TestMain and testDB setup should live in a shared
// 'test' package at the root of your internal modules to avoid duplication.
// For now, we assume 'testDB' is initialized via the TestMain from the repository_test
// package or via a shared setup function.
var testDB *sql.DB // This variable will be initialized by the shared TestMain/Setup

// --- [ Replicate TestMain and cleanDB here if you cannot share cross-package ] ---
// For a clean running example, you MUST copy your TestMain and cleanDB functions here
// or ensure they run before this file. Let's assume you've placed them in a shared utility.
// ---

// setupServiceTest initializes the real repository and the real service.
func setupServiceTest(t *testing.T, db *sql.DB) *service.CourseServiceImplementation {
	// 1. Initialize the Real Repository with the Real Test DB connection
	repo := repository.NewCourseRepository(db)

	// 2. Initialize the Service with the Real Repository
	// Assuming a constructor like: NewCourseService(repo repository.CourseRepository)
	svc := service.NewCourseService(repo, nil)

	return &svc
}

func cleanDB(t *testing.T, db *sql.DB, tables ...string) {
	for _, table := range tables {
		_, err := db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE;")
		if err != nil {
			t.Fatalf("Failed to truncate table %s: %v", table, err)
		}
	}
}

func createDummyCourse(t *testing.T, db *sql.DB) int {
	baseBody := request.CreateCourse{
		CourseName:      "Advanced Go Programming",
		CourseCode:      "01076203",
		ProfessorID:     1,
		CourseProgramID: 1,
		CourseProgram:   "General",
		Sec:             "101",
		SemesterID:      1,
		Semester:        "2/2568",
		ClassdayID:      1,
		Classday:        "Monday",
		ClassStart:      time.Date(2026, time.January, 1, 9, 0, 0, 0, time.UTC).Format("15:04:05"),
		ClassEnd:        time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC).Format("15:04:05"),
		WorkHour:        3,
		CreatedDate:     time.Now(),
	}

	var id int
	// Updated query to include all NOT NULL fields from your schema
	query := `
        INSERT INTO courses (
            course_name, course_code, professor_ID, course_program_ID, 
            course_program, sec, semester_ID, semester, 
            class_day_ID, class_day, class_start, class_end, 
            work_hour, created_date
        ) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) 
        RETURNING course_ID`

	err := db.QueryRow(query,
		baseBody.CourseName,
		baseBody.CourseCode,
		baseBody.ProfessorID,
		baseBody.CourseProgramID,
		baseBody.CourseProgram,
		baseBody.Sec,
		baseBody.SemesterID,
		baseBody.Semester,
		baseBody.ClassdayID,
		baseBody.Classday,
		baseBody.ClassStart,
		baseBody.ClassEnd,
		baseBody.WorkHour,
		baseBody.CreatedDate,
	).Scan(&id)

	assert.NoError(t, err)
	return id
}

func createDummyJobPost(t *testing.T, db *sql.DB, courseID int, taAllocation int) int {
	// Defining dummy data based on your schema requirements
	jobData := request.CreateJobPost{
		ProfessorID:  1, // Must exist in professors table
		Task:         "Assist in Advanced Go Programming labs and grading.",
		TaAllocation: taAllocation,
		Location:     "Engineering Bldg 4",
		CourseID:     courseID,
		GradeID:      4, // e.g., 4 = Grade A requirement
		CreatedDate:  time.Now(),
	}

	var id int
	query := `
        INSERT INTO ta_job_posting (
            professor_ID, task, ta_allocation, location, 
            status_ID, course_ID, grade_ID, created_date
        ) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
        RETURNING id`

	err := db.QueryRow(query,
		jobData.ProfessorID,
		jobData.Task,
		jobData.TaAllocation,
		jobData.Location,
		1,
		jobData.CourseID,
		jobData.GradeID,
		jobData.CreatedDate,
	).Scan(&id)

	assert.NoError(t, err)
	return id
}

func createDummyApplication(t *testing.T, db *sql.DB, studentID int, jobPostID int) int {
	var applicationId int

	// Use the same status ID your service uses
	// If you don't have access to constants here, use the raw ID (e.g., 1)
	statusId := constants.PendingStatusID

	// We use a dummy grade and purpose for the setup
	grade := "A"
	purpose := "Experience"

	query := `INSERT INTO ta_application(
                    student_ID, 
                    status_ID, 
                    job_post_ID,
                    grade,
                    purpose,
                    created_date)
                VALUES($1, $2, $3, $4, $5, $6)
                RETURNING id`

	err := db.QueryRow(query,
		studentID,
		statusId,
		jobPostID,
		grade,
		purpose,
		time.Now(),
	).Scan(&applicationId)

	if err != nil {
		t.Fatalf("Setup Failed: failed to insert dummy application: %v", err)
	}

	return applicationId
}

func createDummyApprove(t *testing.T, db *sql.DB, appId int) {
	// 1. Get the necessary IDs to make a valid TA entry
	var studentId, courseId int
	err := db.QueryRow(`
        SELECT a.student_id, j.course_id 
        FROM ta_application a 
        JOIN ta_job_posting j ON a.job_post_id = j.id 
        WHERE a.id = $1`, appId).Scan(&studentId, &courseId)

	if err != nil {
		t.Fatalf("Failed to get data for dummy approve: %v", err)
	}

	// 2. Update status and insert into ta_courses manually
	_, err = db.Exec("UPDATE ta_application SET status_ID = $1 WHERE id = $2", constants.ApprovedStatusID, appId)
	if err != nil {
		t.Fatalf("Failed to update status in dummy: %v", err)
	}

	_, err = db.Exec("INSERT INTO ta_courses (student_id, course_id) VALUES ($1, $2)", studentId, courseId)
	if err != nil {
		t.Fatalf("Failed to insert ta_course in dummy: %v", err)
	}
}

func verifyCourseInDB(t *testing.T, db *sql.DB, input request.CreateCourse, courseID int) {

	var actualCourseID int
	var profID int
	query := "SELECT course_ID, professor_ID FROM courses WHERE course_ID = $1"

	err := db.QueryRow(query, courseID).Scan(&actualCourseID, &profID)

	assert.NoError(t, err, "Data should exist in database")
	assert.Equal(t, courseID, actualCourseID)
	assert.Equal(t, input.ProfessorID, profID)
}

func verifyApplicationInDB(t *testing.T, db *sql.DB, applicationID int, jobPostID int) {
	var count int

	query := "SELECT COUNT(*) FROM ta_application WHERE id=$1 AND job_post_ID=$2"
	err := db.QueryRow(query, applicationID, jobPostID).Scan(&count)

	assert.Nil(t, err)
	assert.Equal(t, 1, count)
}

func verifyApproveInDB(t *testing.T, db *sql.DB, appId int) {
	var status int
	err := db.QueryRow("SELECT status_ID FROM ta_application WHERE id = $1", appId).Scan(&status)
	if err != nil {
		t.Fatalf("could not find application in DB: %v", err)
	}

	approveStatus := constants.ApprovedStatusID
	assert.Equal(t, approveStatus, status, "Application status should be updated to approved in DB")

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM ta_courses WHERE student_id = (SELECT student_id FROM ta_application WHERE id = $1))", appId).Scan(&exists)
	assert.NoError(t, err)
	assert.True(t, exists, fmt.Sprintf("Student should be present in ta_courses table status:%d", status))
}

func TestCreateCourse_Integration(t *testing.T) {
	// *** CRITICAL ASSUMPTION: testDB must be initialized by TestMain before this runs. ***

	// Check if testDB is available (a safety check)
	if testDB == nil {
		t.Fatal("Test database connection (testDB) is not initialized. Check your TestMain execution.")
	}

	// 1. Setup the real service and clean the database
	svc := setupServiceTest(t, testDB)
	// Reuse the cleanup helper (assuming it's available)

	// Define the common request body
	baseBody := request.CreateCourse{
		CourseName:      "Advanced Go Programming",
		CourseCode:      "01076203",
		ProfessorID:     1,
		CourseProgramID: 1,
		CourseProgram:   "General",
		Sec:             "101",
		SemesterID:      1,
		Semester:        "2/2568",
		ClassdayID:      1,
		Classday:        "Monday",
		// Note: Using fixed dates or a single time.Time instance is better for testing
		ClassStart:  time.Date(2026, time.January, 1, 9, 0, 0, 0, time.UTC).Format("15:04:05"),
		ClassEnd:    time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC).Format("15:04:05"),
		WorkHour:    3,
		CreatedDate: time.Now(),
	}

	// 2. Define the Test Table Scenarios
	tests := []struct {
		name          string
		body          request.CreateCourse
		setupAction   func() // Action to prepare the DB before the primary ACT
		expectedError string
	}{
		{
			name:        "IT001_CreateSuccess",
			body:        baseBody,
			setupAction: func() { cleanDB(t, testDB, "courses") }, /* DB is clean, nothing else to do */
		},
		{
			name: "IT002_DuplicateCourseFailure",
			body: baseBody, // Attempt to insert the same course again
			setupAction: func() {
				cleanDB(t, testDB, "courses")
				_, _ = svc.CreateCourse(baseBody)
			},
			expectedError: "already have this course", // Should fail on the second attempt due to DB constraint
		},
	}

	// 3. Run the Tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure a clean slate for each test case

			// Run the specific setup action for this test case
			tt.setupAction()

			// ACT: Call the Service method (the entry point for this test)
			res, err := svc.CreateCourse(tt.body)

			// ASSERT 1: Check the error result from the Service
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)

				verifyCourseInDB(t, testDB, tt.body, res.Id)
			}
		})
	}
}

func TestApplyJobPost_Integration(t *testing.T) {
	if testDB == nil {
		t.Fatal("Test database connection (testDB) is not initialized. Check your TestMain execution.")
	}

	// 1. Setup the real service and clean the database
	svc := setupServiceTest(t, testDB)

	tests := []struct {
		name          string
		expectedError string
		setUp         func() int
	}{
		{

			name: "IT003_ApplyJobSuccess",
			setUp: func() int {
				cleanDB(t, testDB, "ta_application", "ta_job_posting", "courses")
				courseID := createDummyCourse(t, testDB)
				return createDummyJobPost(t, testDB, courseID, 1)
			},
		},
		{
			name:          "IT004_ApplyJobDuplicate",
			expectedError: "already apply to this job",
			setUp: func() int {
				cleanDB(t, testDB, "ta_application", "ta_job_posting", "courses")

				courseID := createDummyCourse(t, testDB)
				jobPostID := createDummyJobPost(t, testDB, courseID, 1)
				createDummyApplication(t, testDB, 1, jobPostID)
				return jobPostID
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			jobpostID := tt.setUp()

			fNameThai := "pittawat thai"
			lNameThai := "kitmong thai"
			transcript := []byte("hello world")
			baseBody := request.ApplyJobPost{
				JobPostID:        &jobpostID,
				StudentID:        1,
				FirstnameThai:    &fNameThai,
				LastnameThai:     &lNameThai,
				TranscriptBytes:  &transcript,
				BankAccountBytes: &transcript,
				StudentCardBytes: &transcript,
			}

			res, err := svc.ApplyJobPost(baseBody)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				verifyApplicationInDB(t, testDB, res.Id, jobpostID)
			}

		})
	}

}

func TestApproveApplication_Integration(t *testing.T) {
	if testDB == nil {
		t.Fatal("Test database connection (testDB) is not initialized. Check your TestMain execution.")
	}

	// 1. Setup the real service and clean the database
	svc := setupServiceTest(t, testDB)

	tests := []struct {
		name          string
		expectedError string
		setUp         func() int
	}{
		{

			name: "IT005_ApproveSuccess",
			setUp: func() int {
				cleanDB(t, testDB, "ta_application", "ta_job_posting", "courses", "ta_courses")
				courseID := createDummyCourse(t, testDB)
				jonPostID := createDummyJobPost(t, testDB, courseID, 1)
				appID := createDummyApplication(t, testDB, 1, jonPostID)

				var exists bool
				testDB.QueryRow("SELECT EXISTS(SELECT 1 FROM ta_application WHERE id=$1)", appID).Scan(&exists)
				fmt.Printf("DEBUG: AppID %d exists after setup: %v\n", appID, exists)
				return appID
			},
		},
		{
			name:          "IT006_TaAllocationFull",
			expectedError: "ta allocation is full",
			setUp: func() int {
				cleanDB(t, testDB, "ta_application", "ta_job_posting", "courses")

				courseID := createDummyCourse(t, testDB)
				jobPostID := createDummyJobPost(t, testDB, courseID, 1)
				applicationID := createDummyApplication(t, testDB, 1, jobPostID)
				applicationID2 := createDummyApplication(t, testDB, 2, jobPostID)
				createDummyApprove(t, testDB, applicationID2)
				return applicationID
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			applicationID := tt.setUp()

			res, err := svc.ApproveApplication(applicationID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				verifyApproveInDB(t, testDB, applicationID)
			}

		})
	}

}

// func TestService_UpdateCourse_Integration(t *testing.T) {
// 	if testDB == nil {
// 		t.Fatal("Test database connection (testDB) is not initialized. Check your TestMain execution.")
// 	}

// 	// 1. Setup the real service and clean the database
// 	svc := setupServiceTest(t, testDB)
// 	cleanDB(t, testDB, "courses")

// 	baseBody := request.CreateCourse{
// 		CourseName:      "Advanced Go Programming",
// 		CourseCode:      "01076203",
// 		ProfessorID:     1,
// 		CourseProgramID: 1,
// 		CourseProgram:   "General",
// 		Sec:             "101",
// 		SemesterID:      1,
// 		Semester:        "2/2568",
// 		ClassdayID:      1,
// 		Classday:        "Monday",
// 		// Note: Using fixed dates or a single time.Time instance is better for testing
// 		ClassStart:  " time.Date(2026, time.January, 1, 9, 0, 0, 0, time.UTC)",
// 		ClassEnd:    "time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)",
// 		CreatedDate: time.Now(),
// 	}

// 	createResult, err := svc.CreateCourse(baseBody)
// 	if err != nil {
// 		t.Fatalf("Failed to create initial course: %v", err)
// 	}

// 	newCorseName := "programing fundamental"
// 	newCourseCode := "1123456"
// 	updateBody := request.UpdateCourse{
// 		CourseName: &newCorseName,
// 		CourseCode: &newCourseCode,
// 		Id:         createResult.Id,
// 	}

// 	_, err = svc.UpdateCourse(updateBody)
// 	if err != nil {
// 		t.Fatalf("Update course Failed : %v", err)
// 	}

// 	var updatedName string
// 	var updatedId string

// 	err = testDB.QueryRow("SELECT course_name, course_id FROM courses WHERE id = $1 ", createResult.Id).
// 		Scan(&updatedName, &updatedId)

// 	if err != nil {
// 		t.Fatalf("Could not find record after update: %v", err)
// 	}

// 	if updatedName != newCorseName {
// 		t.Errorf("Course mismatch. Want: %s ,got: %s", newCorseName, updatedName)
// 	}
// 	if updatedId != newCourseCode {
// 		t.Errorf("CourseCode mismatch. Want:%s ,got: %s", newCourseCode, updatedId)
// 	}

// }
// func TestService_DeleteCourse_Integration(t *testing.T) {
// 	if testDB == nil {
// 		t.Fatal("Test database connection (testDB) is not initialized. Check your TestMain execution.")
// 	}

// 	// 1. Setup the real service and clean the database
// 	svc := setupServiceTest(t, testDB)
// 	cleanDB(t, testDB, "courses")

// 	baseBody := request.CreateCourse{
// 		CourseName:      "Advanced Go Programming",
// 		CourseCode:      "01076203",
// 		ProfessorID:     1,
// 		CourseProgramID: 1,
// 		CourseProgram:   "General",
// 		Sec:             "101",
// 		SemesterID:      1,
// 		Semester:        "2/2568",
// 		ClassdayID:      1,
// 		Classday:        "Monday",
// 		// Note: Using fixed dates or a single time.Time instance is better for testing
// 		ClassStart:  " time.Date(2026, time.January, 1, 9, 0, 0, 0, time.UTC)",
// 		ClassEnd:    "time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)",
// 		CreatedDate: time.Now(),
// 	}

// 	createResult, err := svc.CreateCourse(baseBody)
// 	if err != nil {
// 		t.Fatalf("Failed to create initial course: %v", err)
// 	}

// 	_, err = svc.DeleteCourse(createResult.Id)
// 	if err != nil {
// 		t.Fatalf("Failed to delete course: %v", err)
// 	}

// 	var count int
// 	err = testDB.QueryRow("SELECT COUNT(*) FROM courses WHERE id = $1 AND deleted_date = NULL", createResult.Id).Scan(&count)

// 	if err != nil {
// 		t.Fatalf("Failed to find course after delete: %v", err)
// 	}
// 	if count != 0 {
// 		t.Errorf("already Deleted but still appear! ")
// 	}
// }

// func TestService_ApplyCourse_Integration(t *testing.T) {
// 	if testDB == nil {
// 		t.Fatal("Test database connection (testDB) is not initialized. Check your TestMain execution.")
// 	}

// 	// 1. Setup the real service and clean the database
// 	svc := setupServiceTest(t, testDB)
// 	cleanDB(t, testDB, "courses")

// 	baseBody := request.CreateCourse{
// 		CourseName:      "Advanced Go Programming",
// 		CourseCode:      "01076203",
// 		ProfessorID:     1,
// 		CourseProgramID: 1,
// 		CourseProgram:   "General",
// 		Sec:             "101",
// 		SemesterID:      1,
// 		Semester:        "2/2568",
// 		ClassdayID:      1,
// 		Classday:        "Monday",
// 		// Note: Using fixed dates or a single time.Time instance is better for testing
// 		ClassStart:  " time.Date(2026, time.January, 1, 9, 0, 0, 0, time.UTC)",
// 		ClassEnd:    "time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)",
// 		CreatedDate: time.Now(),
// 	}

// 	_, err := svc.CreateCourse(baseBody)
// 	if err != nil {
// 		t.Fatalf("Failed to create initial course: %v", err)
// 	}

// 	filePath := "../../../../test_assets/sample.pdf"
// 	fileBytes, err := os.ReadFile(filePath)
// 	if err != nil {
// 		t.Fatalf("Failed to readfile : %v", err)
// 	}
// 	courseId := 1
// 	filename := "sample.pdf"
// 	applyBody := request.ApplyJobPost{
// 		StudentID:       12345,
// 		JobPostID:       &courseId,
// 		TranscriptBytes: &fileBytes,
// 		TranscriptName:  &filename,
// 	}

// 	result, err := svc.ApplyJobPost(applyBody)
// 	if err != nil {
// 		t.Errorf("ApplyCourse failed: %v", err)
// 	}

// 	if result == nil {
// 		t.Error("Expected result to be non-nil")
// 	}
// }

// func TestService_GetApplicationByStudentId_Integration(t *testing.T) {
// 	if testDB == nil {
// 		t.Fatal("Test database connection (testDB) is not initialized. Check your TestMain execution.")
// 	}

// 	// 1. Setup the real service and clean the database
// 	svc := setupServiceTest(t, testDB)
// 	cleanDB(t, testDB, "courses")

// 	baseBody := request.CreateCourse{
// 		CourseName:      "Advanced Go Programming",
// 		CourseCode:      "01076203",
// 		ProfessorID:     1,
// 		CourseProgramID: 1,
// 		CourseProgram:   "General",
// 		Sec:             "101",
// 		SemesterID:      1,
// 		Semester:        "2/2568",
// 		ClassdayID:      1,
// 		Classday:        "Monday",
// 		// Note: Using fixed dates or a single time.Time instance is better for testing
// 		ClassStart:  " time.Date(2026, time.January, 1, 9, 0, 0, 0, time.UTC)",
// 		ClassEnd:    "time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)",
// 		CreatedDate: time.Now(),
// 	}

// 	_, err := svc.CreateCourse(baseBody)
// 	if err != nil {
// 		t.Fatalf("Failed to create initial course: %v", err)
// 	}

// 	filePath := "../../../../test_assets/sample.pdf"
// 	fileBytes, err := os.ReadFile(filePath)
// 	if err != nil {
// 		t.Fatalf("Failed to readfile : %v", err)
// 	}
// 	courseId := 1
// 	filename := "sample.pdf"
// 	studentId := 12345
// 	applyBody := request.ApplyJobPost{
// 		StudentID:       studentId,
// 		JobPostID:       &courseId,
// 		TranscriptBytes: &fileBytes,
// 		TranscriptName:  &filename,
// 	}

// 	result, err := svc.ApplyJobPost(applyBody)
// 	if err != nil {
// 		t.Errorf("ApplyCourse failed: %v", err)
// 	}

// 	if result == nil {
// 		t.Error("Expected result to be non-nil")
// 	}

// 	applications, err := svc.GetApplicationByStudentId(studentId)
// 	if err != nil {
// 		t.Fatalf("Failed to getApplicationByStudentID : %v", err)
// 	}

// 	if applications.Data == nil {
// 		t.Error("Should have application record.")
// 	}
// }

// func TestService_GetApplicationByCourseId_Integration(t *testing.T) {
// 	if testDB == nil {
// 		t.Fatal("Test database connection (testDB) is not initialized. Check your TestMain execution.")
// 	}

// 	// 1. Setup the real service and clean the database
// 	svc := setupServiceTest(t, testDB)
// 	cleanDB(t, testDB, "courses")

// 	baseBody := request.CreateCourse{
// 		CourseName:      "Advanced Go Programming",
// 		CourseCode:      "01076203",
// 		ProfessorID:     1,
// 		CourseProgramID: 1,
// 		CourseProgram:   "General",
// 		Sec:             "101",
// 		SemesterID:      1,
// 		Semester:        "2/2568",
// 		ClassdayID:      1,
// 		Classday:        "Monday",
// 		// Note: Using fixed dates or a single time.Time instance is better for testing
// 		ClassStart:  " time.Date(2026, time.January, 1, 9, 0, 0, 0, time.UTC)",
// 		ClassEnd:    "time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)",
// 		CreatedDate: time.Now(),
// 	}

// 	_, err := svc.CreateCourse(baseBody)
// 	if err != nil {
// 		t.Fatalf("Failed to create initial course: %v", err)
// 	}

// 	filePath := "../../../../test_assets/sample.pdf"
// 	fileBytes, err := os.ReadFile(filePath)
// 	if err != nil {
// 		t.Fatalf("Failed to readfile : %v", err)
// 	}
// 	courseId := 1
// 	filename := "sample.pdf"
// 	applyBody := request.ApplyJobPost{
// 		StudentID:       12345,
// 		JobPostID:       &courseId,
// 		TranscriptBytes: &fileBytes,
// 		TranscriptName:  &filename,
// 	}

// 	result, err := svc.ApplyJobPost(applyBody)
// 	if err != nil {
// 		t.Errorf("ApplyCourse failed: %v", err)
// 	}

// 	if result == nil {
// 		t.Error("Expected result to be non-nil")
// 	}

// 	applications, err := svc.GetApplicationByCourseId(courseId)
// 	if err != nil {
// 		t.Fatalf("Failed to getApplicationByStudentID : %v", err)
// 	}

// 	if applications.Data == nil {
// 		t.Error("Should have application record.")
// 	}
// }

// func TestService_GetApplicationDetail_Integration(t *testing.T) {
// 	if testDB == nil {
// 		t.Fatal("Test database connection (testDB) is not initialized. Check your TestMain execution.")
// 	}

// 	// 1. Setup the real service and clean the database
// 	svc := setupServiceTest(t, testDB)
// 	cleanDB(t, testDB, "courses")

// 	baseBody := request.CreateCourse{
// 		CourseName:      "Advanced Go Programming",
// 		CourseCode:      "01076203",
// 		ProfessorID:     1,
// 		CourseProgramID: 1,
// 		CourseProgram:   "General",
// 		Sec:             "101",
// 		SemesterID:      1,
// 		Semester:        "2/2568",
// 		ClassdayID:      1,
// 		Classday:        "Monday",
// 		// Note: Using fixed dates or a single time.Time instance is better for testing
// 		ClassStart:  " time.Date(2026, time.January, 1, 9, 0, 0, 0, time.UTC)",
// 		ClassEnd:    "time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)",
// 		CreatedDate: time.Now(),
// 	}

// 	_, err := svc.CreateCourse(baseBody)
// 	if err != nil {
// 		t.Fatalf("Failed to create initial course: %v", err)
// 	}

// 	filePath := "../../../../test_assets/sample.pdf"
// 	fileBytes, err := os.ReadFile(filePath)
// 	if err != nil {
// 		t.Fatalf("Failed to readfile : %v", err)
// 	}
// 	courseId := 1
// 	filename := "sample.pdf"
// 	applyBody := request.ApplyJobPost{
// 		StudentID:       12345,
// 		JobPostID:       &courseId,
// 		TranscriptBytes: &fileBytes,
// 		TranscriptName:  &filename,
// 	}

// 	result, err := svc.ApplyJobPost(applyBody)
// 	if err != nil {
// 		t.Errorf("ApplyCourse failed: %v", err)
// 	}

// 	if result == nil {
// 		t.Error("Expected result to be non-nil")
// 	}

// 	applications, err := svc.GetApplicationByCourseId(result.Id)
// 	if err != nil {
// 		t.Fatalf("Failed to getApplicationByStudentID : %v", err)
// 	}

// 	if applications.Data == nil {
// 		t.Error("Should have application record.")
// 	}
// }
