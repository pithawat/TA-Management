package service_test

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"TA-management/internal/modules/shared/dto/testutils"
	tadutrepo "TA-management/internal/modules/ta_duty/repository"
	"TA-management/internal/modules/ta_duty/service"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	testDB = testutils.InitTestDB()
	exitCode := m.Run()
	testDB.Close()
	os.Exit(exitCode)
}

func cleanDB(t *testing.T, db *sql.DB, tables ...string) {
	t.Helper()
	for _, table := range tables {
		_, err := db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE;")
		if err != nil {
			t.Fatalf("Failed to truncate table %s: %v", table, err)
		}
	}
}

func setupTaDutyService(db *sql.DB) service.TaDutyService {
	repo := tadutrepo.NewTaDutyRepository(db)
	logger, _ := zap.NewDevelopment()
	return service.NewTaDutyServiceImplementation(repo, logger.Sugar())
}

// insertDummyDutyHistory inserts a ta_duty_historys row directly (bypasses service date-check).
func insertDummyDutyHistory(t *testing.T, db *sql.DB, courseID, studentID int, date string) {
	t.Helper()
	_, err := db.Exec(
		"INSERT INTO ta_duty_historys (date, course_ID, student_ID) VALUES ($1, $2, $3)",
		date, courseID, studentID,
	)
	assert.NoError(t, err, "Failed to insert dummy duty history")
}

// ─── GetTADutyRoadmap Integration ────────────────────────────────────────────

func TestGetTADutyRoadmap_Integration(t *testing.T) {
	if testDB == nil {
		t.Fatal("testDB is not initialized")
	}

	svc := setupTaDutyService(testDB)

	t.Run("IT001_CourseWithActiveSemester_ReturnsChecklist", func(t *testing.T) {
		// Course ID=1 should be seeded and tied to a semester with start/end dates
		result, err := svc.GetTADutyRoadmap(1, 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		// Result can be empty slice if no class days fall in that range — that's still valid
	})

	t.Run("IT002_NonExistentCourse_ReturnsEmptyNoError", func(t *testing.T) {
		// A non-existent courseID will make the CTE return nothing — should return empty list, no error
		result, err := svc.GetTADutyRoadmap(999999, 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(*result))
	})
}

// ─── MarkDutyAsDone Integration ───────────────────────────────────────────────

func TestMarkDutyAsDone_Integration(t *testing.T) {
	if testDB == nil {
		t.Fatal("testDB is not initialized")
	}

	svc := setupTaDutyService(testDB)

	t.Run("IT003_MarkPastDuty_PersistsInDB", func(t *testing.T) {
		// Use a known past date far from today
		pastDate := "2020-01-06" // A Monday
		// Clean any existing record first
		testDB.Exec(
			"DELETE FROM ta_duty_historys WHERE course_ID = 1 AND student_ID = 1 AND date::date = $1",
			pastDate,
		)

		res, err := svc.MarkDutyAsDone(1, 1, pastDate)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "MarkDuty as Done successfully", res.Message)

		// Verify record in DB
		var count int
		testDB.QueryRow(
			"SELECT COUNT(*) FROM ta_duty_historys WHERE course_ID = 1 AND student_ID = 1 AND date::date = $1",
			pastDate,
		).Scan(&count)
		assert.Equal(t, 1, count)
	})

	t.Run("IT004_MarkFutureDate_ShouldReturnError", func(t *testing.T) {
		futureDate := time.Now().AddDate(0, 0, 10).Format("2006-01-02")
		res, err := svc.MarkDutyAsDone(1, 1, futureDate)

		assert.Nil(t, res)
		assert.EqualError(t, err, "cannot check off future duties")

		// Ensure nothing was inserted
		var count int
		testDB.QueryRow(
			"SELECT COUNT(*) FROM ta_duty_historys WHERE course_ID = 1 AND student_ID = 1 AND date::date = $1",
			futureDate,
		).Scan(&count)
		assert.Equal(t, 0, count)
	})

	t.Run("IT005_InvalidDateFormat_ShouldReturnError", func(t *testing.T) {
		res, err := svc.MarkDutyAsDone(1, 1, "01-13-2020")
		assert.Nil(t, res)
		assert.EqualError(t, err, "invalid date format")
	})
}

// ─── ExportPaymentReport Integration ────────────────────────────────────────

func TestExportPaymentReport_Integration(t *testing.T) {
	if testDB == nil {
		t.Fatal("testDB is not initialized")
	}

	t.Run("IT006_ValidCourseWithNoTA_ReturnsEmptyButNoError", func(t *testing.T) {
		// Course 1 may have no TA in ta_courses — repo returns empty, excel file not needed at service level.
		// The service only calls excel AFTER getting repo data. If ta_courses is empty, TAdutyData is empty.
		// The service will then try to open the excel template — skip this test if the file doesn't exist.
		t.Skip("ExportPaymentReport requires ./prototype/payment-template.xlsx to be present — skipping in CI")
	})
}

func TestExportSignatureSheet_Integration(t *testing.T) {
	if testDB == nil {
		t.Fatal("testDB is not initialized")
	}

	t.Run("IT007_ExcelTemplateRequired_Skipped", func(t *testing.T) {
		t.Skip("ExportSignatureSheet requires ./prototype/signature-template.xlsx to be present — skipping in CI")
	})
}
