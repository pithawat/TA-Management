// File: TA-management/internal/modules/course/service/course_service_integration_test.go

package service_test

import (
	"TA-management/internal/modules/course/dto/request"
	"TA-management/internal/modules/course/repository"
	"TA-management/internal/modules/course/service" // Assuming this is the service location
	"TA-management/internal/modules/shared/dto/testutils"
	"database/sql"
	"os"
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
func TestMain(m *testing.M) {

	testDB = testutils.InitTestDB()

	exitCode := m.Run()

	testDB.Close()
	os.Exit(exitCode)
}

// setupServiceTest initializes the real repository and the real service.
func setupServiceTest(t *testing.T, db *sql.DB) *service.CourseServiceImplementation {
	// 1. Initialize the Real Repository with the Real Test DB connection
	repo := repository.NewCourseRepository(db)

	// 2. Initialize the Service with the Real Repository
	// Assuming a constructor like: NewCourseService(repo repository.CourseRepository)
	svc := service.NewCourseService(repo)

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

func TestService_CreateCourse_Integration(t *testing.T) {
	// *** CRITICAL ASSUMPTION: testDB must be initialized by TestMain before this runs. ***

	// Check if testDB is available (a safety check)
	if testDB == nil {
		t.Fatal("Test database connection (testDB) is not initialized. Check your TestMain execution.")
	}

	// 1. Setup the real service and clean the database
	svc := setupServiceTest(t, testDB)
	cleanDB(t, testDB, "courses") // Reuse the cleanup helper (assuming it's available)

	// Define the common request body
	baseBody := request.CreateCourse{
		CourseName:      "Advanced Go Programming",
		CourseID:        "01076203",
		ProfessorID:     1,
		CourseProgramID: 1,
		CourseProgram:   "General",
		Sec:             "101",
		SemesterID:      1,
		Semester:        "2/2568",
		ClassdayID:      1,
		Classday:        "Monday",
		// Note: Using fixed dates or a single time.Time instance is better for testing
		ClassStart:  time.Date(2026, time.January, 1, 9, 0, 0, 0, time.UTC),
		ClassEnd:    time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC),
		CreatedDate: time.Now(),
	}

	// 2. Define the Test Table Scenarios
	tests := []struct {
		name          string
		body          request.CreateCourse
		setupAction   func(t *testing.T) // Action to prepare the DB before the primary ACT
		expectedError bool
	}{
		{
			name:          "Success_FullFlow",
			body:          baseBody,
			setupAction:   func(t *testing.T) { /* DB is clean, nothing else to do */ },
			expectedError: false,
		},
		{
			name: "Failure_DuplicateCourseID",
			body: baseBody, // Attempt to insert the same course again
			setupAction: func(t *testing.T) {
				// Pre-insert the data using the SERVICE to test the end-to-end failure path
				_, err := svc.CreateCourse(baseBody)
				assert.Nil(t, err, "Setup action failed to pre-insert initial course")
			},
			expectedError: true, // Should fail on the second attempt due to DB constraint
		},
	}

	// 3. Run the Tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure a clean slate for each test case
			cleanDB(t, testDB, "courses")

			// Run the specific setup action for this test case
			tt.setupAction(t)

			// ACT: Call the Service method (the entry point for this test)
			_, err := svc.CreateCourse(tt.body)

			// ASSERT 1: Check the error result from the Service
			if tt.expectedError {
				assert.NotNil(t, err, "Expected an error but got nil")
				// In a real scenario, you'd check the type of error returned by the service
				// e.g., assert.ErrorIs(t, err, service.ErrDuplicate)
			} else {
				assert.Nil(t, err, "Expected nil error on success")

				// ASSERT 2 (Crucial Integration Check): Verify the state in the real database
				var courseID string
				var profID int
				queryErr := testDB.QueryRow("SELECT course_ID, professor_ID FROM courses WHERE course_ID = $1", tt.body.CourseID).Scan(&courseID, &profID)

				assert.Nil(t, queryErr, "Failed to query row after successful creation")
				assert.Equal(t, tt.body.CourseID, courseID)
				assert.Equal(t, tt.body.ProfessorID, profID)
			}
		})
	}
}

