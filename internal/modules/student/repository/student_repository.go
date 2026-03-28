package repository

import (
	"database/sql"
	"fmt"

	"TA-management/internal/modules/student/dto/response"

	"github.com/jmoiron/sqlx"
)

type StudentRepository interface {
	GetStudentByID(studentID int) (*response.StudentProfile, error)
	UpdateStudent(studentID int, firstnameThai, lastnameThai, phoneNumber string) error
	UploadDocument(studentID int, docType string, fileBytes []byte, fileName string) error
	GetDocument(studentID int, docType string) ([]byte, string, error)
	DeleteDocument(studentID int, docType string) error
	CheckDocumentExists(studentID int, docType string) (bool, error)
}

type studentRepositoryImpl struct {
	db *sqlx.DB
}

func NewStudentRepository(db *sqlx.DB) StudentRepository {
	return &studentRepositoryImpl{db: db}
}

func (r *studentRepositoryImpl) GetStudentByID(studentID int) (*response.StudentProfile, error) {
	var (
		id              int
		fnameThai       string
		lnameThai       string
		email           string
		phoneNumber     string
		transcriptName  sql.NullString
		bankAccountName sql.NullString
		studentCardName sql.NullString
	)

	// Get student basic info and document filenames using LEFT JOIN
	query := `
		SELECT 
			s.student_ID, 
			COALESCE(s.firstname_thai, '') as firstname_thai, 
			COALESCE(s.lastname_thai, '') as lastname_thai, 
			COALESCE(s.email, '') as email, 
			COALESCE(s.phone_number, '') as phone_number,
			t.file_name as transcript_filename,
			b.file_name as bank_account_filename,
			c.file_name as student_card_filename
		FROM students s
		LEFT JOIN transcript_storage t ON s.student_ID = t.student_ID
		LEFT JOIN bank_account_storage b ON s.student_ID = b.student_ID
		LEFT JOIN student_card_storage c ON s.student_ID = c.student_ID
		WHERE s.student_ID = $1
	`
	err := r.db.QueryRow(query, studentID).Scan(
		&id,
		&fnameThai,
		&lnameThai,
		&email,
		&phoneNumber,
		&transcriptName,
		&bankAccountName,
		&studentCardName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("student not found")
		}
		return nil, err
	}

	profile := &response.StudentProfile{
		StudentID:           id,
		FirstnameThai:       fnameThai,
		LastnameThai:        lnameThai,
		Email:               email,
		PhoneNumber:         phoneNumber,
		HasTranscript:       transcriptName.Valid,
		TranscriptFileName:  transcriptName.String,
		HasBankAccount:      bankAccountName.Valid,
		BankAccountFileName: bankAccountName.String,
		HasStudentCard:      studentCardName.Valid,
		StudentCardFileName: studentCardName.String,
	}

	return profile, nil
}

func (r *studentRepositoryImpl) UpdateStudent(studentID int, firstnameThai, lastnameThai, phoneNumber string) error {
	query := `
		UPDATE students
		SET firstname_thai = $1, lastname_thai = $2, phone_number = $3
		WHERE student_ID = $4
	`
	result, err := r.db.Exec(query, firstnameThai, lastnameThai, phoneNumber, studentID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("student not found")
	}

	return nil
}

func (r *studentRepositoryImpl) UploadDocument(studentID int, docType string, fileBytes []byte, fileName string) error {
	var query string
	var tableName string

	switch docType {
	case "transcript":
		tableName = "transcript_storage"
	case "bank_account":
		tableName = "bank_account_storage"
	case "student_card":
		tableName = "student_card_storage"
	default:
		return fmt.Errorf("invalid document type")
	}

	// Check if document exists
	exists, err := r.CheckDocumentExists(studentID, docType)
	if err != nil {
		return err
	}

	if exists {
		// Update existing document
		query = fmt.Sprintf(`
			UPDATE %s
			SET file_bytes = $1, file_name = $2
			WHERE student_ID = $3
		`, tableName)
	} else {
		// Insert new document
		query = fmt.Sprintf(`
			INSERT INTO %s (file_bytes, file_name, student_ID)
			VALUES ($1, $2, $3)
		`, tableName)
	}

	_, err = r.db.Exec(query, fileBytes, fileName, studentID)
	return err
}

func (r *studentRepositoryImpl) GetDocument(studentID int, docType string) ([]byte, string, error) {
	var tableName string
	var fileBytes []byte
	var fileName string

	switch docType {
	case "transcript":
		tableName = "transcript_storage"
	case "bank_account":
		tableName = "bank_account_storage"
	case "student_card":
		tableName = "student_card_storage"
	default:
		return nil, "", fmt.Errorf("invalid document type")
	}

	query := fmt.Sprintf(`
		SELECT file_bytes, file_name
		FROM %s
		WHERE student_ID = $1
	`, tableName)

	err := r.db.QueryRow(query, studentID).Scan(&fileBytes, &fileName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", fmt.Errorf("document not found")
		}
		return nil, "", err
	}

	return fileBytes, fileName, nil
}

func (r *studentRepositoryImpl) DeleteDocument(studentID int, docType string) error {
	var tableName string

	switch docType {
	case "transcript":
		tableName = "transcript_storage"
	case "bank_account":
		tableName = "bank_account_storage"
	case "student_card":
		tableName = "student_card_storage"
	default:
		return fmt.Errorf("invalid document type")
	}

	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE student_ID = $1
	`, tableName)

	result, err := r.db.Exec(query, studentID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("document not found")
	}

	return nil
}

func (r *studentRepositoryImpl) CheckDocumentExists(studentID int, docType string) (bool, error) {
	var tableName string

	switch docType {
	case "transcript":
		tableName = "transcript_storage"
	case "bank_account":
		tableName = "bank_account_storage"
	case "student_card":
		tableName = "student_card_storage"
	default:
		return false, fmt.Errorf("invalid document type")
	}

	query := fmt.Sprintf(`
		SELECT EXISTS(SELECT 1 FROM %s WHERE student_ID = $1)
	`, tableName)

	var exists bool
	err := r.db.QueryRow(query, studentID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
