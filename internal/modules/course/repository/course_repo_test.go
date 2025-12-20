package repository_test

import (
	"TA-management/internal/modules/course/dto/request"
	"TA-management/internal/modules/course/repository"
	"TA-management/internal/modules/shared/dto/testutils"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var testDB *sql.DB

func TestMain(m *testing.M) {

	testDB = testutils.InitTestDB()

	exitCode := m.Run()

	testDB.Close()
	os.Exit(exitCode)
}

func cleanDB(t *testing.T, db *sql.DB, tables ...string) {
	for _, table := range tables {
		_, err := db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE;")
		if err != nil {
			t.Fatalf("Failed to truncate table %s: %v", table, err)
		}
	}
}

func TestRepo_CreateCourse_Success(t *testing.T) {
	cleanDB(t, testDB, "courses")
	repo := repository.NewCourseRepository(testDB)

	inputBody := request.CreateCourse{
		CourseName:      "program fundamental 1",
		CourseID:        "01076103",
		ProfessorID:     1,
		CourseProgramID: 1,
		CourseProgram:   "General",
		Sec:             "101",
		SemesterID:      1,
		Semester:        "2/2568",
		ClassdayID:      1,
		Classday:        "Monday",
		ClassStart:      "13:00",
		ClassEnd:        "17:00",
		CreatedDate:     time.Now(),
	}

	_, err := repo.CreateCourse(inputBody)

	// ASSERT 1: No error from the repository call
	assert.Nil(t, err)

	// ASSERT 2: Verify the data exists in the database
	var courseID string
	var profID int
	err = testDB.QueryRow("SELECT course_ID, professor_ID FROM courses WHERE course_ID=$1", "01076103").Scan(&courseID, &profID)

	assert.Nil(t, err)
	assert.Equal(t, "01076103", courseID)
	assert.Equal(t, 1, profID)
}

func TestRepo_CreateCourse_FailsOnDuplicate(t *testing.T) {
	cleanDB(t, testDB, "courses")

	repo := repository.NewCourseRepository(testDB)
	inputBody := request.CreateCourse{
		CourseName:      "program fundamental 1",
		CourseID:        "01076103",
		ProfessorID:     1,
		CourseProgramID: 1,
		CourseProgram:   "General",
		Sec:             "101",
		SemesterID:      1,
		Semester:        "2/2568",
		ClassdayID:      1,
		Classday:        "Monday",
		ClassStart:      "13:00",
		ClassEnd:        "14:00",
		CreatedDate:     time.Now(),
	}
	_, err := repo.CreateCourse(inputBody)
	assert.Nil(t, err, "First insert should success")

	_, err = repo.CreateCourse(inputBody)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "course already exists")
}