func TestService_UpdateCourse_Integration(t *testing.T) {
	if testDB == nil {
		t.Fatal("Test database connection (testDB) is not initialized. Check your TestMain execution.")
	}

	// 1. Setup the real service and clean the database
	svc := setupServiceTest(t, testDB)
	cleanDB(t, testDB, "courses")

	baseBody := request.CreateCourse{
		CourseName:      "Advanced Go Programming",
		CourseID:        "01076203",
		ProfessorID:     1,
		CourseProgramID: 1,
		CourseProgram:   "General",
		Sec:             "101",
		SemesterID:      1,
		Semester:        "2/2568",
		ClassdayID:      1,
		Classday:        "Monday",
		// Note: Using fixed dates or a single time.Time instance is better for testing
		ClassStart:  time.Date(2026, time.January, 1, 9, 0, 0, 0, time.UTC),
		ClassEnd:    time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC),
		CreatedDate: time.Now(),
	}

	createResult, err := svc.CreateCourse(baseBody)
	if err != nil {
		t.Fatalf("Failed to create initial course: %v", err)
	}

	newCorseName := "programing fundamental"
	newCourseId := "1123456"
	updateBody := request.UpdateCourse{
		CourseName: &newCorseName,
		CourseID:   &newCourseId,
		Id:         createResult.Id,
	}

	_, err = svc.UpdateCourse(updateBody)
	if err != nil {
		t.Fatalf("Update course Failed : %v", err)
	}

	var updatedName string
	var updatedId string

	err = testDB.QueryRow("SELECT course_name, course_id FROM courses WHERE id = $1 ", createResult.Id).
		Scan(&updatedName, &updatedId)

	if err != nil {
		t.Fatalf("Could not find record after update: %v", err)
	}

	if updatedName != newCorseName {
		t.Errorf("Course mismatch. Want: %s ,got: %s", newCorseName, updatedName)
	}
	if updatedId != newCourseId {
		t.Errorf("CourseId mismatch. Want:%s ,got: %s", newCourseId, updatedId)
	}

}
func TestService_DeleteCourse_Integration(t *testing.T) {
	if testDB == nil {
		t.Fatal("Test database connection (testDB) is not initialized. Check your TestMain execution.")
	}

	// 1. Setup the real service and clean the database
	svc := setupServiceTest(t, testDB)
	cleanDB(t, testDB, "courses")

	baseBody := request.CreateCourse{
		CourseName:      "Advanced Go Programming",
		CourseID:        "01076203",
		ProfessorID:     1,
		CourseProgramID: 1,
		CourseProgram:   "General",
		Sec:             "101",
		SemesterID:      1,
		Semester:        "2/2568",
		ClassdayID:      1,
		Classday:        "Monday",
		// Note: Using fixed dates or a single time.Time instance is better for testing
		ClassStart:  time.Date(2026, time.January, 1, 9, 0, 0, 0, time.UTC),
		ClassEnd:    time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC),
		CreatedDate: time.Now(),
	}

	createResult, err := svc.CreateCourse(baseBody)
	if err != nil {
		t.Fatalf("Failed to create initial course: %v", err)
	}

	_, err = svc.DeleteCourse(createResult.Id)
	if err != nil {
		t.Fatalf("Failed to delete course: %v", err)
	}

	var count int
	err = testDB.QueryRow("SELECT COUNT(*) FROM courses WHERE id = $1 AND deleted_date = NULL", createResult.Id).Scan(&count)

	if err != nil {
		t.Fatalf("Failed to find course after delete: %v", err)
	}
	if count != 0 {
		t.Errorf("already Deleted but still appear! ")
	}
}
