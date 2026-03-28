package service_test

import (
	"database/sql"
	"os"
	"testing"

	annrequest "TA-management/internal/modules/announce/dto/request"
	annrepo "TA-management/internal/modules/announce/repository"
	"TA-management/internal/modules/announce/service"
	"TA-management/internal/modules/shared/dto/testutils"

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
	t.Helper()
	for _, table := range tables {
		_, err := db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE;")
		if err != nil {
			t.Fatalf("Failed to truncate table %s: %v", table, err)
		}
	}
}

func setupAnnounceService(db *sql.DB) service.AnnouncementService {
	repo := annrepo.NewAnnouncementRepository(db)
	// discordClient = nil; tests that require discord will test error paths only
	return service.NewAnnouncementService(repo, nil)
}

// insertEmailHistory inserts a row directly into email_history for test setup.
func insertEmailHistory(t *testing.T, db *sql.DB, subject, body, receivedName string, nReceived, statusID int) int {
	t.Helper()
	var id int
	err := db.QueryRow(
		`INSERT INTO email_history (subject, body, received_name, n_received, status_ID, created_date)
		 VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id`,
		subject, body, receivedName, nReceived, statusID,
	).Scan(&id)
	assert.NoError(t, err)
	return id
}

// ─── GetEmailHistory Integration ──────────────────────────────────────────────

func TestGetEmailHistory_Integration(t *testing.T) {
	if testDB == nil {
		t.Fatal("testDB is not initialized")
	}

	svc := setupAnnounceService(testDB)

	t.Run("IT001_EmptyHistory_ReturnsEmptySlice", func(t *testing.T) {
		cleanDB(t, testDB, "email_history")

		result, err := svc.GetEmailHistory()

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 0, len(*result))
	})

	t.Run("IT002_WithEmailHistory_ReturnsLatestFirst", func(t *testing.T) {
		cleanDB(t, testDB, "email_history")
		// statusID 1 and 2 are seeded in the status lookup table
		insertEmailHistory(t, testDB, "First Email", "Body A", "All TA", 3, 1)
		insertEmailHistory(t, testDB, "Second Email", "Body B", "Course A", 1, 2)

		result, err := svc.GetEmailHistory()

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 2, len(*result))
		// ORDER BY created_date DESC → Second Email inserted last should appear first
		assert.Equal(t, "Second Email", (*result)[0].Subject)
	})

	t.Run("IT003_MoreThanTenRecords_ReturnsMaxTen", func(t *testing.T) {
		cleanDB(t, testDB, "email_history")
		for i := 0; i < 12; i++ {
			insertEmailHistory(t, testDB, "Subject", "Body", "Recipient", 1, 1)
		}

		result, err := svc.GetEmailHistory()

		assert.NoError(t, err)
		assert.NotNil(t, result)
		// Repo LIMIT 10
		assert.Equal(t, 10, len(*result))
	})
}

// ─── JoinDiscordChannel Integration ──────────────────────────────────────────

func TestJoinDiscordChannel_Integration(t *testing.T) {
	if testDB == nil {
		t.Fatal("testDB is not initialized")
	}

	svc := setupAnnounceService(testDB)

	t.Run("IT004_CourseWithNoDiscordEntry_ReturnsError", func(t *testing.T) {
		// Course 999999 has no discord channel record → repo returns error
		inviteLink, err := svc.JoinDiscordChannel(999999)
		assert.Empty(t, inviteLink)
		assert.Error(t, err)
	})
}

// ─── SendMailToAllCourse Integration (fire-and-forget goroutine) ──────────────

func TestSendMailToAllCourse_Integration(t *testing.T) {
	if testDB == nil {
		t.Fatal("testDB is not initialized")
	}

	svc := setupAnnounceService(testDB)

	t.Run("IT005_NoTAInDB_ShouldNotPanic", func(t *testing.T) {
		// Even if ta_courses is empty, the service dispatches email in a goroutine.
		// We just verify it does not panic.
		rq := annrequest.MailForAllCourse{Subject: "Integration Test", Body: "Hello All"}
		assert.NotPanics(t, func() {
			svc.SendMailToAllCourse(rq)
		})
	})
}

// ─── SendMailToCourse Integration ────────────────────────────────────────────

func TestSendMailToCourse_Integration(t *testing.T) {
	if testDB == nil {
		t.Fatal("testDB is not initialized")
	}

	svc := setupAnnounceService(testDB)

	t.Run("IT006_ExistingCourse_ShouldNotPanic", func(t *testing.T) {
		rq := annrequest.MailForCourse{CourseID: 1, Subject: "Course Announcement", Body: "Details here"}
		assert.NotPanics(t, func() {
			svc.SendMailToCourse(rq)
		})
	})
}

// ─── SendMailToTA Integration ─────────────────────────────────────────────────

func TestSendMailToTA_Integration(t *testing.T) {
	if testDB == nil {
		t.Fatal("testDB is not initialized")
	}

	svc := setupAnnounceService(testDB)

	t.Run("IT007_ExistingStudent_ShouldNotPanic", func(t *testing.T) {
		rq := annrequest.MailForTA{StudentID: 1, Subject: "TA Notice", Body: "Notice details"}
		assert.NotPanics(t, func() {
			svc.SendMailToTA(rq)
		})
	})

	t.Run("IT008_NonExistentStudent_ShouldNotPanic", func(t *testing.T) {
		rq := annrequest.MailForTA{StudentID: 999999, Subject: "TA Notice", Body: "Notice details"}
		assert.NotPanics(t, func() {
			svc.SendMailToTA(rq)
		})
	})
}
