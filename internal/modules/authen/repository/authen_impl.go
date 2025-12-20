package repository

import (
	"TA-management/internal/modules/authen/dto/request"
	"database/sql"
	"fmt"
	"strings"
)

type AuthenRepositoryImplementation struct {
	db *sql.DB
}

func NewAuthenRepository(DB *sql.DB) AuthenRepositoryImplementation {
	return AuthenRepositoryImplementation{db: DB}
}

func (r AuthenRepositoryImplementation) CheckUserRole(name string) (string, error) {
	parts := strings.Split(name, " ")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid name format")
	}
	firstname := parts[0]
	lastname := parts[1]

	var isProf bool
	query := "SELECT EXISTS(SELECT 1 FROM professors WHERE firstname = $1 AND lastname = $2)"
	err := r.db.QueryRow(query, firstname, lastname).Scan(&isProf)
	if err != nil {
		return "", err
	}
	if isProf {
		return "PROFESSOR", nil
	}

	var isAccount bool
	query = "SELECT EXISTS(SELECT 1 FROM accountants WHERE firstname = $1 AND lastname = $2)"
	err = r.db.QueryRow(query, firstname, lastname).Scan(&isProf)
	if err != nil {
		return "", err
	}
	if isAccount {
		return "FINANCE", nil
	}

	return "STUDENT", nil
}

func (r AuthenRepositoryImplementation) AddStudent(rq request.CreateStudent) error {

	var count int
	checkQuery := `SELECT count(*) FROM students WHERE student_ID = $1`

	err := r.db.QueryRow(checkQuery, rq.StudentID).Scan(&count)
	if err != nil {
		return err
	}

	if count >= 1 {
		fmt.Println("this user already signup")
		return nil
	}

	query := `INSERT INTO students(student_ID, firstname, lastname) VALUES($1, $2, $3)`
	_, err = r.db.Exec(query,
		rq.StudentID,
		rq.Firstname,
		rq.Lastname)

	if err != nil {
		return err
	}

	return nil
}

func (r AuthenRepositoryImplementation) GetUserIDByName(name string, role string) (string, error) {

	var table string
	switch strings.ToUpper(role) {
	case "PROFESSOR":
		table = "professors"
	case "FINANCE":
		table = "accounts"
	default:
		return "", fmt.Errorf("unsupported role for name lookup: %s", role)
	}

	parts := strings.Split(name, " ")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid name format")
	}
	firstname := parts[0]
	lastname := parts[1]

	var id string
	query := fmt.Sprintf(`SELECT id FROM %s WHERE firstname = $1 AND lastname =$2`, table)

	err := r.db.QueryRow(query, firstname, lastname).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, err
}
