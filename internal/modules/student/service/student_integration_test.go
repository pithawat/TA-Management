package service_test

import (
	"bytes"
	"database/sql"
	"io"
	"os"
	"testing"

	"TA-management/internal/modules/shared/dto/testutils"
	sreq "TA-management/internal/modules/student/dto/request"
	studentrepo "TA-management/internal/modules/student/repository"
	"TA-management/internal/modules/student/service"

	"github.com/jmoiron/sqlx"
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

func setupStudentService(db *sql.DB) service.StudentService {
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := studentrepo.NewStudentRepository(sqlxDB)
	return service.NewStudentService(repo)
}

func toReader(b []byte) io.Reader {
	return bytes.NewReader(b)
}

// ─── GetStudentProfile Integration ───────────────────────────────────────────

func TestGetStudentProfile_Integration(t *testing.T) {
	if testDB == nil {
		t.Fatal("testDB is not initialized")
	}

	svc := setupStudentService(testDB)

	t.Run("IT001_ExistingStudent_ReturnsProfile", func(t *testing.T) {
		// Student ID=1 is a seed student in the test DB (from professors table foreign key)
		profile, err := svc.GetStudentProfile(1)

		assert.NoError(t, err)
		assert.NotNil(t, profile)
		assert.Equal(t, 1, profile.StudentID)
	})

	t.Run("IT002_NonExistentStudent_ReturnsError", func(t *testing.T) {
		profile, err := svc.GetStudentProfile(999999)

		assert.Nil(t, profile)
		assert.EqualError(t, err, "student not found")
	})
}

// ─── UpdateStudentProfile Integration ─────────────────────────────────────────

func TestUpdateStudentProfile_Integration(t *testing.T) {
	if testDB == nil {
		t.Fatal("testDB is not initialized")
	}

	svc := setupStudentService(testDB)

	t.Run("IT003_UpdateSuccess_VerifyInDB", func(t *testing.T) {
		req := &sreq.UpdateProfile{
			FirstnameThai: "สมชาย",
			LastnameThai:  "ทดสอบ",
			PhoneNumber:   "0812345678",
		}

		err := svc.UpdateStudentProfile(1, req)
		assert.NoError(t, err)

		// Verify directly in DB
		var fname, lname, phone string
		dbErr := testDB.QueryRow(
			"SELECT COALESCE(firstname_thai,''), COALESCE(lastname_thai,''), COALESCE(phone_number,'') FROM students WHERE student_ID = $1", 1,
		).Scan(&fname, &lname, &phone)
		assert.NoError(t, dbErr)
		assert.Equal(t, "สมชาย", fname)
		assert.Equal(t, "ทดสอบ", lname)
		assert.Equal(t, "0812345678", phone)
	})

	t.Run("IT004_NonExistentStudent_ShouldFail", func(t *testing.T) {
		req := &sreq.UpdateProfile{
			FirstnameThai: "สมชาย",
			LastnameThai:  "ทดสอบ",
			PhoneNumber:   "0812345678",
		}
		err := svc.UpdateStudentProfile(999999, req)
		assert.EqualError(t, err, "student not found")
	})
}

// ─── Document Lifecycle Integration ──────────────────────────────────────────

func TestUploadDocument_Integration(t *testing.T) {
	if testDB == nil {
		t.Fatal("testDB is not initialized")
	}

	svc := setupStudentService(testDB)

	t.Run("IT005_UploadTranscript_Success", func(t *testing.T) {
		// Ensure clean slate
		testDB.Exec("DELETE FROM transcript_storage WHERE student_ID = 1")

		content := []byte("fake-pdf-content")
		err := svc.UploadDocument(1, "transcript", toReader(content), "transcript.pdf")
		assert.NoError(t, err)

		// Verify in DB
		var count int
		testDB.QueryRow("SELECT COUNT(*) FROM transcript_storage WHERE student_ID = 1").Scan(&count)
		assert.Equal(t, 1, count)
	})

	t.Run("IT006_UpdateExistingTranscript_ShouldReplaceNotDuplicate", func(t *testing.T) {
		// Upload second time — should update, not insert a new row
		content := []byte("updated-pdf")
		err := svc.UploadDocument(1, "transcript", toReader(content), "transcript_v2.pdf")
		assert.NoError(t, err)

		// Must still be exactly 1 row
		var count int
		testDB.QueryRow("SELECT COUNT(*) FROM transcript_storage WHERE student_ID = 1").Scan(&count)
		assert.Equal(t, 1, count)
	})
}

func TestGetDocument_Integration(t *testing.T) {
	if testDB == nil {
		t.Fatal("testDB is not initialized")
	}

	svc := setupStudentService(testDB)

	t.Run("IT007_GetExistingDocument_ReturnsBytes", func(t *testing.T) {
		testDB.Exec("DELETE FROM transcript_storage WHERE student_ID = 1")
		content := []byte("hello-pdf")
		_ = svc.UploadDocument(1, "transcript", toReader(content), "hello.pdf")

		gotBytes, gotName, err := svc.GetDocument(1, "transcript")
		assert.NoError(t, err)
		assert.Equal(t, content, gotBytes)
		assert.Equal(t, "hello.pdf", gotName)
	})

	t.Run("IT008_GetNonExistentDocument_ReturnsError", func(t *testing.T) {
		testDB.Exec("DELETE FROM bank_account_storage WHERE student_ID = 1")
		_, _, err := svc.GetDocument(1, "bank_account")
		assert.EqualError(t, err, "document not found")
	})
}

func TestDeleteDocument_Integration(t *testing.T) {
	if testDB == nil {
		t.Fatal("testDB is not initialized")
	}

	svc := setupStudentService(testDB)

	t.Run("IT009_DeleteExistingDocument_Success", func(t *testing.T) {
		testDB.Exec("DELETE FROM transcript_storage WHERE student_ID = 1")
		content := []byte("to-be-deleted")
		_ = svc.UploadDocument(1, "transcript", toReader(content), "del.pdf")

		err := svc.DeleteDocument(1, "transcript")
		assert.NoError(t, err)

		// Verify gone from DB
		var count int
		testDB.QueryRow("SELECT COUNT(*) FROM transcript_storage WHERE student_ID = 1").Scan(&count)
		assert.Equal(t, 0, count)
	})

	t.Run("IT010_DeleteNonExistentDocument_ReturnsError", func(t *testing.T) {
		testDB.Exec("DELETE FROM bank_account_storage WHERE student_ID = 1")
		err := svc.DeleteDocument(1, "bank_account")
		assert.EqualError(t, err, "document not found")
	})
}
